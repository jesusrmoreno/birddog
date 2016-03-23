package mdl

import "encoding/json"

// TagPostPair ...
type TagPostPair struct {
	Tag     string
	Post    Post
	TagType string
}

// Post is post that actually gets saved to the database
type Post struct {
	SubredditTags []string
	GlobalTags    []string
	Domain        string  `json:"domain"`
	Subreddit     string  `json:"subreddit"`
	Selftext      string  `json:"selftext"`
	ID            string  `json:"id"`
	Gilded        int     `json:"gilded"`
	Author        string  `json:"author"`
	Score         int     `json:"score"`
	Over18        bool    `json:"over_18"`
	NumComments   int     `json:"num_comments"`
	Thumbnail     string  `json:"thumbnail"`
	SubredditID   string  `json:"subreddit_id"`
	Downs         int     `json:"downs"`
	PostHint      string  `json:"post_hint"`
	Permalink     string  `json:"permalink"`
	Name          string  `json:"name"`
	Created       float64 `json:"created"`
	URL           string  `json:"url"`
	Title         string  `json:"title"`
	CreatedUTC    float64 `json:"created_utc"`
	Ups           int     `json:"ups"`
	Priority      int64   `json:"priority"`
	IsSelf        bool    `json:"is_self"`
}

// AsByteSlice returns the post as a []byte
func (p *Post) AsByteSlice() ([]byte, error) {
	var slice []byte
	var err error
	slice, err = json.Marshal(p)
	if err != nil {
		return slice, err
	}
	return slice, nil
}

// PostWrapper ...
type PostWrapper struct {
	Kind string `json:"kind"`
	Data Post   `json:"data"`
}

// Subreddit is the json from the subreddit.
type Subreddit struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string        `json:"modhash"`
		Children []PostWrapper `json:"children"`
	} `json:"data"`
}
