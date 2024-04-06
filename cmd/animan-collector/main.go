package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

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

	// pageTitle := doc.Find("article > h1").First()

	introAreaTexts := findThreadsText(doc, false)
	mainAreaTexts := findThreadsText(doc, true)
	commentAreaText := findCommentAreaText(doc)

	texts := append(introAreaTexts, mainAreaTexts...)
	texts = append(texts, commentAreaText...)

	csvData := formatTextsToCSV(texts)

	return csvData, nil
}

func findThreadsText(doc *goquery.Document, isMainArea bool) []string {
	texts := []string{}

	section := "#maintext"
	if !isMainArea {
		section = "#introtext"
	}

	doc.Find(fmt.Sprintf("%s > .res .t_b", section)).Each(func(i int, s *goquery.Selection) {
		// Remove empty lines
		if utf8.RuneCountInString(strings.TrimSpace(s.Text())) != 0 {
			texts = append(texts, s.Text())
		}
	})

	return texts
}

func findCommentAreaText(doc *goquery.Document) []string {
	texts := []string{}

	doc.Find(fmt.Sprintln("#commentarea > .commentwrap > .commentbody")).Each(func(i int, s *goquery.Selection) {
		// Remove empty lines
		if utf8.RuneCountInString(strings.TrimSpace(s.Text())) != 0 {
			// Replace newline characters with spaces
			text := strings.ReplaceAll(s.Text(), "\n", " ")
			texts = append(texts, text)
		}
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
