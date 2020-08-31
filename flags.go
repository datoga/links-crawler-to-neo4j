package main

import "flag"

var uriFlag = flag.String("uri", "", "URI to crawl")

var levelsFlag = flag.Int("levels", 1, "Number of levels to crawl")

var asyncFlag = flag.Bool("async", false, "Wether if the crawling is async or not")

func init() {
	flag.Parse()
}
