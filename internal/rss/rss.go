package rss

import (
	"fmt"
	"html"
	"strings"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (feed *RSSFeed) parseFeed() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for idx, item := range feed.Channel.Item {
		feed.Channel.Item[idx].Description = html.UnescapeString(item.Description)
		feed.Channel.Item[idx].Title = html.UnescapeString(item.Title)
	}
}

func (feed *RSSFeed) PrintFeed() {
	var b strings.Builder
	fmt.Fprintf(&b, "Title: %v\n", feed.Channel.Title)
	fmt.Fprintf(&b, "Link: %v\n", feed.Channel.Link)
	fmt.Fprintf(&b, "Description: %v\n", feed.Channel.Description)
	fmt.Fprintln(&b, "")
	for _, item := range feed.Channel.Item {
		b = item.printItemToBuffer(b)
		fmt.Fprintln(&b, "")
	}
	fmt.Print(b.String())
}

func (item *RSSItem) PrintItem() {
	var b strings.Builder
	fmt.Fprintf(&b, "Title: %v\n", item.Title)
	fmt.Fprintf(&b, "Link: %v\n", item.Link)
	fmt.Fprintf(&b, "Description: %v\n", item.Description)
	fmt.Fprintf(&b, "PubDate: %v\n", item.PubDate)
	fmt.Print(b.String())
}

func (item *RSSItem) printItemToBuffer(b strings.Builder) strings.Builder {
	fmt.Fprintf(&b, "Title: %v\n", item.Title)
	fmt.Fprintf(&b, "Link: %v\n", item.Link)
	fmt.Fprintf(&b, "Description: %v\n", item.Description)
	fmt.Fprintf(&b, "PubDate: %v\n", item.PubDate)
	return b
}
