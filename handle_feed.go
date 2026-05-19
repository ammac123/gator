package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ammac123/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v [refresh frequency (e.g. 1s, 1m, 1h25m)]\n", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feedname> <url>", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	ctx := context.Background()

	now := time.Now()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
	})
	if err != nil {
		return fmt.Errorf("error adding feed")
	}
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed")
	}

	fmt.Printf("%+v\n", feed)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("feed to follow does not exist")
	}

	now := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed")
	}

	fmt.Printf("User (%v) is following feed (%v)\n", feedFollow.User, feedFollow.Feed)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	ctx := context.Background()

	followedFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error finding followed feeds")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "User (%v) following %v feeds:\n", user.Name, len(followedFeeds))

	for _, feed := range followedFeeds {
		fmt.Fprintf(w, "\t%c %s\n", '➤', feed.Feed)
	}
	w.Flush()
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feedUrl := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("could not find feed.\n")
	}

	err = s.db.DeleteFeedFollowRecord(ctx, database.DeleteFeedFollowRecordParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed.\n")
	}

	fmt.Printf("Successfully deleted feed '%s'\n", feed.Name)
	return nil

}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}

	var limit int
	if len(cmd.Args) == 1 {
		limit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("could not parse %v as int", limit)
		}
	} else {
		limit = 2
	}

	ctx := context.Background()
	feed, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}

	for _, item := range feed {
		printPost(item)
	}
	return nil

}
