package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func readLines() ([]string, error) {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range strings.Split(string(stdin), "\n") {
		if len(line) != 0 {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

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

func main() {
	feeds := []*gofeed.Feed{}
	files, err := readLines()
	if err != nil {
		panic(err)
	}
	// Read all files
	for _, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't open file:", err)
			continue
		}
		parser := gofeed.NewParser()
		feed, err := parser.Parse(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not parse feed in ", filepath, ":", err)
			continue
		} else {
			feeds = append(feeds, feed)
		}
	}

	// Create list containting all feed items
	news := []News{}
	for _, feed := range feeds {
		for _, itm := range feed.Items {
			news = append(news, News{itm, feed})
		}
	}

	// Sort by date
	sort.Slice(news, func(i, j int) bool {
		timeI := news[i].item.PublishedParsed
		timeJ := news[j].item.PublishedParsed
		return timeI.Compare(*timeJ) == 1
	})

	// Write HTML code to stdout
	// Table columns: Title (Link) and Date/Feed, Description
	feedHTML := "<table><tr><th>News</th><th>Description</th></tr>\n"
	for _, new := range news {
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

		feedHTML += "<tr>"

		feedHTML += "<td class='mariurss-content-main'>"
		feedHTML += "<a class='mariurss-content-title' href='" + new.item.Link + "'>"
		feedHTML += new.item.Title + "</a><br>"
		feedHTML += "<span class='mariurss-content-date'>(" + new.item.PublishedParsed.Format(time.RFC1123) + ")</span>"
		feedHTML += "<br><a class='mariurss-content-feed' href='" + new.feed.Link + "'>"
		feedHTML += new.feed.Title + "</a></td>"

		feedHTML += "<td class='mariurss-content-description'><p>" + contentText + "</p></td>"

		feedHTML += "</tr>\n"
	}
	feedHTML += "</table>"

	fmt.Println(feedHTML)

}
