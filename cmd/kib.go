package cmd

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var re = regexp.MustCompile(`<img .*>`)

var suffix = []string{
	"SP",
	"S",
	"R",
}

var baseRarity = []string{
	"C",
	"CC",
	"CR",
	"FR",
	"MR",
	"PR",
	"PS",
	"R",
	"RE",
	"RR",
	"RR+",
	"TD",
	"U",
	"AR",
}

var triggersMap = map[string]string{
	"soul":     "SOUL",
	"salvage":  "COMEBACK",
	"draw":     "DRAW",
	"stock":    "POOL",
	"treasure": "TREASURE",
	"shot":     "SHOT",
	"bounce":   "RETURN",
	"gate":     "GATE",
	"standby":  "STANDBY",
	"choice":   "CHOICE",
}

func parseInt(st string) string {
	res := strings.Split(st, "：")[1]
	if strings.Contains(res, "-") {
		res = "0"
	}
	return res
}

// ExtractData extract data to card
func ExtractData(mainHTML *goquery.Selection) Card {
	var imgPlaceHolder string
	trigger := []string{}
	sa := []string{}
	ability := []string{}
	complex := mainHTML.Find("h4 span").Last().Text()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic for %v", complex)
		}
	}()
	set := strings.Split(complex, "/")[0]
	side := strings.Split(complex, "/")[1][0]
	setInfo := strings.Split(strings.Split(complex, "/")[1][1:], "-")
	setName := strings.TrimSpace(strings.Split(mainHTML.Find("h4").Text(), ") -")[1])
	imageCardURL, _ := mainHTML.Find("a img").Attr("src")
	abilityNode, _ := mainHTML.Find("span").Last().Html()
	imgURL, has := mainHTML.Find("span").Last().Find("img").Attr("src")

	if has {
		_, _imgPlaceHolder := path.Split(imgURL)
		_imgPlaceHolder = strings.Split(_imgPlaceHolder, ".")[0]
		imgPlaceHolder = fmt.Sprintf("[%v]", triggersMap[_imgPlaceHolder])
	}

	for _, line := range strings.Split(abilityNode, "<br/>") {
		ability = append(ability, re.ReplaceAllString(line, imgPlaceHolder))
	}

	infos := mainHTML.Find(".unit").Map(func(i int, s *goquery.Selection) string {
		// Color
		if s.Text() == "色：" {
			_, colorName := path.Split(s.Children().AttrOr("src", "yay"))
			return strings.ToUpper(strings.Split(colorName, ".")[0])
		}
		// Card type
		if strings.HasPrefix(s.Text(), "種類：") {
			cType := strings.TrimSpace(strings.Split(s.Text(), "種類：")[1])

			switch cType {
			case "イベント":
				return "EV"
			case "キャラ":
				return "CH"
			case "クライマックス":
				return "CX"
			}
		}
		// Soul
		if strings.HasPrefix(s.Text(), "ソウル：") {
			return strconv.Itoa(s.Children().Length())
		}
		// Trigger
		if strings.HasPrefix(s.Text(), "トリガー：") {
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				if i != 0 {
					res.WriteString(" ")
				}
				_, trigger := path.Split(ss.AttrOr("src", "yay"))
				res.WriteString(triggersMap[strings.Split(trigger, ".")[0]])
			})
			return strings.ToUpper(res.String())
		}
		// Trait
		if strings.HasPrefix(s.Text(), "特徴：") {
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				res.WriteString(ss.Text())
			})
			if strings.Contains(res.String(), "-") {
				return ""
			}
			return res.String()
		}
		return s.Text()
	})

	if infos[8] != "" {
		trigger = strings.Split(infos[8], " ")
	}

	if infos[9] != "" {
		sa = strings.Split(infos[9], "・")
	}
	card := Card{
		JpName:        strings.TrimSpace(mainHTML.Find("h4 span").First().Text()),
		Set:           set,
		SetName:       setName,
		Side:          string(side),
		Release:       setInfo[0],
		ID:            setInfo[1],
		CardType:      infos[1],
		Level:         parseInt(infos[2]),
		Colour:        infos[3],
		Power:         parseInt(infos[4]),
		Soul:          infos[5],
		Cost:          parseInt(infos[6]),
		Rarity:        strings.Split(infos[7], "：")[1],
		Trigger:       trigger,
		SpecialAttrib: sa,
		Ability:       ability,
		Version:       CardModelVersion,
		Cardcode:      complex,
		ImageURL:      imageCardURL,
	}
	return card
}

// IsbaseRarity check if a card is a C / U / R / RR
func IsbaseRarity(card Card) bool {
	for _, rarity := range baseRarity {
		if rarity == card.Rarity && isTrullyNotFoil(card) {
			return true
		}
	}
	return false
}

func isTrullyNotFoil(card Card) bool {
	for _, _suffix := range suffix {
		if strings.HasSuffix(card.ID, _suffix) {
			return false
		}
	}
	return true
}
