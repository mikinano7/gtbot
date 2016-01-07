package xvideos

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func Xvideos(query []string) []string {
	url := fmt.Sprintf(
		"http://jp.xvideos.com/?k=%s",
		strings.Join(query, "+"),
	)
	doc, _ := goquery.NewDocument(url)

	var arr []string
	doc.Find("div.mozaique").Children().Each(func(_ int, s *goquery.Selection) {
		cld := s.Find("div.thumbInside")
		title := cld.Find("p a[href]").Text()
		url, _ := cld.Find("a[href]").Attr("href")
		arr = append(
			arr,
			fmt.Sprintf(
				"%s - http://jp.xvideos.com/%s",
				title,
				url,
			),
		)
	})

	return arr
}
