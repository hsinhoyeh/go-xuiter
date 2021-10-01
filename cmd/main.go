package main

import (
	"flag"

	goxuiter "github.com/hsinhoyeh/go-xuiter"
	"github.com/hsinhoyeh/go-xuiter/sites"
)

var concurrency = flag.Int64("concurrency", 1, "number of concurrency")
var albumSite = flag.String("album", "", "album site url")
var albumPassword = flag.String("password", "", "album password")
var destinationPrefix = flag.String("destination", "/album", "output file dir")

func main() {
	flag.Parse()

	c := goxuiter.NewCollyController(&goxuiter.NoOpController{}, int(*concurrency))
	alb := sites.NewXuiteAlbumController(c, *destinationPrefix, *albumPassword)
	alb.AddSite(*albumSite)
	c.Run()
}
