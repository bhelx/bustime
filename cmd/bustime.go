package main

import (
	"github.com/bhelx/bustime"
	"github.com/bhelx/bustime/lib"
)

func main() {
	config := bustime.GetConfig()
	scraper := lib.NewScraper(config)
	defer scraper.Close()
	scraper.Start()
}
