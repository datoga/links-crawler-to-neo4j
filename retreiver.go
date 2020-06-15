package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type link struct {
	source string
	target string
}

type retriever struct {
	events  map[string][]chan link
	visited map[string]bool
}

func (b *retriever) addEvent(e string, ch chan link) {
	if b.events == nil {
		b.events = make(map[string][]chan link)
	}
	if _, ok := b.events[e]; ok {
		b.events[e] = append(b.events[e], ch)
	} else {
		b.events[e] = []chan link{ch}
	}
}

func (b *retriever) removeEvent(e string, ch chan link) {
	if _, ok := b.events[e]; ok {
		for i := range b.events[e] {
			if b.events[e][i] == ch {
				b.events[e] = append(b.events[e][:i], b.events[e][i+1:]...)
				break
			}
		}
	}
}

func (b *retriever) emit(e string, response link) {
	if _, ok := b.events[e]; ok {
		for _, handler := range b.events[e] {
			go func(handler chan link) {
				handler <- response
			}(handler)
		}
	}
}

func (b *retriever) crawl(uri string) {

	links, _ := b.retrieve(uri)

	for _, l := range links {
		if !b.visited[l] {
			b.emit("newLink", link{
				source: uri,
				target: l,
			})
			b.visited[uri] = true
			b.crawl(l)
		}
	}
}

func (b *retriever) retrieve(uri string) ([]string, error) {
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
