#!/bin/sh
# This file generates temporary sample files to show how mariurss works.
# It downloads some sample feeds and aggregates the feeds into a sample html file.

echo Note that you must first build the aggregator and htmled! \(See README.md\)

feeds=$(mktemp)
cat > "$feeds"<< EOF
https://yt.cdaut.de/feed/channel/UCNLjRiychUaaCU0uJTYDUYA # yewtu.be is an alternative frontend for YouTube
https://xkcd.com/rss.xml
EOF

index=$(mktemp --suffix=.html)
cat > "$index"<< EOF
<html>
<head>
<style>
tr, td {
	border-style: solid;
	border-width: 1px;
}
table { border-collapse: collapse; }
html {
	font-family: sans-serif;
	max-width: 600px;
	margin: auto;
}
</style>
<title>Demo of mariurss</title>
</head>
<body>
<h1 id="mariurss-time"> PLACEHOLDER </h1>
<div id="mariurss-content"> PLACEHOLDER </div>
</body></html>
EOF

cat "$feeds" | update/update | aggregate/aggregate | htmled/htmled "$index" "#mariurss-content"
date | htmled/htmled "$index" "#mariurss-time"

echo "See sample HTML file at: file://$index"
