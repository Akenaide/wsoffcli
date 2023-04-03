// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"golang.org/x/net/publicsuffix"

	"github.com/Akenaide/biri"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxWorker int = 5

type furniture struct {
	Jobs      chan string
	Values    url.Values
	Kanseru   *bool
	Wg        *sync.WaitGroup
	Jar       http.CookieJar
	Transport *http.Transport
}

func responseWorker(
	id int,
	furni furniture,
	respChannel chan *http.Response,
	writeChan chan *goquery.Selection,
	retry chan<- string,
) {
	for resp := range respChannel {
		log.Printf("Start page: %v", resp.Request.URL)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			retry <- resp.Request.URL.String()
			log.Println("goquery error: ", err)
			furni.Wg.Done()
			continue
		}
		resultTable := doc.Find(".search-result-table tr")

		if resultTable.Length() == 0 && resp.StatusCode == 200 {
			*furni.Kanseru = true
		} else {
			resultTable.Each(func(i int, s *goquery.Selection) {
				furni.Wg.Add(1)
				writeChan <- s
			})
		}
		furni.Wg.Done()
		log.Printf("Finish page: %v", resp.Request.URL)
	}
}

func writeWorker(id int, furni furniture, writeChan chan *goquery.Selection) {
	for s := range writeChan {

		card := ExtractData(s)

		if !allRarity {
			if !IsbaseRarity(card) {
				furni.Wg.Done()
				continue
			}
		}

		res, errMarshal := json.Marshal(card)
		if errMarshal != nil {
			log.Println("error marshal", errMarshal)
			furni.Wg.Done()
			continue
		}
		var buffer bytes.Buffer
		cardName := fmt.Sprintf("%v-%v%v-%v.json", card.Set, card.Side, card.Release, card.ID)
		dirName := filepath.Join(card.Set, fmt.Sprintf("%v%v", card.Side, card.Release))
		os.MkdirAll(dirName, 0744)
		out, err := os.Create(filepath.Join(dirName, cardName))
		if err != nil {
			log.Println("write error", err.Error())
			furni.Wg.Done()
			continue
		}
		json.Indent(&buffer, res, "", "\t")
		buffer.WriteTo(out)
		out.Close()
		furni.Wg.Done()
		log.Println("Finish card- : ", cardName)
	}
}

func worker(id int, furni furniture, respChannel chan *http.Response, retry chan<- string) {
	for link := range furni.Jobs {

		log.Println("ID :", id, "Fetch page : ", link, "with params : ", furni.Values)
		proxy := biri.GetClient()
		proxy.Client.Jar = furni.Jar

		resp, err := proxy.Client.PostForm(link, furni.Values)
		if err != nil || resp.StatusCode != 200 {
			log.Println("Ban proxy:", err)
			proxy.Ban()
			retry <- link
		} else {
			if resp.StatusCode == 302 {
				*furni.Kanseru = true
				log.Printf("Kanseru by : %v", link)
			} else {
				proxy.Readd()
				respChannel <- resp
			}
		}
	}
	log.Println("Nani", id)
}

func getLastPage(doc *goquery.Document) int {
	fmt.Print(doc.Filter(".pager").Html())
	all := doc.Find(".pager .next")

	all.Each(func(i int, s *goquery.Selection) {
		log.Printf("getLastPage %v/ text: %v\n", i, s.Text())
	})

	last, _ := strconv.Atoi(all.Prev().First().Text())
	log.Printf("Go for %v pages\n", last)
	return last
}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch cards",
	Long: `Fetch cards

Use global switches to specify the set, by default it will fetch all sets.`,
	Run: func(cmd *cobra.Command, args []string) {
		page := viper.GetInt("page")
		fmt.Println("fetch called")
		log.Printf("Starting from page %v\n", page)
		biri.Config.PingServer = "https://ws-tcg.com/"
		biri.Config.TickMinuteDuration = 2
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		kanseru := false
		respChannel := make(chan *http.Response)
		writeChannel := make(chan *goquery.Selection)
		jobs := make(chan string)
		retry := make(chan string, 50)

		biri.ProxyStart()

		values := url.Values{
			"cmd":             {"search"},
			"show_page_count": {"100"},
			"show_small":      {"0"},
		}
		if serieNumber != "" {
			values.Add("expansion", serieNumber)
		}
		if neo != "" {
			values.Add("title_number", fmt.Sprintf("##%v##", neo))
		}

		furni := furniture{
			Jobs:    jobs,
			Kanseru: &kanseru,
			Values:  values,
			Wg:      &wg,
			Jar:     jar,
		}

		proxy := biri.GetClient()
		proxy.Client.Jar = furni.Jar

		resp, err := http.PostForm(fmt.Sprintf("%v?page=%d", Baseurl, 1), furni.Values)
		if err != nil {
			log.Fatal("Error on getting last page")
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal("Error on getting last page parse")
		}
		maxPage := getLastPage(doc)
		loop_num := maxPage - page + 1
		log.Printf("Number of loop %v\n", loop_num)
		wg.Add(loop_num)

		for i := 0; i < maxWorker; i++ {
			go worker(i, furni, respChannel, retry)
			go writeWorker(i, furni, writeChannel)
			go writeWorker(i, furni, writeChannel)
			go responseWorker(i, furni, respChannel, writeChannel, retry)
		}

		go func() {
			for i := page; i <= maxPage; i++ {
				jobs <- fmt.Sprintf("%v?page=%d", Baseurl, i)
			}
		}()

		go func() {
			for v := range retry {
				jobs <- v
				log.Printf("Retry: %v", v)

			}
		}()

		log.Println("Waiting...")
		wg.Wait()
		close(jobs)
		biri.Done()

		if viper.GetBool("incremental") {
			// write in config
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fetchCmd.Flags().IntP("page", "p", 1, "Starting page")
	fetchCmd.Flags().BoolP("incremental", "i", false, "Save last page")

	viper.BindPFlag("page", fetchCmd.Flags().Lookup("page"))
	viper.BindPFlag("incremental", fetchCmd.Flags().Lookup("incremental"))
	// fetchCmd.Flags().BoolP("reverse", "r", false, "Reverse order")
}
