package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var visited = make(map[string]bool)

func crawl(uri string) {

	links, _ := retrieve(uri)

	for _, l := range links {
		if !visited[l] {
			fmt.Println("Fetching", l)
			visited[uri] = true
			crawl(l)
		}
	}
}

func retrieve(uri string) ([]string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	doc, readerErr := goquery.NewDocumentFromReader(resp.Body)
	if readerErr != nil {
		fmt.Println("Error:", readerErr)
		return nil, readerErr
	}
	u, parseErr := url.Parse(uri)
	if parseErr != nil {
		fmt.Println("Error:", parseErr)
		return nil, parseErr
	}
	host := u.Host

	links := []string{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		lu, err := url.Parse(href)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if isInternalURL(host, lu) {
			links = append(links, u.ResolveReference(lu).String())
		}

	})

	return unique(links), nil
}

// insures that the link is internal
func isInternalURL(host string, lu *url.URL) bool {

	if lu.IsAbs() {
		return strings.EqualFold(host, lu.Host)
	}
	return len(lu.Host) == 0
}

// insures that there is no repetition
func unique(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
