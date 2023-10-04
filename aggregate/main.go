package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func crampString(str string, bounds int) string {
	indicator := "..."
	if len(str) > bounds {
		return str[0:bounds-len(indicator)] + indicator
	} else {
		return str
	}
}

type News struct {
	item *gofeed.Item
	feed *gofeed.Feed
}

func parseNewsFileAsync(path string, news chan News, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
	}
	defer file.Close()

	parser := gofeed.NewParser()
	feed, err := parser.Parse(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing file:", err)
	}

	for _, item := range feed.Items {
		news <- News{item, feed}
	}
}

func collectNewsItems(newsDB *[]News, newsChan chan News) {
	for {
		*newsDB = append(*newsDB, <-newsChan)
	}
}

func readAndParse() []News {
	news := []News{}
	newsChan := make(chan News)
	var wg sync.WaitGroup
	// This'll read from 'newsChan' to 'news' while avoiding race condition
	go collectNewsItems(&news, newsChan)

	// Read line-by-line, parse each file async
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		go parseNewsFileAsync(scanner.Text(), newsChan, &wg)
	}

	wg.Wait()
	return news
}

func sortNews(news *[]News) {
	// Sort by date, newest at the top
	sort.Slice(*news, func(i, j int) bool {
		timeI := (*news)[i].item.PublishedParsed
		timeJ := (*news)[j].item.PublishedParsed
		return timeI.Compare(*timeJ) == 1
	})
}

func writeTable(news *[]News) {
	// Write HTML code to stdout
	// Table columns: Title (Link) and Date/Feed, Description
	// Important: First line is table header, following lines are one
	//            news entry each, last line is table closing tag
	fmt.Println("<table><tr><th>News</th><th>Description</th></tr>")
	for _, new := range *news {
		var line string
		reader := strings.NewReader(new.item.Content)
		contentNode, err := goquery.NewDocumentFromReader(reader)
		var contentText string
		if err != nil {
			contentText = "[Preview failed.]"
		} else {
			contentText = contentNode.Text()
			contentText = strings.Trim(contentText, "\n ")
			contentText = strings.ReplaceAll(contentText, "\n", "&nbsp;")
			contentText = crampString(contentText, 350)
		}

		line += "<tr>"

		line += "<td class='mariurss-content-main'>"
		line += "<a class='mariurss-content-title' href='" + new.item.Link + "'>"
		line += new.item.Title + "</a><br>"
		line += "<span class='mariurss-content-date'>(" + new.item.PublishedParsed.Format(time.RFC1123) + ")</span>"
		line += "<br><a class='mariurss-content-feed' href='" + new.feed.Link + "'>"
		line += new.feed.Title + "</a></td>"

		line += "<td class='mariurss-content-description'><p>" + contentText + "</p></td>"

		line += "</tr>"
		fmt.Println(line)
	}

	fmt.Println("</table>")
}

func main() {
	news := readAndParse()
	sortNews(&news)
	writeTable(&news)
}
