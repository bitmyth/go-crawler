package page

import (
    "net/url"
    "os"
    "path"
    "strings"
)

var Prefix string

func init() {
    Prefix, _ = os.Getwd()
    Prefix += "/downloads/"
}

func IsDir(link string) bool {
    suffix := strings.LastIndex(link, ".")
    if suffix > 0 {
        t := link[suffix:]
        // for some file end with common file types ,not with version numbers like 1.7.3
        if string(t[len(t)-1]) > "a" {
            //is file
            return false
        }
        //is dir
        return true
    }
    return true
}

type FileName struct {
    Name   string
    Suffix string
}

func (f FileName) IsMedia() bool {
    mediaTypes := []string{
        "jpeg",
        "jpg",
        "png",
        "gif",
    }

    suffix := strings.Trim(f.Suffix, ".")

    for _, b := range mediaTypes {
        if b == suffix {
            return true
        }
    }

    return false
}

func GetFileName(link Link) *FileName {
    u, _ := url.Parse(link.Url)

    var dir string

    dir = Prefix + strings.Trim(u.Path, "/")
    if IsDir(u.Path) {
        dir += "/"
    }

    dir = path.Dir(dir)

    err := os.MkdirAll(dir, os.ModePerm)

    if err != nil {
        println(err)
    }

    filename := Prefix + u.Path

    return &FileName{
        filename,
        Suffix(u.Path),
    }
    //filename += Suffix(u.Path)
    //filename := dir + fmt.Sprintf("%v-%d", time.Now().Unix(), link.Id)
    //return filename
}

func Suffix(link string) (suffix string) {
    link = strings.Trim(link, "/")

    parsedUrl, _ := url.Parse(link)

    link = parsedUrl.Path

    if link == "" {
        suffix = "/index.html"
        return
    }

    if i := strings.LastIndex(link, "."); !IsDir(link) && i > 0 {
        suffix = link[i:]
        return
    } else {
        suffix = "/index.html"
        return
    }
    return
}
