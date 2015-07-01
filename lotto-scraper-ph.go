package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type LottoResults struct {
	Results []LottoItem
}

type LottoItem struct {
	Game        string
	Combination string
	Date        string
	Jackpot     string
	Winners     string
}

const CLR_0 = "\x1b[39;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_X = "\x1b[38;1m"
const CLR_N = "\x1b[0m"

func main() {
	wg := &sync.WaitGroup{}
	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("."+string(filepath.Separator)+"results", 0777)
	}
	runtime.GOMAXPROCS(8)
	wg.Add(8)
	go StartScrape("6-55results.asp", CLR_0, "results/6-55.csv", wg)
	go StartScrape("6-49results.asp", CLR_R, "results/6-49.csv", wg)
	go StartScrape("6-45results.asp", CLR_G, "results/6-45.csv", wg)
	go StartScrape("6-42results.asp", CLR_Y, "results/6-42.csv", wg)
	go StartScrape("6-dresults.asp", CLR_B, "results/6-d.csv", wg)
	go StartScrape("4-dresults.asp", CLR_M, "results/4-d.csv", wg)
	go StartScrape("3-dresults.asp", CLR_C, "results/3-d.csv", wg)
	go StartScrape("2-dresults.asp", CLR_W, "results/2-d.csv", wg)
	wg.Wait()
}

func StartScrape(path string, color string, filename string, wg *sync.WaitGroup) {
	fmt.Println(color + "Creating " + filename + CLR_N)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	doc, queryError := goquery.NewDocument("http://pcso-lotto-results-and-statistics.webnatin.com/" + path)
	if queryError != nil {
		fmt.Println("\nError fetching data..\n" + queryError.Error())
		return
	}
	insertedCounter := 0

	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		game := ""
		date := ""
		jackpot := ""
		winners := ""
		combination := ""

		s.Find("td").Each(func(j int, t *goquery.Selection) {
			if j == 0 {
				game = strings.TrimSpace(t.Text())
			}
			if j == 2 {
				date = strings.TrimSpace(t.Text())
			}
			if j == 3 {
				combination = strings.TrimSpace(t.Text())
			}
			if j == 4 {
				jackpot = strings.TrimSpace(t.Text())
			}
			if j == 5 {
				winners = strings.TrimSpace(t.Text())
			}
		})
		f.WriteString("\"" + game + "\",\"" + date + "\",\"" + combination + "\",\"" + jackpot + "\",\"" + winners + "\"\n")
		fmt.Print(color + "●" + CLR_N)
		insertedCounter++
	})
	f.Close()
	fmt.Println(color + "\nCreated csv file :" + filename + " -> " + strconv.Itoa(insertedCounter) + " records " + CLR_N)
	wg.Done()
}
