package page

import (
    "sync"
    "testing"
)

func TestParse(t *testing.T) {

    html := `<a href="centos/">centos/</a>`

    Base = "http://test.com/b"

    p := &Page{
        Depth:   1,
        Content: html,
        Url:     Base,
    }

    r := p.ExtractUrl()

    for i, u := range r {
        println(i, "page:", p.Url, "find url:", u.Url)
    }
}

func TestParse2(t *testing.T) {

    html := `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Index of linux/ubuntu/dists/</title>
</head>
<body>
<h1>Index of linux/ubuntu/dists/</h1>
<hr>
<pre><a href="../">../</a>
<a href="artful/">artful/</a>
<a href="bionic/">bionic/</a>
</pre><hr></body></html>
`

    Base = "https://download.docker.com/linux/ubuntu/dists/"

    p := &Page{
        Depth:   1,
        Content: html,
        Url:     Base,
    }

    r := p.ExtractUrl()

    for i, u := range r {
        println(i, "page:", p.Url, "find url:", u.Url)
    }
}
func C(id int) chan int {
    ch := make(chan int)
    go func() {
        println("send----", id)
        ch <- id
        close(ch)
        println("send", id)
    }()
    return ch
}

func Test3(t *testing.T) {
    chans := make([]<-chan int, 3)

    for i := 0; i < 3; i++ {
        chans[i] = C(i)
    }

    var wg sync.WaitGroup
    wg.Add(len(chans))

    out := make(chan int)

    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }

    for _, ch := range chans {
        go output(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()

    for i := range out {
        println("get", i)
    }
    println("===== over =====")
}

func Test4(t *testing.T) {
    chans := make([]chan int, 3)

    for i := 0; i < 3; i++ {
        chans[i] = C(i)
    }

    var wg sync.WaitGroup
    wg.Add(len(chans))

    out := make(chan int)

    for _, ch := range chans {
        go func(ch chan int) {
            for n := range ch {
                out <- n
            }
            wg.Done()
        }(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()

    for i := range out {
        println("get", i)
    }
    println("===== over =====")
}
