package config

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
)

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}
type RSSFeed struct {
	Channel Channel `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "rss-aggregator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if res.StatusCode > 299 {
		msg := fmt.Sprintf("Response failed with status code: %d\nbody: %s\n", res.StatusCode, body)
		err = errors.New(msg)
		return nil, err
	}
	var result RSSFeed
	err = xml.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cleanResult(&result), nil
}

func cleanResult(feed *RSSFeed) *RSSFeed {
	items := []RSSItem{}
	for _, item := range feed.Channel.Item {
		items = append(items, RSSItem{
			Title:       html.UnescapeString(item.Title),
			Link:        item.Link,
			Description: html.UnescapeString(item.Description),
			PubDate:     item.PubDate,
		})
	}
	cleaned := RSSFeed{
		Channel: Channel{
			Title:       html.UnescapeString(feed.Channel.Title),
			Link:        feed.Channel.Link,
			Description: html.UnescapeString(feed.Channel.Description),
			Item:        items,
		},
	}
	return &cleaned
}
