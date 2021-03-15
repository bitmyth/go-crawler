package page

import (
    "fmt"
    "github.com/bitmyth/go-crawler/src/client"
    "io/ioutil"
    "math"
    "os"
    "regexp"
    "strings"
)

type Page struct {
    Depth   int
    Content string
    Url     string
}

type Link struct {
    Depth int
    Url   string
    Id    int
}

var Depth int

var UrlCount int

type Crawler struct {
    SuccessUrlCount int
    FailedUrlCount  int
    Done            chan Crawler
    Id              int
}

func New(id int) *Crawler {
    return &Crawler{
        SuccessUrlCount: 0,
        FailedUrlCount:  0,
        Done:            make(chan Crawler),
        Id:              id,
    }
}

func (c Crawler) GoCrawl(links chan Link, pages chan Page) chan Crawler {
    go c.Crawl(links, pages)

    return c.Done
}

// consume urls, produce pages
func (c Crawler) Crawl(links chan Link, pages chan Page) {
    for link := range links {

        if link.Depth == Depth {
            fmt.Printf("%d link Depth %d exceed limit\n", c.Id, link.Depth)
            fmt.Printf("%d total urls: %d\n", c.Id, c.SuccessUrlCount)
            c.Done <- c
            close(c.Done)
            return
        }

        println(c.Id, "downloading", link.Id, link.Url)

        resp, err := client.Get(link.Url)

        if err != nil {
            println(err.Error())
            c.FailedUrlCount++
            continue
        }

        c.SuccessUrlCount++

        content := resp.Resp

        if string(content) != "" {
            println(c.Id, "downloaded size:", len(content))

            //pages <- Page{link.Depth, string(content)}

            filename := GetFileName(link)
            println("file", filename.Name, "suffix:", filename.Suffix)

            if len(content) > 10*1024 && filename.IsMedia() {
                err1 := ioutil.WriteFile(filename.Name, []byte(string(content)+link.Url), os.ModePerm)
                if err1 != nil {
                    println(c.Id, "error:", err1.Error())
                    continue
                } else {
                    println(c.Id, "save to file", filename, len(content))
                }
            }

            if filename.Suffix == "/index.html" {
                if resp.IsRedirected {
                    link.Url = resp.FinalUrl
                }
                println(c.Id, "about send page", link.Url)

                go func(link Link, content []byte) {
                    pages <- Page{link.Depth, string(content), link.Url}
                    println("send page ----->  ", link.Url)
                }(link, content)
            }
        }
    }
}

func (page Page) ExtractUrl() []Link {
    matches := r.FindAllString(page.Content, math.MaxInt64)

    var urls []Link

    for _, absoluteUrl := range matches {

        urlTrimmed := strings.Trim(absoluteUrl, "//\"'")

        fmt.Println("page Depth:", page.Depth, "find url:", UrlCount, urlTrimmed)

        urls = append(urls, Link{page.Depth + 1, urlTrimmed, UrlCount})

        UrlCount++
    }

    var r = regexp.MustCompile(`(href|src) *= *("|') *(?P<url>[^ "';<>\n>]+)`)
    result := r.FindAllSubmatch([]byte(page.Content), -1)

    for _, relativeUrl := range result {
        currentUrl := formatRelativeUrl(string(relativeUrl[3]), page.Url)
        if currentUrl != page.Url {
            urls = append(urls, Link{page.Depth + 1, currentUrl, UrlCount})
            UrlCount++
        }
    }

    return urls
}

func formatRelativeUrl(value string, baseUrl string) string {
    currentUrl := strings.Trim(value, " //\"'.")
    baseUrl = strings.Trim(baseUrl, "/")

    if strings.Index(value, "http") == 0 {
        return currentUrl
    }

    //u, err := url.Parse(baseUrl)
    //if err != nil {
    //    panic(err)
    //}
    //fmt.Printf("baseUrl: %v\n", baseUrl)
    //fmt.Printf("u.Path: %v\n", u.Path)

    //currentUrl = strings.Join([]string{u.Scheme + "://" + u.Host + path, currentUrl}, "/")
    currentUrl = baseUrl + "/" + currentUrl

    return currentUrl
}
