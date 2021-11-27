package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGetLastPage(t *testing.T) {
	f, err := os.Open("mockws/bd.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	var last = getLastPage(doc)
	if last != 69 {
		t.Errorf("%v is not last", last)
	}
}
