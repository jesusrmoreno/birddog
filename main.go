package birddog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jesusrmoreno/birddog/Godeps/_workspace/src/github.com/BurntSushi/toml"
	"github.com/jesusrmoreno/birddog/Godeps/_workspace/src/github.com/beefsack/go-rate"
	"github.com/jesusrmoreno/birddog/Godeps/_workspace/src/github.com/jesusrmoreno/gordita"
	"github.com/jesusrmoreno/birddog/Godeps/_workspace/src/github.com/siddontang/ledisdb/ledis"
	"github.com/jesusrmoreno/birddog/models"
)

const (
	dataURL      = "https://www.reddit.com/r/%s/new/.json"
	invalidChars = "'[](){}<>:,,،、-._?\";/\\&@*"
	seenStore    = "seen_______seen______seen%s"
)

func (ctx *Monitor) dispatchError(err error) {
	for _, handler := range ctx.Failers {
		handler.Fail(err)
	}
	if ctx.Config.ExitOnError {
		os.Exit(1)
	}
}

func (ctx *Monitor) dispatchPost(post mdl.Post) {
	for _, handler := range ctx.PostHandlers {
		go handler.Post(post)
	}
}

func (ctx *Monitor) dispatchGlobalTag(tp mdl.TagPostPair) {
	for _, handler := range ctx.GlobalTagHandlers {
		go handler.GlobalTag(tp)
	}
}

func (ctx *Monitor) dispatchTag(tp mdl.TagPostPair) {
	for _, handler := range ctx.AllTagHandlers {
		if tp.TagType == "global" {
			ctx.dispatchGlobalTag(tp)
		} else if tp.TagType == "subreddit" {
			ctx.dispatchSubTag(tp)
		}
		go handler.Tag(tp)
	}
}

func (ctx *Monitor) dispatchSubTag(tp mdl.TagPostPair) {
	for _, handler := range ctx.SubTagHandlers {
		go handler.SubredditTag(tp)
	}
}

func (ctx *Monitor) getSubreddit(sub string) {
	client := &http.Client{}
	u := fmt.Sprintf(dataURL, sub)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		ctx.dispatchError(err)
		return
	}
	req.Header.Set("User-Agent", ctx.Config.UserAgent)
	res, err := client.Do(req)
	if err != nil {
		ctx.dispatchError(err)
		return
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.dispatchError(err)
		return
	}
	s := mdl.Subreddit{}
	if err := json.Unmarshal(contents, &s); err != nil {
		ctx.dispatchError(err)
		return
	}
	normlSub := strings.ToLower(sub)
	globalTags := ctx.Config.GlobalTags
	subTags := ctx.Config.AlertTags[normlSub]
	globalTagMap := map[string]bool{}
	subTagMap := map[string]bool{}
	for _, tag := range globalTags {
		globalTagMap[strings.ToLower(tag)] = true
	}
	for _, tag := range subTags {
		subTagMap[strings.ToLower(tag)] = true
	}
	for i := range s.Data.Children {
		child := s.Data.Children[i]
		postKey := []byte(fmt.Sprintf(seenStore, child.Data.ID))
		storedVal, err := ctx.DB.Get(postKey)
		if err != nil {
			ctx.dispatchError(err)
			return
		}
		if storedVal != nil {
			continue
		}
		finalPost := &child.Data
		normlTitle := stripchars(child.Data.Title, invalidChars)
		wordsInTitle := strings.Split(normlTitle, " ")
		if len(globalTags) > 0 || len(subTags) > 0 {
			for _, word := range wordsInTitle {
				word = strings.ToLower(stripchars(word, " "))
				if globalTagMap[word] == true {
					finalPost.GlobalTags = append(finalPost.GlobalTags, word)
				}
				if subTagMap[word] == true {
					finalPost.SubredditTags = append(finalPost.SubredditTags, word)
				}
			}
		}
		postAsSlice, err := child.Data.AsByteSlice()
		if err != nil {
			ctx.dispatchError(err)
			return
		}
		err = ctx.DB.Set(postKey, postAsSlice)
		if err != nil {
			ctx.dispatchError(err)
			return
		}
		ctx.dispatchPost(child.Data)
		for _, word := range child.Data.GlobalTags {
			ctx.dispatchTag(mdl.TagPostPair{
				Post:    child.Data,
				Tag:     word,
				TagType: "global",
			})
		}
		for _, word := range child.Data.SubredditTags {
			ctx.dispatchTag(mdl.TagPostPair{
				Post:    child.Data,
				Tag:     word,
				TagType: "subreddit",
			})
		}
	}
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

// short ...
func short(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}

// Monitor holds the values we need ...
type Monitor struct {
	ShouldRun         bool
	PQ                *gordita.PQ
	DB                *ledis.DB
	Config            *Config
	Failers           []mdl.Failer
	PostHandlers      []mdl.PostHandler
	GlobalTagHandlers []mdl.GlobalTagHandler
	SubTagHandlers    []mdl.SubTagHandler
	AllTagHandlers    []mdl.TagHandler
}

// Config is the configuration for the app
type Config struct {
	ExitOnError bool
	UserAgent   string
	DBPath      string
	Subreddits  []string
	GlobalTags  []string            `toml:"globalAlertTags"`
	AlertTags   map[string][]string `toml:"subredditAlertTags"`
	Throttling  struct {
		ConcurrentRequests int
		MonitorFrequency   time.Duration
	}
}

// RegisterHandler should be used to register handlers that have implemented
// one or more of the exposed interfaces...
func (ctx *Monitor) RegisterHandler(handler interface{}) {
	var valid bool
	if _, ok := interface{}(handler).(mdl.Failer); ok {
		ctx.Failers = append(ctx.Failers, handler.(mdl.Failer))
		valid = true
	}
	if _, ok := interface{}(handler).(mdl.PostHandler); ok {
		ctx.PostHandlers = append(ctx.PostHandlers, handler.(mdl.PostHandler))
		valid = true
	}
	if _, ok := interface{}(handler).(mdl.TagHandler); ok {
		ctx.AllTagHandlers = append(ctx.AllTagHandlers, handler.(mdl.TagHandler))
		valid = true
	}
	if _, ok := interface{}(handler).(mdl.GlobalTagHandler); ok {
		ctx.GlobalTagHandlers = append(ctx.GlobalTagHandlers, handler.(mdl.GlobalTagHandler))
		valid = true
	}
	if _, ok := interface{}(handler).(mdl.SubTagHandler); ok {
		ctx.SubTagHandlers = append(ctx.SubTagHandlers, handler.(mdl.SubTagHandler))
		valid = true
	}
	if !valid {
		log.Fatal("Invalid handler does not implement any interfaces.")
	}
}

func readConfig(f string) (*Config, error) {
	if _, err := os.Stat(f); err != nil {
		return nil, err
	}
	conf := Config{}
	if _, err := toml.DecodeFile(f, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// New creates a new monitor
func New(db *ledis.DB, configFilePath string) (*Monitor, error) {
	fmt.Println(configFilePath)
	conf, err := readConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	q := gordita.NewPQ(db)
	ctx := Monitor{
		PQ:     q,
		DB:     db,
		Config: conf,
	}
	return &ctx, nil
}

// Run starts off the process of getting things
func (ctx *Monitor) Run() {
	interval := ctx.Config.Throttling.MonitorFrequency
	numRequests := ctx.Config.Throttling.ConcurrentRequests
	rl := rate.New(int(numRequests), time.Second*interval)
	subreddits := ctx.Config.Subreddits
	for _, sub := range subreddits {
		ctx.PQ.Push([]byte(sub), timeScore())
	}
	ctx.ShouldRun = true
	for ctx.ShouldRun {
		if ok, _ := rl.Try(); ok {
			if !ctx.PQ.IsEmpty() {
				sub := ctx.PQ.Pop()
				go ctx.getSubreddit(string(sub.Value))
				ctx.PQ.Push([]byte(sub.Value), timeScore())
			}
		}
	}
}

// Stop stops the monitor
func (ctx *Monitor) Stop() {
	ctx.ShouldRun = false
}

func timeScore() int64 {
	return time.Now().UnixNano()
}
