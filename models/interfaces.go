package mdl

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
