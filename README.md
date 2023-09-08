# mariurss

A program that collects RSS feeds and inserts them into an HTML page.

Basic usage:
```sh
go build .

cat my-feed-urls.txt | ./mariurss-update.sh | ./mariurss -html=index.html
```
Make sure that your HTML file contains an element with the id "mariurss-content" where the feed can be placed.

For more details, see: https://gecero.de/rss/ (German).