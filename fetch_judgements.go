package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	req, err := http.NewRequest("GET", "https://main.sci.gov.in/judgments", nil)
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

	req, err = http.NewRequest("POST", "https://main.sci.gov.in/php/captcha_num.php", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://main.sci.gov.in")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://main.sci.gov.in/judgments")
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
	captcha, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	reqbody := strings.NewReader("JBJfrom_date=01-01-2021&JBJto_date=23-04-2021&jorrop=J&ansCaptcha=" + string(captcha))
	req, err = http.NewRequest("POST", "https://main.sci.gov.in/php/v_judgments/getJBJ.php", reqbody)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://main.sci.gov.in")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://main.sci.gov.in/judgments")
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
	columnnames := "#"
	delimited := ""
	delimiter := ";"
	firstpass := true
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		c := s.Children()
		switch c.Length() {
		case 4:
			c.Each(func(i2 int, s2 *goquery.Selection) {
				if i2 == 0 {
					delimited = delimited + "\n" + s2.Text()
				}
				if i2 == 1 && firstpass {
					columnnames = columnnames + delimiter + s2.Text()
				}
				if i2 == 2 {
					delimited = delimited + delimiter + s2.Text()
				}
			})
		case 3, 2:
			c.Each(func(i2 int, s2 *goquery.Selection) {
				if i2 == 0 && firstpass {
					columnnames = columnnames + delimiter + s2.Text()
				}
				if i2 == 1 {
					delimited = delimited + delimiter + s2.Text()
				}
			})
		case 1:
			firstpass = false
		}
	})
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}
	dt := time.Now()
	f, err := os.Create("./data/judgements" + dt.Format("01022006-150405") + ".csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = fmt.Fprintf(w, columnnames+delimited)

	w.Flush()

	if err != nil {
		fmt.Println(err)
		return
	}
}
