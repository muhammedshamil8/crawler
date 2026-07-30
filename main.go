package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	rawURL := args[0]

	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("invalid value for maxConcurrency")
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("invalid value for maxPages")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", rawURL)

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("invalid URL:", err)
		os.Exit(1)
	}
	pages := make(map[string]PageData)
	cfg := &config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawURL)

	cfg.wg.Wait()

	if err := writeJSONReport(cfg.pages, "report.json"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("============== CRAWLING COMPLETED ==============")
	for url, pageData := range pages {
		fmt.Printf("URL: %s, Heading: %s, First Paragraph: %s\n", url, pageData.Heading, pageData.FirstParagraph)
	}
}
