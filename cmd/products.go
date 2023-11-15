/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Akenaide/biri"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

const PRODUCTS_URL = "https://ws-tcg.com/products/"

var BAN_PRODUCT = []string{
	"new_title_ws",
	"resale_news",
}

// ProductInfo represents the extracted information from the HTML
type ProductInfo struct {
	ReleaseDate string
	Title       string
	SetCode     string
	Image       string
}

func getResponse(url string) *http.Response {
	var resp *http.Response

	for {
		var err error
		proxy := biri.GetClient()
		resp, err = proxy.Client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			log.Println("Error on fetch page: ", err)
			proxy.Ban()
			continue
		}
		proxy.Readd()
		break
	}

	return resp
}

func extractProductInfo(url string) ProductInfo {
	resp := getResponse(url)

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error parse product page: ", err)
		return ProductInfo{}
	}

	releaseDate := strings.Split(strings.TrimSpace(doc.Find(".release strong").Text()), "(")[0]
	titleAndWorkNumber := strings.TrimSpace(doc.Find(".release").Text())

	titleAndWorkNumberArray := strings.Split(titleAndWorkNumber, "/ ")
	// title := strings.Split(titleAndWorkNumberArray[0], "：")[1]
	setCode := strings.Split(titleAndWorkNumberArray[1], "：")[1]

	// Remove last char "】"
	setCode = strings.Replace(setCode, "】", "", -1)
	return ProductInfo{
		ReleaseDate: releaseDate,
		Title:       doc.Find(".entry-content > h3").Text(),
		SetCode:     setCode,
		Image:       doc.Find(".product-detail .alignright img").AttrOr("src", "notfound"),
	}
}

func fetchProduct() {
	productList := []ProductInfo{}
	resp := getResponse(PRODUCTS_URL)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error on parsing page ", err)

	}
	defer resp.Body.Close()

	doc.Find(".product-list .show-detail a").Each(func(i int, s *goquery.Selection) {
		productDetail := s.AttrOr("href", "nope")
		for _, ban := range BAN_PRODUCT {
			if strings.Contains(productDetail, ban) {
				return
			}
		}
		log.Println("Extract :", productDetail)
		productList = append(productList, extractProductInfo(productDetail))
	})

	res, errMarshal := json.Marshal(productList)
	if errMarshal != nil {
		log.Println("error marshal", errMarshal)
	}
	var buffer bytes.Buffer
	out, err := os.Create("product.json")
	if err != nil {
		log.Println("write error", err.Error())
	}
	json.Indent(&buffer, res, "", "\t")
	buffer.WriteTo(out)
	out.Close()
	log.Println("Finished")
}

// productsCmd represents the products command
var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("products called")
		biri.Config.PingServer = "https://ws-tcg.com/"
		biri.Config.TickMinuteDuration = 1
		biri.Config.Timeout = 25
		biri.ProxyStart()

		fetchProduct()

	},
}

func init() {
	rootCmd.AddCommand(productsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// productsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// productsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
