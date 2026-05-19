package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ammac123/gator/internal/database"
	"github.com/ammac123/gator/internal/rss"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	ctx, close := context.WithTimeout(ctx, 15*time.Second)
	defer close()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt:     now,
	})
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	err = savePostsInFeed(s, ctx, *rssFeed, nextFeed.ID)
	if err != nil {
		return err
	}

	return nil
}

func savePostsInFeed(
	s *state,
	ctx context.Context,
	feed rss.RSSFeed,
	feedID uuid.UUID,
) error {
	for _, item := range feed.Channel.Item {
		err := savePost(s, ctx, item, feedID)
		if err != nil {
			return err
		}
	}
	return nil
}

func savePost(s *state, ctx context.Context, item rss.RSSItem, feedID uuid.UUID) error {
	now := time.Now()

	pubTime, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		return err
	}

	_, err = s.db.CreatePost(ctx, database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Title: sql.NullString{
			String: item.Title, Valid: true,
		},
		Url: item.Link,
		Description: sql.NullString{
			String: item.Description, Valid: true,
		},
		PublishedAt: sql.NullTime{
			Time: pubTime, Valid: true,
		},
		FeedID: feedID,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "(23505)") {
			log.Printf("%v", err)
		}
	}
	return nil
}

func printPost(post database.GetPostsForUserRow) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n(%s)\n", post.FeedName)
	fmt.Fprintf(w, "Title: '%s'\n", post.Title.String)
	fmt.Fprintf(w, "Published: %v\n", post.PublishedAt.Time.Format(time.RFC1123))
	fmt.Fprintf(w, "Post:\n")
	fmt.Fprintf(w, "\t%s\n", post.Description.String)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 25))
	w.Flush()
}
