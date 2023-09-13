#!/bin/sh
# This file generates temporary sample files to show how mariurss works.
# It downloads two sample feeds from invidious (alternative frontend to YT)
# and aggregates the feeds into a sample html file.

echo Note that you must first build the aggregator and htmled!

feeds=$(mktemp)
cat > "$feeds"<< EOF
https://invidious.snopyta.org/feed/channel/UCNLjRiychUaaCU0uJTYDUYA
https://invidious.snopyta.org/feed/channel/UCCJ-NJtqLQRxuaxHZA9q6zg
EOF

index=$(mktemp --suffix=.html)
cat > "$index"<< EOF
<html><body>
<h1 id="mariurss-time"></h1>
<div id="mariurss-content"></div>
</body></html>
EOF

cat "$feeds"|update/update|aggregate/aggregate|htmled/htmled "$index" "#mariurss-content"
date|htmled/htmled "$index" "#mariurss-time"

echo "See sample html file at: $index"
