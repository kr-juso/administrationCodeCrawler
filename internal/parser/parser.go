package parser

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"time"

	"github.com/kr-juso/administrationCodeCrawler/internal/model"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func Parse(eucKrBytes []byte) []model.AdministrationCode {
	decodedBytes, _, err := transform.Bytes(korean.EUCKR.NewDecoder(), eucKrBytes)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(decodedBytes)

	scanner := bufio.NewScanner(r)

	regex := regexp.MustCompile("\\s{2,}")

	results := make([]model.AdministrationCode, 0)

	scanner.Scan() // skip first row
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		replacedStr := regex.ReplaceAllString(strings.Replace(line, " ", "\t", 1), "\t")
		splits := strings.Split(replacedStr, "\t")

		code := splits[0][:10] // ex: 1100000000
		city := splits[1]      // ex: 제주특별자치도
		var state string
		var town string
		if len(splits) > 3 {
			state = splits[2]
		}

		if len(splits) > 4 {
			town = splits[3]
		}

		dateStr := strings.TrimSpace(splits[len(splits)-1])
		createDate, err := time.Parse("20060102", dateStr) // yyyyMMdd
		if err != nil {
			createDate, _ = time.Parse("20060102", "99991231")
		}

		administrationCode := model.AdministrationCode{
			Code:        code,
			City:        city,
			State:       state,
			Town:        town,
			CreateDate:  createDate,
			DestroyDate: nil,
		}

		results = append(results, administrationCode)
	}

	return results
}
