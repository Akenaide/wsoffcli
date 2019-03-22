package cmd

import (
	"bytes"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractData extract data to card
func ExtractData(mainHtml *goquery.Selection) Card {
	complex := mainHtml.Find("h4 span").Last().Text()
	set := strings.Split(complex, "/")[0]
	side := strings.Split(complex, "/")[1][0]
	setInfo := strings.Split(strings.Split(complex, "/")[1][1:], "-")
	var ability, _ = mainHtml.Find("span").Last().Html()

	infos := mainHtml.Find(".unit").Map(func(i int, s *goquery.Selection) string {
		if s.Text() == "色：" {
			_, colorName := path.Split(s.Children().AttrOr("src", "yay"))
			return strings.Split(colorName, ".")[0]
		}
		if strings.HasPrefix(s.Text(), "種類：") {
			var cType = strings.TrimSpace(strings.Split(s.Text(), "種類：")[1])

			switch cType {
			case "イベント":
				return "Event"
			case "キャラ":
				return "Character"
			case "クライマックス":
				return "Climax"
			}
		}
		if s.Text() == "ソウル：" {
			return strconv.Itoa(s.Children().Length())
		}
		if strings.HasPrefix(s.Text(), "トリガー：") {
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				if i != 0 {
					res.WriteString(" ")
				}
				_, trigger := path.Split(ss.AttrOr("src", "yay"))
				res.WriteString(strings.Split(trigger, ".")[0])
			})
			return strings.ToUpper(res.String())
		}
		if strings.HasPrefix(s.Text(), "特徴：") {
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				res.WriteString(ss.Text())
			})
			return res.String()
		}
		return s.Text()
	})

	card := Card{
		JpName:        strings.TrimSpace(mainHtml.Find("h4 span").First().Text()),
		Set:           set,
		SetName:       set,
		Side:          string(side),
		Release:       setInfo[0],
		ID:            setInfo[1],
		CardType:      infos[1],
		Level:         strings.Split(infos[2], "：")[1],
		Color:         infos[3],
		Power:         strings.Split(infos[4], "：")[1],
		Soul:          infos[5],
		Cost:          strings.Split(infos[6], "：")[1],
		Rarity:        strings.Split(infos[7], "：")[1],
		Trigger:       strings.Split(infos[8], " "),
		SpecialAttrib: strings.Split(infos[9], "・"),
		Ability:       strings.Split(ability, "<br/>"),
	}
	return card

}
