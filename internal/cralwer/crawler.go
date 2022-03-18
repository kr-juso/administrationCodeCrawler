package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"archive/zip"
	"bytes"
	"strings"
)

const moisURI = "https://www.mois.go.kr"
const listPage = "https://www.mois.go.kr/frt/bbs/type001/commonSelectBoardList.do?bbsId=BBSMSTR_000000000052" // 주민등록, 인감, 행정사 리스트

func DownloadCode() []byte {
	// detailPageUrl := "https://www.mois.go.kr/frt/bbs/type001/commonSelectBoardArticle.do?bbsId=BBSMSTR_000000000052&nttId=89611"
	detailPageUrl := getDetailLink()
	detailBytes := downloadDetailPage(detailPageUrl)

	zipLink, zipFileName := extractZipLink(detailBytes)
	zipBytes := downloadZip(zipLink)

	filename := fmt.Sprintf("KIKcd_H.%s", zipFileName[6:14]) // yyyyMMdd from "jscode20220103.zip"
	return extractFile(zipBytes, filename)
}

func getDetailLink() string {
	res, err := http.Get(listPage)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	querySelector := ".table_wrap .table_style1 tr td a"

	path := ""
	isFirst := false

	doc.Find(querySelector).Each(func(i int, s *goquery.Selection) {
		isContain := strings.HasPrefix(s.Text(), "주민등록업무 행정기관 및 관할구역 변경내역")
		if !isFirst && isContain { // 가장 최근 변경내역 페이지만 추출
			path, _ = s.Attr("href")
			isFirst = true
		}
	})

	return fmt.Sprintf("%s%s", moisURI, path)
}

func extractZipLink(detailBytes []byte) (string, string) {
	detailReader := bytes.NewReader(detailBytes)

	doc, err := goquery.NewDocumentFromReader(detailReader)
	if err != nil {
		panic(err)
	}

	regex := regexp.MustCompile("jscode\\d+\\.zip")

	var zipFilename string
	var codeDownloadLink string
	doc.Find(".fileList li a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		isJscode := regex.FindString(text) != ""

		if isJscode {
			codeDownloadLink, _ = s.Attr("href")
			zipFilename = regex.FindString(s.Text())
		}
	})

	return fmt.Sprintf("%s%s", moisURI, codeDownloadLink), zipFilename
}

func downloadDetailPage(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("status code error.")
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return result
}

func downloadZip(link string) []byte {
	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return result
}

func extractFile(zipBytes []byte, fileName string) []byte {
	bReader := bytes.NewReader(zipBytes)

	r, err := zip.NewReader(bReader, int64(len(zipBytes)))
	if err != nil {
		panic(err)
	}

	for _, f := range r.File {
		if f.Name == fileName {
			fileReader, err := f.Open()
			if err != nil {
				panic(err)
			}

			result, err := ioutil.ReadAll(fileReader)
			if err != nil {
				return nil
			}

			return result
		}
	}

	panic("file not exist.")
}
