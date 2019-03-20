package cmd

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractData extract data to card
func ExtractData(s *goquery.Selection) Card {
	complex := s.Find("h4 span").Last().Text()
	// OVL/S62-001
	// sdf
	set := strings.Split(complex, "/")[0]
	side := strings.Split(complex, "/")[1][0]
	setInfo := strings.Split(strings.Split(complex, "/")[1][1:], "-")

	infos := s.Find(".unit").Map(func(i int, s *goquery.Selection) string {
		if s.Text() == "色：" {
			_, colorName := path.Split(s.Children().AttrOr("src", "yay"))
			return strings.Split(colorName, ".")[0]
		}
		if s.Text() == "ソウル：" {
			return strconv.Itoa(s.Length())
		}
		return s.Text()
	})

	card := Card{
		JpName:  strings.TrimSpace(s.Find("h4 span").First().Text()),
		Set:     set,
		SetName: set,
		Side:    string(side),
		Release: setInfo[0],
		ID:      setInfo[1],
		Level:   strings.Split(infos[2], "：")[1],
		Color:   infos[3],
		Power:   strings.Split(infos[4], "：")[1],
		Soul:    infos[5],
	}
	fmt.Println(infos)
	return card

}
