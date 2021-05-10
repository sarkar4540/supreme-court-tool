package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	req, err := http.NewRequest("GET", "https://main.sci.gov.in/display-board", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("If-Modified-Since", "Fri, 23 Apr 2021 17:52:15 GMT")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err = http.NewRequest("POST", "https://main.sci.gov.in/php/display/get_board.php?ctype=v", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")
	req.Header.Set("Origin", "https://main.sci.gov.in")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://main.sci.gov.in/display-board")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	for _, c := range resp.Cookies() {
		req.AddCookie(c)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	delimited := ""
	delimiter := ";"
	doc.Find(".record").Each(func(i int, s *goquery.Selection) {
		c := s.Children()
		c.Each(func(i2 int, s2 *goquery.Selection) {
			if i2 == 0 {
				delimited = delimited + strings.TrimSpace(s2.Text())
			} else {
				delimited = delimited + delimiter + strings.TrimSpace(s2.Text())
			}
		})
		delimited = delimited + "\n"
	})
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}
	dt := time.Now()
	f, err := os.Create("./data/listings" + dt.Format("01022006-150405") + ".csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = fmt.Fprintf(w, delimited)

	w.Flush()

	if err != nil {
		fmt.Println(err)
		return
	}
}
