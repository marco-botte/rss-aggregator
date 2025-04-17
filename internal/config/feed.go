package config

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"rss-aggregator/internal/database"
	"time"
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
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
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

func ScrapeFeeds(s *State, cmd CommandInput) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error getting next feed to fetch: %v", err)
		return nil
	}
	now := time.Now()
	params := database.MarkFeedFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
	}
	s.Db.MarkFeedFetched(context.Background(), params)
	rss_feed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("Error fetching feed: %v", err)
		return nil
	}
	fmt.Printf("Titles from : %s\n", rss_feed.Channel.Title)
	for _, item := range rss_feed.Channel.Item {
		fmt.Printf("* %s\n", item.Title)

	}
	return nil
}
