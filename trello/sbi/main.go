package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Board struct {
	Lists []List `json:"lists"`
	Cards []Card `json:"cards"`
}

type List struct {
	Closed bool   `json:"closed"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

type Card struct {
	Closed   bool    `json:"closed"`
	IdList   string  `json:"idList"`
	IdShort  int64   `json:"idShort"`
	Labels   []Label `json:"labels"`
	Name     string  `json:"name"`
	ShortUrl string  `json:"shortUrl"`
}

type Label struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Sbi struct {
	IdList string
	Number int64
	Name   string
	Url    string
	Impl   float64
	Review float64
}

// parseHour returns the number of hours from label text like "2.5h".
func parseHour(label string) (float64, bool) {
	l := len(label)
	last := label[l-1 : l]
	if last != "h" {
		return 0, false
	}
	hour, err := strconv.ParseFloat(label[0:l-1], 64)
	if err != nil {
		return 0, false
	}
	return hour, true
}

func parseCard(card Card) (*Sbi, bool) {
	if card.Closed {
		return nil, false
	}
	sbi := Sbi{
		IdList: card.IdList,
		Number: card.IdShort,
		Name:   card.Name,
		Url:    card.ShortUrl,
		Impl:   0,
		Review: 0,
	}
	for _, label := range card.Labels {
		if label.Color == "yellow" {
			hour, ok := parseHour(label.Name)
			if ok {
				sbi.Impl += hour
			}
		} else if label.Color == "orange" {
			hour, ok := parseHour(label.Name)
			if ok {
				sbi.Review += hour
			}
		}
	}
	return &sbi, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v backlog.json\n", os.Args[0])
		os.Exit(1)
	}
	raw, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var board Board
	err = json.Unmarshal(raw, &board)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	listMap := make(map[string]List)
	for _, list := range board.Lists {
		listMap[list.Id] = list
	}
	for _, card := range board.Cards {
		sbi, ok := parseCard(card)
		if !ok {
			continue
		}
		list := listMap[sbi.IdList]
		if !list.Closed && (sbi.Impl > 0 || sbi.Review > 0) {
			fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\n",
				list.Name, sbi.Number, sbi.Name, sbi.Url, sbi.Impl, sbi.Review)
		}
	}
}
