package csv

import (
	"bufio"
	"encoding/csv"
	"github.com/kr-juso/administrationCodeCrawler/internal/model"
	"os"
)

func SaveCsv(items []model.AdministrationCode) {
	f, err := os.Create("./administrationCode.tsv")
	if err != nil {
		panic(err)
	}

	wr := csv.NewWriter(bufio.NewWriter(f))
	wr.Comma = '\t'
	defer wr.Flush()

	wr.Write([]string{
		"행정동코드",
		"시도명",
		"시군구명",
		"읍면동명",
		"생성일자",
		"말소일자",
	})

	for _, item := range items {
		wr.Write(item.ToCsvRow())
	}
}
