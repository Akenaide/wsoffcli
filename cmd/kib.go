package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
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

func processInt(st string) string {
	if strings.Contains(st, "-") {
		st = "0"
	}
	return st
}

// ExtractData extract data to card
func ExtractData(config siteConfig, mainHTML *goquery.Selection) Card {
	var imgPlaceHolder string
	ability := []string{}
	complex := mainHTML.Find("h4 span").Last().Text()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic for %v. Error=%v", complex, err)
		}
	}()
	log.Println("Start card:", complex)
	var set string
	var setInfo []string
	if strings.Contains(complex, "/") {
		set = strings.Split(complex, "/")[0]
		setInfo = strings.Split(strings.Split(complex, "/")[1], "-")
	} else {
		// TODO: deal with "BSF2024-03 PR" and similar cards
		log.Println("Can't get set info from:", complex)
	}
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

	infos := make(map[string]string)
	mainHTML.Find(".unit").Each(func(i int, s *goquery.Selection) {
		txt := strings.TrimSpace(s.Text())
		switch {
		// Color
		case strings.HasPrefix(txt, "色：") || strings.HasPrefix(txt, "[Color]:"):
			_, colorName := path.Split(s.Children().AttrOr("src", "yay"))
			infos["color"] = strings.ToUpper(strings.Split(colorName, ".")[0])
			// Card type
		case strings.HasPrefix(txt, "種類：") || strings.HasPrefix(txt, "[Card Type]:"):
			cType := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "種類："), "[Card Type]:"))

			switch cType {
			case "イベント", "Event":
				infos["type"] = "EV"
			case "キャラ", "Character":
				infos["type"] = "CH"
			case "クライマックス", "Climax":
				infos["type"] = "CX"
			}
			// Cost
		case strings.HasPrefix(txt, "コスト：") || strings.HasPrefix(txt, "[Cost]:"):
			cost := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "コスト："), "[Cost]:"))
			infos["cost"] = cost
			// Flavor text
		case strings.HasPrefix(txt, "フレーバー：") || strings.HasPrefix(txt, "[Flavor Text]:"):
			flvr := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "フレーバー："), "[Flavor Text]:"))
			infos["flavourText"] = flvr
			// Level
		case strings.HasPrefix(txt, "レベル：") || strings.HasPrefix(txt, "[Level]:"):
			lvl := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "レベル："), "[Level]:"))
			infos["level"] = lvl
			// Power
		case strings.HasPrefix(txt, "パワー：") || strings.HasPrefix(txt, "[Power]:"):
			pwr := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "パワー："), "[Power]:"))
			infos["power"] = pwr
			// Rarity
		case strings.HasPrefix(txt, "レアリティ：") || strings.HasPrefix(txt, "[Rarity]:"):
			rarity := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(txt, "レアリティ："), "[Rarity]:"))
			infos["rarity"] = rarity
			// Side
		case strings.HasPrefix(txt, "サイド：") || strings.HasPrefix(txt, "[Side]:"):
			_, side := path.Split(s.Children().AttrOr("src", "yay"))
			infos["side"] = strings.ToUpper(strings.Split(side, ".")[0])
			// Soul
		case strings.HasPrefix(txt, "ソウル：") || strings.HasPrefix(txt, "[Soul]:"):
			infos["soul"] = strconv.Itoa(s.Children().Length())
			// Trigger
		case strings.HasPrefix(txt, "トリガー：") || strings.HasPrefix(txt, "[Trigger]:"):
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				if i != 0 {
					res.WriteString(" ")
				}
				_, trigger := path.Split(ss.AttrOr("src", "yay"))
				res.WriteString(triggersMap[strings.Split(trigger, ".")[0]])
			})
			infos["trigger"] = strings.ToUpper(strings.TrimSpace(res.String()))
			// Trait
		case strings.HasPrefix(txt, "特徴：") || strings.HasPrefix(txt, "[Special Attribute]:"):
			var res bytes.Buffer
			s.Children().Each(func(i int, ss *goquery.Selection) {
				res.WriteString(strings.TrimSpace(ss.Text()))
			})
			if strings.Contains(res.String(), "-") {
				infos["specialAttribute"] = ""
			} else {
				infos["specialAttribute"] = strings.TrimSpace(res.String())
			}
		default:
			log.Println("Unknown:", txt)
		}
	})

	card := Card{
		Name:        strings.TrimSpace(mainHTML.Find("h4 span").First().Text()),
		Set:         set,
		SetName:     setName,
		Side:        infos["side"],
		CardType:    infos["type"],
		Level:       processInt(infos["level"]),
		FlavourText: infos["flavourText"],
		Colour:      infos["color"],
		Power:       processInt(infos["power"]),
		Soul:        infos["soul"],
		Cost:        processInt(infos["cost"]),
		Rarity:      infos["rarity"],
		Ability:     ability,
		Version:     CardModelVersion,
		Cardcode:    complex,
	}
	if fullURL, err := url.JoinPath(config.baseURL, imageCardURL); err == nil {
		card.ImageURL = fullURL
	} else {
		log.Printf("Couldn't form full image URL: %v\n", err)
		card.ImageURL = imageCardURL
	}
	if infos["specialAttribute"] != "" {
		card.SpecialAttrib = strings.Split(infos["specialAttribute"], "・")
	}
	if infos["trigger"] != "" {
		card.Trigger = strings.Split(infos["trigger"], " ")
	}
	if len(setInfo) > 1 {
		card.Release = setInfo[0]
		card.ID = setInfo[1]
	}
	if config.languageCode == "JP" {
		card.JpName = card.Name
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
