package page

import (
    "testing"
)

func TestFile(t *testing.T) {
    link := Link{
        Url: "https://download.docker.com/",
    }

    filename := GetFileName(link).Name

    t.Log(filename)

    if filename != Prefix+"/index.html" {
        t.Error("not equal")
    }

}

func TestFile2(t *testing.T) {
    link := Link{
        Url: "https://download.docker.com/mac",
    }

    filename := GetFileName(link)

    t.Log(filename)
}

func TestIsDir(t *testing.T) {
    str := ".tar.gz.9"
    println(str[len(str)-1])
    println(string(str[len(str)-1]))

    if string(str[len(str)-1]) > "a" {
        println("true > a")
        t.Fatal("9 should not >  a")
    } else {
        println("false > a")
    }

}
