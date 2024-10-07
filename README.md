# Gator CLI blog aggregator

A multi-user command line tool for aggregating RSS feeds and viewing posts.




## Installation

Make sure you have the latest [Go toolchain](https://golang.org/dl/) installed. You will also need a local Postgres database.
You can then install gator with:

```bash
  go install ...
```
## Config

Create a '.gatorconfig.json' file in your home directory. It should utilize the following structure:
```json
{
    "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```
Replace the values with your database connection string.
## Usage/Examples

Create a new user:
```bash
gator register <name>
```

Add a feed:
```bash
gator addfeed <url>
```

Start the aggregator:
```bash
gator agg 30s
```

View your posts:
```bash
gator browse [limit]
```

There are a few other commands to interact with your data.

- 'gator login <name>' - Log in as a registered user. Certain commands require you to be logged in.
- 'gator users' - Lists all registered users.
- 'gator feeds' - Lists all feeds.
- 'gator follow <url>' - Follow a feed that has been added to the database.
- 'gator unfollow <url>' - Unfollow a feed that has been added to the database.