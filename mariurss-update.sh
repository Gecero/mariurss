#!/bin/sh
# MARIURSS-UPDATE.SH
# Requirements: Curl, Coreutils, should work in bourne shell
# Downloads/Syncs given remote files and stores them in a 
# mariurss-compatible format.
#
# Usage: This script needs two inputs, one as command-line parameter
#        and one via stdin.
# Command-line parameter: The location of the 'store'. This is the 
#                         directory where the feed data will be stored
# Stdin: The URLs to fetch the data from. One URL per line.
if [ $# -ne 1 ]
then store=/tmp/mariurss-store/
else store=$1
fi
mkdir -p "$store"
cd "$store"

while IFS=$'\n' read -r url; do
    # Remove spaces / turn comments (using '#') into valid urls
    url=$(echo "$url" | sed "s/ //g")
    ID=$(echo "$url" | sha1sum | head -c 40)
    file=$ID

    curl --silent --connect-timeout 60 --output "$file" --time-cond "$file" "$url"
    if [ $? -ne 0 ]; then
        echo "Could not download file #$ID under URL $url" >&2 # to stderr
    fi

    if [ -e "$file" ]; then
        echo "$(pwd)/$file"
    fi
done