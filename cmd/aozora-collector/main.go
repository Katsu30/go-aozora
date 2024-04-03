package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	InfoURL  string
	ZipURL   string
}

func findEntries(url string) ([]Entry, error) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)

	doc.Find("ol li a").Each(func(n int, elem *goquery.Selection) {
		token := pat.FindStringSubmatch(elem.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/cards%s.html", token[1], token[2])
		println(pageURL)
	})

	return nil, nil
}

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
}
