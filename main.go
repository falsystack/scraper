package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

var baseURL = "https://www.wantedly.com/projects?type=mixed&occupation_types%5B%5D=jp__engineering&keywords%5B%5D=golang&page=1"

func main() {
	pages := getPages()
	for i := 1; i <= pages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageURL := baseURL + "&page=" + strconv.Itoa(page)
	fmt.Println("リクエストURL：", pageURL)

	resp, err := http.Get(pageURL)
	hasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	hasErr(err)

	doc.Find(".projects-index-single").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("data-project-id")
		title := s.Find(".project-title").Text()
		excerpt := s.Find(".project-excerpt").Text()
		fmt.Println(id, title, excerpt)
	})
}

func getPages() int {
	pages := 0
	resp, err := http.Get(baseURL)
	hasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	hasErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// paginationクラス中の全てのaタグを探し出す。
		pages = s.Find("a").Length()
	})
	return pages
}

func hasErrCodes(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln("request failed with status : ", resp.StatusCode)
	}
}

func hasErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
