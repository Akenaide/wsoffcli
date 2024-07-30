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
	"strings"
	"sync"

	"golang.org/x/net/publicsuffix"

	"github.com/Akenaide/biri"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxWorker int = 5

type siteConfig struct {
	baseURL                  string
	cardListURL              string
	cardSearchURL            string
	languageCode             string
	nextButtonDistinguisher  string
	resultTableDistinguisher string
}

var siteConfigs = map[string]siteConfig{
	"EN": {
		baseURL:                  "https://en.ws-tcg.com/",
		cardListURL:              "https://en.ws-tcg.com/cardlist/list/",
		cardSearchURL:            "https://en.ws-tcg.com/cardlist/cardsearch/exec",
		languageCode:             "EN",
		nextButtonDistinguisher:  "[rel=\"next\"]",
		resultTableDistinguisher: "#searchResult-table tr",
	},
	"JP": {
		baseURL:                  "https://ws-tcg.com/",
		cardListURL:              "https://ws-tcg.com/cardlist/",
		cardSearchURL:            "https://ws-tcg.com/cardlist/search",
		languageCode:             "JP",
		nextButtonDistinguisher:  ".pager .next",
		resultTableDistinguisher: ".search-result-table tr",
	},
}

type writerWorkerStruct struct {
	mode       string
	wg         *sync.WaitGroup
	writeChan  chan *goquery.Selection
	boosterMap map[string]booster
}

type booster struct {
	code  string
	cards []Card
}

type furniture struct {
	Jobs     chan string
	Values   url.Values
	Wg       *sync.WaitGroup
	Jar      http.CookieJar
	lastPage int
}

func (w *writerWorkerStruct) run(config siteConfig) {
	switch w.mode {
	case "card":
		go w.card(config)
		go w.card(config)
	case "booster":
		go w.populateBooster(config)
	}
}

func (w *writerWorkerStruct) populateBooster(config siteConfig) {
	for s := range w.writeChan {
		newCard := ExtractData(config, s)
		boosterCode := newCard.Release
		boosterObj := w.boosterMap[boosterCode]

		boosterObj.cards = append(boosterObj.cards, newCard)
		w.boosterMap[boosterCode] = boosterObj
		w.wg.Done()
	}
}

func (w *writerWorkerStruct) write(config siteConfig) {
	fmt.Println("Start write in mode: ", w.mode)
	if w.mode == "booster" {
		for k, v := range w.boosterMap {
			log.Println("Found booster :", k)
			dirName := filepath.Join(viper.GetString("boosterDir"), config.languageCode)
			os.MkdirAll(dirName, 0o744)
			filename := filepath.Join(dirName, k+".json")
			updatedData, err := json.Marshal(v.cards)
			if err != nil {
				log.Println("Error marshal struct: ", k)
			}
			if err := os.WriteFile(filename, updatedData, 0o644); err != nil {
				log.Println("Error writing :", k)
			}
		}
	}
}

func (w *writerWorkerStruct) card(config siteConfig) {
	for s := range w.writeChan {

		card := ExtractData(config, s)

		res, errMarshal := json.Marshal(card)
		if errMarshal != nil {
			log.Println("error marshal", errMarshal)
			w.wg.Done()
			continue
		}
		var buffer bytes.Buffer
		cardName := fmt.Sprintf("%v-%v-%v.json", card.Set, card.Release, card.ID)
		dirName := filepath.Join(viper.GetString("cardDir"), config.languageCode, card.Set, card.Release)
		os.MkdirAll(dirName, 0o744)
		out, err := os.Create(filepath.Join(dirName, cardName))
		if err != nil {
			log.Println("write error", err.Error())
			w.wg.Done()
			continue
		}
		json.Indent(&buffer, res, "", "\t")
		buffer.WriteTo(out)
		out.Close()
		w.wg.Done()
		log.Println("Finish card- : ", cardName)
	}
}

func responseWorker(
	id int,
	furni *furniture,
	config siteConfig,
	respChannel chan *http.Response,
	writeChan chan *goquery.Selection,
) {
	for resp := range respChannel {
		log.Printf("Start page: %v", resp.Request.URL)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			furni.Jobs <- resp.Request.URL.String()
			log.Println("goquery error: ", err, "for page: ", resp.Request.URL)
			continue
		}
		resultTable := doc.Find(config.resultTableDistinguisher)

		if resultTable.Length() == 0 && resp.StatusCode == 200 {
			fmt.Println("No cards on response page")
		} else {
			log.Println("Found cards !!", resp.Request.URL)
			resultTable.Each(func(i int, s *goquery.Selection) {
				furni.Wg.Add(1)
				writeChan <- s
			})
		}
		furni.Wg.Done()
		log.Printf("Finish page: %v", resp.Request.URL)
	}
}

func worker(id int, furni *furniture, respChannel chan *http.Response) {
	for link := range furni.Jobs {

		log.Println("ID :", id, "Fetch page : ", link, "with params : ", furni.Values)
		proxy := biri.GetClient()
		log.Println("Got proxy")
		proxy.Client.Jar = furni.Jar

		resp, err := proxy.Client.PostForm(link, furni.Values)
		if err != nil || resp.StatusCode != 200 {
			log.Println("Ban proxy:", err)
			proxy.Ban()
			furni.Jobs <- link
		} else {
			if resp.StatusCode != 302 {
				proxy.Readd()
				respChannel <- resp
			}
		}
	}
	log.Println("Nani", id)
}

func (f *furniture) getLastPage(config siteConfig) int {
	resp, err := http.PostForm(fmt.Sprintf("%v?page=%d", config.cardSearchURL, 1), f.Values)
	if err != nil {
		log.Fatal("Error on getting last page")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error on getting last page parse")
	}
	all := doc.Find(config.nextButtonDistinguisher)

	last, _ := strconv.Atoi(all.Prev().First().Text())
	// default is 1, there no .pager .next if it's the only page
	if last == 0 {
		last = 1
	}
	log.Println("Last pages is", last, " for :", f.Values)
	f.lastPage = last
	return last
}

func getRecentFurnitures(doc *goquery.Document) []furniture {
	var furnitures = []furniture{}
	// Find all <a> elements with onclick attributes within the <ul> element
	doc.Find("div.system > ul.expansion-list a[onclick]").Each(func(i int, sel *goquery.Selection) {
		onclickAttr, exists := sel.Attr("onclick")
		if exists {
			// Extract the integer value from the onclick attribute
			parts := strings.Split(onclickAttr, "('")
			if len(parts) >= 2 {
				value := strings.TrimSuffix(parts[1], "')")
				urlValues := url.Values{
					"cmd":             {"search"},
					"show_page_count": {"100"},
					"show_small":      {"0"},
					// TODO: this should use viper
					"parallel":  {"0"},
					"expansion": {value},
				}
				furnitures = append(furnitures, furniture{Values: urlValues})
			}
		}
	})
	return furnitures
}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch cards",
	Long: `Fetch cards

Use global switches to specify the set, by default it will fetch all sets.`,
	Run: func(cmd *cobra.Command, args []string) {
		lang := viper.GetString("lang")
		var config siteConfig
		if c, ok := siteConfigs[lang]; !ok {
			log.Fatalf("Unsupported language: %q\n", lang)
		} else {
			config = c
			log.Printf("Fetching %s cards\n", lang)
		}

		iter := viper.GetInt("iter")
		loopNum := 0
		fmt.Println("fetch called")
		fmt.Printf("Settings: %v\n", viper.AllSettings())
		biri.Config.PingServer = config.baseURL
		biri.Config.TickMinuteDuration = 1
		biri.Config.Timeout = 25
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}

		furnitures := []*furniture{}

		var wg sync.WaitGroup
		respChannel := make(chan *http.Response, maxWorker)
		writeChannel := make(chan *goquery.Selection)

		biri.ProxyStart()

		values := url.Values{
			"cmd":             {"search"},
			"show_page_count": {"100"},
			"show_small":      {"0"},
			"parallel":        {"0"},
		}
		if serieNumber != "" {
			values.Add("expansion", serieNumber)
		}
		if neo != "" {
			values.Add("title_number", fmt.Sprintf("##%v##", neo))
		}

		if !viper.GetBool("allrarity") {
			values.Set("parallel", "1")
		}
		// Default furniture
		furni := furniture{
			Values: values,
			Wg:     &wg,
			Jar:    jar,
		}

		if viper.GetBool("recent") {
			resp, err := http.Get(config.cardListURL)
			if err != nil {
				log.Fatal("Error on get recent")
			}
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal("Error on parse recent")
			}
			for _, recentfurni := range getRecentFurnitures(doc) {
				copyFurni := furni
				copyFurni.Values = recentfurni.Values
				log.Println(furni, recentfurni)
				furnitures = append(furnitures, &copyFurni)
			}
		} else {
			furnitures = append(furnitures, &furni)

		}

		for _, furni := range furnitures {
			lastPage := furni.getLastPage(config)
			loopNum += lastPage
			furni.Jobs = make(chan string, lastPage)
		}

		if iter == 0 {
			iter = loopNum
		} else {
			loopNum = iter
		}

		log.Printf("Number of loop %v\n", loopNum)
		wg.Add(loopNum)

		writerWorker := writerWorkerStruct{
			mode:       viper.GetString("export"),
			wg:         &wg,
			writeChan:  writeChannel,
			boosterMap: make(map[string]booster),
		}
		writerWorker.run(config)
		for _, furni := range furnitures {
			for i := 0; i < maxWorker; i++ {
				go worker(i, furni, respChannel)
				go responseWorker(i, furni, config, respChannel, writeChannel)
			}
			go func(furni *furniture, iter int) {
				fmt.Println("furni.lastPage:", furni.lastPage)
				if viper.GetBool("reverse") {
					for i := 0; i < iter; i++ {
						// hotfix, more work needed to have a clean implementation
						if furni.lastPage-i > 0 {
							furni.Jobs <- fmt.Sprintf("%v?page=%d", config.cardSearchURL, furni.lastPage-i)
						}
					}
				} else {
					for i := 1; i <= iter; i++ {
						// hotfix, more work needed to have a clean implementation
						if i <= furni.lastPage {
							furni.Jobs <- fmt.Sprintf("%v?page=%d", config.cardSearchURL, i)
						}
					}
				}
				log.Print("Finished loop")
			}(furni, iter)
		}

		log.Println("Waiting...")
		wg.Wait()
		// close(jobs)
		biri.Done()
		writerWorker.write(config)
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
	fetchCmd.Flags().StringP("boosterDir", "", "boosters", "Directory to put fetched booster information into")
	fetchCmd.Flags().StringP("cardDir", "", "cards", "Directory to put fetched card information into")
	fetchCmd.Flags().IntP("iter", "i", 0, "Number of iteration")
	fetchCmd.Flags().BoolP("reverse", "r", false, "Reverse order")
	fetchCmd.Flags().BoolP("allrarity", "a", false, "get all rarity (sp, ssp, sbr, etc...)")
	fetchCmd.Flags().StringP("export", "e", "card", "export value: card, booster, all")
	fetchCmd.Flags().StringP("lang", "l", "JP", "Site language to pull from. Options are EN or JP. JP is default")
	fetchCmd.Flags().BoolP("recent", "", false, "get all recent products")

	viper.BindPFlag("boosterDir", fetchCmd.Flags().Lookup("boosterDir"))
	viper.BindPFlag("cardDir", fetchCmd.Flags().Lookup("cardDir"))
	viper.BindPFlag("page", fetchCmd.Flags().Lookup("page"))
	viper.BindPFlag("iter", fetchCmd.Flags().Lookup("iter"))
	viper.BindPFlag("reverse", fetchCmd.Flags().Lookup("reverse"))
	viper.BindPFlag("allrarity", fetchCmd.Flags().Lookup("allrarity"))
	viper.BindPFlag("export", fetchCmd.Flags().Lookup("export"))
	viper.BindPFlag("lang", fetchCmd.Flags().Lookup("lang"))
	viper.BindPFlag("recent", fetchCmd.Flags().Lookup("recent"))
}
