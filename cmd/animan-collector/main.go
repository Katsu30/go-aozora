package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type CsvRow struct {
	VoiceID int
	Text    string
}

func main() {
	_, err := scraping()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func scraping() ([]CsvRow, error) {
	if len(os.Args) != 2 {
		return nil, errors.New("please specify the ID of the site to scrape")
	}

	threadID := os.Args[1]
	path := fmt.Sprintf("https://animanch.com/archives/%s.html", threadID)

	doc, err := goquery.NewDocument(path)
	if err != nil {
		return nil, err
	}

	pageTitle := doc.Find("article > h1").First()

	introAreaTexts := findThreadsText(doc, false)
	mainAreaTexts := findThreadsText(doc, true)

	texts := append(introAreaTexts, mainAreaTexts...)
	csvData := formatTextsToCSV(texts)

	fmt.Println(pageTitle.Text(), csvData, path)
	return nil, nil
}

func findThreadsText(doc *goquery.Document, isMainArea bool) []string {
	texts := []string{}

	section := "#maintext"
	if !isMainArea {
		section = "#introtext"
	}

	doc.Find(fmt.Sprintf("%s > .res .t_b", section)).Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		texts = append(texts, s.Text())
	})

	return texts
}

func formatTextsToCSV(texts []string) []CsvRow {
	csvData := []CsvRow{}

	for i, text := range texts {
		csvData = append(csvData, CsvRow{VoiceID: i%5 + 1, Text: text})
	}

	return csvData
}
