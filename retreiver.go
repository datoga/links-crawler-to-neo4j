package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Link struct {
	source string
	target string
}

type Retriever struct {
	Links chan Link
	depth int
	async bool
}

func NewRetriever(depth int, async bool) *Retriever {
	chLinks := make(chan Link)

	return &Retriever{
		Links: chLinks,
		depth: depth,
		async: async,
	}
}

func (b Retriever) emit(link Link) {
	b.Links <- link
}

func (b *Retriever) finish() {
	close(b.Links)
}

func (b Retriever) Crawl(uri string) {

	fmt.Println("Configuring collector")
	fmt.Println("---------------------")
	fmt.Println("Max depth:", b.depth)
	fmt.Println("Async:", b.async)
	fmt.Println("---------------------")

	c := colly.NewCollector(
		colly.MaxDepth(b.depth),
	)

	c.Async = b.async

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		parsedLink := e.Request.AbsoluteURL(href)

		if parsedLink == "" {
			return
		}

		visited, err := e.Request.HasVisited(parsedLink)

		if err != nil {
			fmt.Println(err)
			return
		}

		if visited {
			return
		}

		b.emit(Link{
			source: e.Request.URL.String(),
			target: parsedLink,
		})

		e.Request.Visit(href)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(uri)

	c.Wait()

	b.finish()
}
