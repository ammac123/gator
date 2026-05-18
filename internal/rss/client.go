package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

func getData(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	data, err := getData(ctx, feedURL)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed := RSSFeed{}
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed.parseFeed()

	return &feed, nil
}
