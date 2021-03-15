package page

import (
    "fmt"
    "os"
    "regexp"
    "strings"
)

var Base string
var Urls = make(map[string]bool)

// var r, _ = regexp.Compile(`((http://|https://)(www\.)?[^ "';<>\n>]+?\.(\w{2,3})[^ "';<>\n>]+?\.jpg)[ "'>]`)
var r, _ = regexp.Compile(`((http://|https://)(www\.)?[^ "';<>\n>]+?\.(\w{2,3})[^ "';<>\n>]+?)[ "'>]`)

// consume pages, produce urls
func Parse(pages chan Page, links chan Link) {
    storage, _ := os.Create("urls.txt")
    defer storage.Close()

    for page := range pages {

        extractedUrls := page.ExtractUrl()

        for _, link := range extractedUrls {

            if _, present := Urls[link.Url]; present {
                println("!!!!!!! present", link.Url)
                continue
            }

            println("find url =============", link.Id, link.Url, "in page", page.Url)
            if strings.LastIndex(link.Url, ".js") > 0 {
                continue
            }
            if strings.LastIndex(link.Url, ".css") > 0 {
                continue
            }
            if strings.LastIndex(link.Url, ".dtd") > 0 {
                continue
            }
            println("about send link", link.Url)
            links <- link
            println("send link", link.Url)
            Urls[link.Url] = true

            _, _ = storage.WriteString(fmt.Sprintf("%d %d %s", link.Depth, link.Id, link.Url+"\n"))
        }
    }
}
