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

	"golang.org/x/net/publicsuffix"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch cards",
    Long: `Fetch cards

Use global switches to specify the set, by default it will fetch all sets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fetch called")
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}
		page := 1
		client := &http.Client{Jar: jar}
		values := url.Values{
			"cmd":             {"search"},
			"show_page_count": {"100"},
		}
		if serieNumber != "" {
			values.Add("expansion", serieNumber)
		}
		if neo != "" {
			values.Add("title_number", fmt.Sprintf("##%v##", neo))
		}
		for {
			resp, err := client.PostForm(fmt.Sprintf("%v?page=%d", Baseurl, page), values)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode == 404 {
				break
			}
			log.Println("Fetch page : ", page, "with params : ", values)
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			doc.Find(".search-result-table tr").Each(func(i int, s *goquery.Selection) {
				var buffer bytes.Buffer
				card := ExtractData(s)

				if !allRarity {
					if !IsbaseRarity(card) {
						return
					}
				}

				res, errMarshal := json.Marshal(card)
				if errMarshal != nil {
					log.Println(errMarshal)
				}
				// fmt.Println(fmt.Sprintf("%v-%v%v-%v.json", card.Set, card.Side, card.Release, card.ID))
				var cardName = fmt.Sprintf("%v-%v%v-%v.json", card.Set, card.Side, card.Release, card.ID)
				var dirName = filepath.Join(card.Set, fmt.Sprintf("%v%v", card.Side, card.Release))
				os.MkdirAll(dirName, 0744)
				out, err := os.Create(filepath.Join(dirName, cardName))
				if err != nil {
					log.Println(err.Error())
				}
				defer out.Close()
				json.Indent(&buffer, res, "", "\t")
				buffer.WriteTo(out)
			})
			page = page + 1

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
}
