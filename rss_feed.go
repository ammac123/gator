package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ammac123/gator/internal/database"
	"github.com/ammac123/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v [url]\n", cmd.Name)
	}

	// feedURL := html.EscapeString(cmd.Args[0])
	feedURL := html.EscapeString("https://www.wagslane.dev/index.xml")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	feed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	data, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	return nil

}

func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feedname> <url>", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	userName := s.cfg.CurrentUserName
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, userName)
	if (err != nil) || (user.Name != userName) {
		return fmt.Errorf("error retrieving user")
	}

	now := time.Now()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed")
	}

	fmt.Printf("%+v", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	ctx := context.Background()

	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving feeds")
	}
	widths := func(rows []database.GetAllFeedsRow) []string {
		widthOut := make([]int, 3)
		for _, r := range rows {
			if len(r.Name) > widthOut[0] {
				widthOut[0] = len(r.Name)
			}
			if len(r.Url) > widthOut[1] {
				widthOut[1] = len(r.Url)
			}
			if len(r.User) > widthOut[2] {
				widthOut[2] = len(r.User)
			}
		}
		widths := []string{
			strings.Repeat("-", widthOut[0]),
			strings.Repeat("-", widthOut[1]),
			strings.Repeat("-", widthOut[2]),
		}
		return widths
	}(feeds)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "TITLE\tURL\tUSER\n")
	fmt.Fprintf(w, "%s\t%s\t%s\n", widths[0], widths[1], widths[2])
	for _, item := range feeds {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.Url, item.User)
	}
	w.Flush()
	return nil

}
