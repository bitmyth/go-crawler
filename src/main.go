package main

import (
    "flag"
    "fmt"
    "github.com/bitmyth/go-crawler/src/page"
    "sync"
)

var depth int

var maxConnections *int

func main() {
    urls := make(chan page.Link, 1)
    pages := make(chan page.Page, 1)

    entry := flag.String("url", "http://image.baidu.com/", "index page to crawl")
    flag.IntVar(&depth, "depth", 1, "depth to search,zero means no limit")
    maxConnections = flag.Int("c", 1, "max connections")

    flag.Parse()

    if entry == nil {
        println("no entry specified")
        return
    }

    page.Base = *entry
    page.Depth = depth

    fmt.Printf("Begin crawl %s for depth %v use %d connections\n", *entry, depth, *maxConnections)

    var chans = make([]chan page.Crawler, *maxConnections)

    for i := 0; i < *maxConnections; i++ {
        crawler := page.New(i)
        chans[i] = crawler.GoCrawl(urls, pages)
    }

    if depth > 0 {
        go page.Parse(pages, urls)
    }

    urls <- page.Link{Depth: 0, Url: *entry, Id: 0}

    for n := range merge(chans[:]) {
        println("Crawler ID:", n.Id, "Success URL count:", n.SuccessUrlCount, "Failed URL count:", n.FailedUrlCount)
    }

    println("exit succeed")

    //// Graceful shutdown
    //quit := make(chan os.Signal, 1)
    //signal.Notify(quit, os.Interrupt)
    //<-quit
    //println(page.UrlCount, len(page.Urls))
}

func merge(cs []chan page.Crawler) chan page.Crawler {
    var wg sync.WaitGroup
    wg.Add(len(cs))

    out := make(chan page.Crawler)

    for _, c := range cs {
        go func(ch chan page.Crawler) {
            var t page.Crawler
            for n := range ch {
                t = n
                out <- n
            }
            println("wg.Down ", "-------------Crawler ID:", t.Id, "Success URL count:", t.SuccessUrlCount, "Failed URL count:", t.FailedUrlCount)
            wg.Done()
        }(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
