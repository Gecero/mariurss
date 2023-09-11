# mariurss

A Unix program that collects RSS feeds and inserts them into an HTML page.
[(Demo)](https://gecero.de/rss/)

## Usage
To run this program you will need a bourne-shell compatible shell, GNU coreutils, curl and golang.   
Build by running:
```sh
git clone https://codeberg.org/mdalp/mariurss
cd mariurss
go build .
```

The two tools, mariurss-update and mariurss are designed to work together, not standalone.

You can them like this:
```
cat my-feed-urls.txt | ./mariurss-update.sh | ./mariurss -html=index.html
```

### mariurss-update
The tool that downloads the feeds to disk.   
Give feed URLs to stdin. These will then be downloaded by this tool. The feed URLs may contain comments in the form of non-special-character texts behind a ``#``. You may optionally specify a path as first parameter, to which the feeds will be downloaded, if you don't, they will be stored in ``/tmp/mariurss-store/``. Download errors come in stderr. On a successful download, the path of the downloaded feed will be given to stdout. Exit codes are not used.   

### mariurss
The tool that modifies an HTML file to contain the latest feed information.   
Give file paths of RSS/Atom files via stdin. Ideally, you pipe the output of mariurss-update into mariurss. As command line parameter, you must specify the path of the HTML file to manipulate, by writing ``-html=<path of html file>``. The given HTML file must be in UTF-8 format. Stdout will output nothing on success or may contain error messages. Exit codes are not used.   
The tool will look for the HTML element ``#mariurss-content``, clear it, and insert a table featuring all feed news entries found in all files specified in stdin, sorted chronologically (latest at top, oldest at bottom). It will also look for the HTML element ``#mariurss-time``, clear it and insert a timestamp.   
The HTML elements inserted by mariurss have class labels so you can stylize them using CSS. These classes are: ``.mariurss-content-main`` (left table column), ``.mariurss-content-description`` (right table column), ``.mariurss-content-feed`` (left column, news feed title), ``.mariurss-content-date`` (left column, news entry date), ``.mariurss-content-title`` (left column, news entry title).   

