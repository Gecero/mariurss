package main

import (
	"flag"
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

/*
Program algorithm:
1. Read the files in the paths given via stdin
2. Parse the files as RSS feeds
3. Sort them chronologically
4. Open the HTML file specified as first command line parameter (MUST be UTF-8 encoded!)
5. Write chronological RSS report to #mariurss-content
6. Write current time to #mariurss-time
7. Save changes to HTML file
*/

type News struct {
	item *gofeed.Item
	feed *gofeed.Feed
}

func main() {
	htmlPath := flag.String("html", "", "The .html file to put the rss feeds into")
	flag.Parse()
	if len(*htmlPath) == 0 {
		panic("You must specify an html file path via '-html=<my file path>'")
	}

	feeds := []*gofeed.Feed{}
	files, err := readLines()
	if err != nil {
		panic(err)
	}
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

	// All feeds are read now. Next step is creating a chronologic list
	news := []News{}
	for _, feed := range feeds {
		for _, itm := range feed.Items {
			news = append(news, News{itm, feed})
		}
	}

	sort.Slice(news, func(i, j int) bool {
		timeI := news[i].item.PublishedParsed
		timeJ := news[j].item.PublishedParsed
		return timeI.Compare(*timeJ) == 1
	})

	htmlFile, err := os.OpenFile(*htmlPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(fmt.Sprint("Can't open html file:", err))
	}
	root, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		panic(fmt.Sprint("Can't parse html file:", err))
	}

	timeNode := root.Find("#mariurss-time")
	if timeNode.Length() >= 1 {
		timeNode.Empty()
		timeNode.SetText(time.Now().Format(time.RFC822))
	}

	contentNode := root.Find("#mariurss-content")
	if contentNode.Length() >= 1 {
		contentNode.Empty()
		// Table columns: Title (Link) and Date/Feed, Description
		feedHTML := "<table><tr><th>News</th><th>Description</th></tr>"
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
				length := 350
				if len(contentText) < length {
					length = len(contentText)
				}
				contentText = contentText[0:length] + "..."
			}

			feedHTML += "<tr>"

			feedHTML += "<td class='mariurss-content-main'>"
			feedHTML += "<a class='mariurss-content-title' href='" + new.item.Link + "'>"
			feedHTML += new.item.Title + "</a><br>"
			feedHTML += "<span class='mariurss-content-date'>(" + new.item.PublishedParsed.Format(time.RFC1123) + ")</span>"
			feedHTML += "<br><a class='mariurss-content-feed' href='" + new.feed.Link + "'>"
			feedHTML += new.feed.Title + "</a></td>"

			feedHTML += "<td class='mariurss-content-description'><p>" + contentText + "</p></td>"

			feedHTML += "</tr>"
		}
		feedHTML += "</table>"
		contentNode.AppendHtml(feedHTML)
	}

	if err = htmlFile.Truncate(0); err != nil {
		panic(fmt.Sprint("Can't clear HTML file:", err))
	}
	if _, err := htmlFile.Seek(0, 0); err != nil {
		panic(fmt.Sprint("Can't move writing cursor to beginning of HTML file:", err))
	}

	if err = goquery.Render(htmlFile, root.Selection); err != nil {
		panic(fmt.Sprint("Can't render/write HTML:", err))
	}

	htmlFile.Close()
}
