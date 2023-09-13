package main

import (
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) != 3 {
		panic("Illegal number of arguments. Syntax: htmled [html file] [query selector]")
	}

	filepath := os.Args[1]
	query := os.Args[2]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filepath, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		panic(err)
	}

	// Manipulate HTML
	node := doc.Find(query)
	if node.Length() >= 1 {
		node.Empty()
		node.AppendHtml(string(input))
	} else {
		panic(fmt.Sprintln("Element not found:", query))
	}

	// Clear file
	if err = file.Truncate(0); err != nil {
		panic(err)
	}
	if _, err = file.Seek(0, 0); err != nil {
		panic(err)
	}

	// Write to file
	err = goquery.Render(file, doc.Selection)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
