# RSS Aggregator CLI

A command-line RSS feed aggregator written in Go. This app allows users to register, follow RSS feeds, and browse the latest articles from their followed feeds.

---

## Prerequisites

Before you begin, ensure the following are installed on your system:

- **Go** (v1.20+): https://go.dev/doc/install
- **PostgreSQL**: https://www.postgresql.org/download/

---

## Installation

To install the CLI globally:

```bash
go install github.com/your-username/rss-aggregator@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH` to use the `rss-aggregator` command.

---

## Configuration

The app uses a simple config file stored in your home directory:

```ini
# ~/.rss-aggregator.conf
DBurl=postgres://user:password@localhost:5432/rss?sslmode=disable
Username=
```

- Replace the connection string with your actual PostgreSQL credentials.
- `Username` will be auto-filled once you register/login.

---

## Usage

Run commands like this:

```bash
rss-aggregator <command> [args]
```

### Available Commands

#### Setup & Users

- `register <username>` – Create a new user and log in.
- `login <username>` – Switch to an existing user.
- `users` – List all registered users.
- `reset` – Reset the database (delete all users and orphaned feeds).

#### Feeds

- `feeds` – List all existing feeds.
- `addfeed <name> <url>` – Add a new feed and follow it (logged-in users only).
- `follow <feed-name>` – Follow an existing feed (logged-in users only).
- `unfollow <feed-name>` – Unfollow a feed (logged-in users only).
- `following` – List all feeds the current user follows.

#### Reading Posts

- `browse [limit]` – Show the latest posts for the logged-in user. Optional limit defaults to 2.

#### Aggregation

- `agg <interval>` – Start fetching and storing posts from followed feeds at the given interval (e.g., `10s`, `1m`).

---

## Example

```bash
rss-aggregator register alice
rss-aggregator addfeed golang https://blog.golang.org/feed.atom
rss-aggregator browse 5
```

---

## Notes

- The aggregator will only work properly if the feed URLs are valid and publicly accessible.
- Be sure to set up your Postgres schema correctly (use migrations as needed).
