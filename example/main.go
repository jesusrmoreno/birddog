package main

import (
	"fmt"
	"time"

	"github.com/jesusrmoreno/birddog"
	"github.com/jesusrmoreno/birddog/models"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

type postHandler struct{}

func (p postHandler) Tag(tp mdl.TagPostPair) {
	if tp.Tag == "exampleTag" {
		fmt.Println(tp.Post.ID, tp.Post.Subreddit, tp.Post.Title, tp.Post.Author, tp.Post.URL)
	}
}

func main() {

	// Create our ledisdb database
	cfg := lediscfg.NewConfigDefault()
	l, _ := ledis.Open(cfg)
	db, _ := l.Select(0)

	// Initialize the monitor with the database and config path
	monitor, _ := birddog.New(db, "./config.toml")

	// Initialize our handler
	p := postHandler{}
	// Register our handler
	monitor.RegisterHandler(p)
	// Empty the prior priotity queue
	monitor.PQ.Empty()

	// Star the monitor in its own goroutine
	go monitor.Run()

	time.Sleep(10 * time.Second)

	// After 10 seconds kill the monitor loop
	monitor.Stop()
}
