package client

import "testing"

func TestGet(t *testing.T) {
    result, _ := Get("https://download.docker.com/linux/centos/7.0/")
    t.Log(string(result.Resp))
    t.Log(result.IsRedirected)
    t.Log(result.FinalUrl)
}
