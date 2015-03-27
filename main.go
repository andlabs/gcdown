// 27 march 2015
package main

import (
	"fmt"
//	"os"
	"net/http"
	"golang.org/x/net/html"
)

const URL = "https://code.google.com/p/bsnes/downloads/list?can=1&q=&colspec=Filename+Summary+Uploaded+ReleaseDate+Size+DownloadCount"

func main() {
	var f func(*html.Node)

	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			href := ""
			is := false
			for _, a := range n.Attr {
				if a.Key == "title" && a.Val == "Download" {
					is = true
				}
				if a.Key == "href" {
					href = a.Val
				}
			}
			if is {
				fmt.Println(href)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
