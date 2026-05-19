# Gator

A command-line RSS feed aggregator. Follow feeds, run a background aggregator, and browse posts — all from your terminal.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

```sh
go install github.com/ammac123/gator@latest
```

## Configuration

Gator reads its config from `~/.gatorconfig.json`. Create the file with your Postgres connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Create the database before running:

```sh
createdb gator
```

## Usage

### User management

```sh
gator register <name>   # create a new user and set as current
gator login <name>      # switch to an existing user
gator users             # list all users
```

### Feeds

```sh
gator addfeed <name> <url>   # add a new feed and follow it
gator feeds                  # list all feeds
gator follow <url>           # follow an existing feed
gator following              # list feeds you currently follow
gator unfollow <url>         # unfollow a feed
```

### Aggregation

Start the aggregator to fetch new posts on a schedule. It runs continuously until stopped (`Ctrl+C`).

```sh
gator agg 1m     # fetch feeds every minute
gator agg 30s    # fetch feeds every 30 seconds
```

### Browsing posts

```sh
gator browse         # show the 2 most recent posts from your feeds
gator browse 10      # show the 10 most recent posts
```

### Reset

```sh
gator reset   # delete all users and feeds (destructive)
```