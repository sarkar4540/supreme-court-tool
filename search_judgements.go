package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir("./data")
	if err != nil {
		log.Fatal(err)
	}

	m := make([][]string, 0)
	fmt.Println("Enter the search term:")
	var term string
	std_scanner := bufio.NewScanner(os.Stdin)
	if std_scanner.Scan() {
		term = strings.ToLower(std_scanner.Text())
	}
	if len(term) < 4 {
		log.Fatal("Search term too short")
		return
	}
	f_no := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "judgements") {
			f_no = f_no + 1
			file, _ := os.Open("./data/" + f.Name())
			scanner := bufio.NewScanner(file)
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			first_line := true
			for scanner.Scan() {
				line := scanner.Text()
				if first_line {
					if f_no == 1 {
						fmt.Println(strings.Replace(line, ";", "\t", -1))
					}
					first_line = false
					continue
				}
				row := strings.Split(line, ";")
				if strings.Contains(strings.ToLower(line), term) || strings.Contains(strings.ReplaceAll(strings.ToLower(line), " ", ""), term) {
					there := false
					for _, e := range m {
						//comparing against case number
						if strings.Compare(e[2], row[2]) == 0 {
							there = true
						}
					}
					if !there {
						m = append(m, row)
					}
				}
			}

		}
	}
	if len(m) == 0 {
		fmt.Printf("Not found '%s'\n", term)
		return
	}
	for _, e := range m {
		fmt.Println(strings.Join(e, "\t"))
	}
}
