# mariurss

A Unix toolchain that collects RSS feeds and inserts them into an HTML page.
[(Demo)](https://gecero.de/r/)

## Usage
To run this program you will need a bourne-shell compatible shell, GNU coreutils, curl and golang.   
Build by running:
```sh
git clone https://codeberg.org/mdalp/mariurss
cd mariurss/aggregate
go build .
cd ../htmled
go build .
cd ..
```

The tools follow unix philosophy and are bundled together via pipes. For example:   
```
cat my-feed-urls.txt | update/update | aggregate/aggregate | htmled/htmled index.html "#rss-feed"
```

### ``update [store path]``
This tool downloads feeds to disk.   
Give feed URLs to stdin. These will then be downloaded by this tool. The feed URLs may contain comments in the form of non-special-character texts (allowed characters: 'a-zA-z0-9,.- ') behind a ``#``. You may optionally specify a path as first parameter, to which the feeds will be downloaded, if you don't, they will be stored in ``/tmp/mariurss-store/``. Download errors come from stderr. On a successful download, the path of the downloaded feed will be given to stdout. Exit codes are not used.   

### ``aggregate``
This tool aggregates RSS feeds and creates an HTML table respectively.   
Give file paths of RSS/Atom files via stdin. The given HTML file must be in UTF-8 format. On success, stdout will output the HTML code of a table containing the the news feeds entries. The first line is boilerplate (eg &lt;table&gt;), followed by one news entry per line, follwed by a final boilerplate line. On failure, an error log is written to stderr. Exit codes are not used.   
The table will be in chronological order (newest at the top). It uses some CSS classes so you can stylize it to your liking: ``.mariurss-content-main`` (left table column), ``.mariurss-content-description`` (right table column), ``.mariurss-content-feed`` (left column, news feed title), ``.mariurss-content-date`` (left column, news entry date), ``.mariurss-content-title`` (left column, news entry title).

### ``htmled html-file query-selector``
This tool manipulates HTML files.   
The given HTML file must be in UTF-8 format. If more than one element that fits the query selector is found, all of them will be replaced. What you input via stdin will replace the contents of all HTML elements that fit the given [query selector](https://developer.mozilla.org/en-US/docs/Web/API/Document_object_model/Locating_DOM_elements_using_selectors). Exit code is zero on success, non-zero on failure.   
Also, this is the only tool from this toolchain that can be used sensibly outside of this toolchain. For instance, you can also use it to write a timestamp into the HTML file after updating the RSS feed (``date | htmled/htmled index.html "#rss-timestamp"``).
