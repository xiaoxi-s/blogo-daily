package models

import "encoding/xml"

type GoogleNewsEntry struct {
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Title       string `xml:"title"`
}

type GoogleNewsFeed struct {
	XMLName xml.Name `xml:"rss"`

	Items []GoogleNewsEntry `xml:"channel>item"`
}

type Entry struct {
	Link        string
	Description string
	Title       string
}
