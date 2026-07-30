package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"


	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages int
}


func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	heading := doc.Find("h1").First()
	if heading.Length() > 0 {
		return heading.Text()
	}
	heading = doc.Find("h2").First()
	if heading.Length() > 0 {
		return heading.Text()
	}
	return ""
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	main := doc.Find("main")
	if main.Length() > 0 {
		return main.Find("p").First().Text()
	}
	paragraph := doc.Find("p")
	if paragraph.Length() > 0 {
		return paragraph.First().Text()
	}

	return ""
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}
		resolvedURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, resolvedURL.String())
	})
	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}
		resolvedURL := baseURL.ResolveReference(parsedURL)
		images = append(images, resolvedURL.String())
	})
	return images, nil
}

func extractPageData(html, pageURL string) PageData {
	heading := getHeadingFromHTML(html)
	paragraph := getFirstParagraphFromHTML(html)
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		parsedURL = &url.URL{}
	}
	links, err := getURLsFromHTML(html, parsedURL)
	if err != nil {
		links = []string{}
	}
	images, err := getImagesFromHTML(html, parsedURL)
	if err != nil {
		images = []string{}
	}
	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: paragraph,
		OutgoingLinks:  links,
		ImageURLs:      images,
	}
}

func getHTML(rawURL string) (string, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", errors.New("HTTP error: " + resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", errors.New("invalid content type: " + contentType)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil

}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	}
	cfg.pages[normalizedURL] = PageData{}
	return true
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	parsedURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("Error parsing current URL:", err)
		return
	}
	if parsedURL.Host != cfg.baseURL.Host {
		return
	}
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("Error normalizing current URL:", err)
		return
	}
	isFirst := cfg.addPageVisit(normalizedCurrentURL)
	if !isFirst {
		return
	}
	fmt.Printf("Crawling: %s\n", normalizedCurrentURL)
	res, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("Error fetching HTML:", err)
		return
	}
	cfg.mu.Lock()
	cfg.pages[normalizedCurrentURL] = extractPageData(res, rawCurrentURL)
	 cfg.mu.Unlock()

	urls, err := getURLsFromHTML(res, parsedURL)
	if err != nil {
		fmt.Println("Error fetching URLs:", err)
		return
	}
	for _, link := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)

	}

}
