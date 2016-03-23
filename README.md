# Birddog
# Inspired by Buzzfeed's rss-puppy
[Buzzfeed rss-puppy](https://github.com/buzzfeed-openlab/rss-puppy)

```bash
go get github.com/jesusrmoreno/birddog
```

##### A watchdog tool for monitoring subreddits
This library is designed to monitor Subreddits in bulk, and to generate machine friendly
notification when new posts appear and can search for keywords in the title of posts.

This monitor can be run on any cloud service provider, and only requires Go and LedisDB. Also, it is trivial to add output handlers which can pipe subreddit post data to any service you use.

### Set up a database
This monitor uses a ledisdb database to store seen posts and to back the priority queue it uses to determine which subreddits to check next. The priorty of each subreddit is determined by the last time it was checked using a UNIX-Nano timestamp and is prioritized by subreddits that were checked further back in time being popped from the queue.

An example configuration can be found in the example subfolder.

### Configure your subreddits and tags
```toml
# If this is true the main loop will exit after sending the errors to the
# handlers
exitOnError = true

# The user agent that will be presented to reddit
userAgent = "<platform>:<app ID>:<version string> (by /u/<reddit username>)"

# A list of subreddits to subscribe to. This list is from
# http://redditlist.com and represents the subreddits with the most growth
# on Wed Mar 23 2:22 AM
subreddits = [
  "AskReddit",
  "The_Donald",
  "politics",
  "funny",
  "worldnews",
  "pics",
  "AdviceAnimals",
  "videos",
  "nba",
  "todayilearned",
  "gifs",
  ...
]

# All tags should be single words with no punctuation
# Global tags are looked for in every single post regardless of subreddit
globalAlertTags = [
  "globalTag",
]

# Subreddit tags are applied for each post in an individual subreddit
[subredditAlertTags]
  AskReddit = [
    "programming",
    "exampleTag",
  ]

# Reddit says you should not be making more than 60 requests per minute
[throttling]
  concurrentRequests = 10
  monitorFrequency = 5
```

### Configure your outputs
Outputs are modules of code that implement interfaces for events that the monitor executes and do something useful with the resulting data.

There are several different kinds of interfaces:
```Go
// Failer should be implemented by any thing that wants to handle errors
type Failer interface {
	Fail(e error)
}

// TagHandler should be implemeneted by any thing that wants to handl
// all tag notifications.
type TagHandler interface {
	Tag(tp TagPostPair)
}

// PostHandler should be implemented by any thing that wants to handle new posts
type PostHandler interface {
	// Called in goroutine
	Post(p Post)
}

// GlobalTagHandler should be implemented by any thing that wants to handle
// global tag notifications
type GlobalTagHandler interface {
	// Called in goroutine
	GlobalTag(tp TagPostPair)
}

// SubTagHandler should be implemented by any thing that wants to handle
// Subreddit tag notifications
type SubTagHandler interface {
	// Called in goroutine
	SubredditTag(tp TagPostPair)
}
```
An implementation of the TagHandler can be found in the example subfolder
