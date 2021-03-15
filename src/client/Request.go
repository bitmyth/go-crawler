package client

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

type Request struct {
    Resp         []byte
    IsRedirected bool
    FinalUrl     string
}

func New() *Request {
    return &Request{
        IsRedirected: false,
    }
}

func (r *Request) HandleRedirect(req *http.Request, via []*http.Request) error {
    r.FinalUrl = req.URL.String()
    r.IsRedirected = true
    return nil
}

func (r *Request) Send(method string, url string, data string, headers map[string]string) (*Request, error) {
    fmt.Println(method, url)

    client := &http.Client{
        CheckRedirect: r.HandleRedirect,
    }

    req, err := http.NewRequest(method, url, strings.NewReader(data))

    if err != nil {
        println(err.Error())
        return r, err
    }

    req.Close = true

    req.Header.Set("Accept", "*")

    req.Header.Set("User-Agent", "go")

    for key, value := range headers {
        req.Header.Set(key, value)
    }

    resp, err := client.Do(req)
    if err != nil {
        println(err.Error())
        return r, err
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    fmt.Println(method, url, resp.StatusCode, resp.Request.URL)

    r.Resp = body
    return r, err
}

func Get(url string) (*Request, error) {
    headers := make(map[string]string)
    return New().Send("GET", url, "", headers)
}
