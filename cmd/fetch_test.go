package cmd

import (
	"log"
	"os"
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// func TestGetLastPage(t *testing.T) {
// 	f, err := os.Open("mockws/bd.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	doc, err := goquery.NewDocumentFromReader(f)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	last := getLastPage(doc)
// 	if last != 69 {
// 		t.Errorf("%v is not last", last)
// 	}
// }

func TestRecentSwitch(t *testing.T) {
	expectedExpension := []string{
		"444",
		"439",
		"443",
		"442",
		"438",
		"437",
		"441",
		"440",
	}
	f, err := os.Open("mockws/recent.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	recentFurnitures := getRecentFurnitures(doc)
	if len(recentFurnitures) != 8 {
		t.Errorf("Should be equal to 8: %v", recentFurnitures)
	}

	for _, furni := range recentFurnitures {
		expansion := furni.Values.Get("expansion")
		if slices.Contains(expectedExpension, expansion) == false {
			t.Errorf("Did not expect %v expansion", expansion)
		}

	}
}