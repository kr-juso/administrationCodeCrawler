package main

import (
	"github.com/kr-juso/administrationCodeCrawler/internal/cralwer"
	"github.com/kr-juso/administrationCodeCrawler/internal/csv"
	"github.com/kr-juso/administrationCodeCrawler/internal/parser"
)

func main() {
	c := crawler.DownloadCode()

	items := parser.Parse(c)
	csv.SaveCsv(items)
}
