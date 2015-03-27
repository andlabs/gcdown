// 27 march 2015
package main

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"path/filepath"
	"golang.org/x/net/html"
)

const URL = "https://code.google.com/p/bsnes/downloads/list?can=1&q=&colspec=Filename+Summary+Uploaded+ReleaseDate+Size+DownloadCount"

type Download struct {
	URL		string
	Filename	string
}

func collect(url string) (downloads []*Download, err error) {
	var f func(*html.Node)

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
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
				downloads = append(downloads, &Download{
					URL:			"https:" + href,
					Filename:		filepath.Base(href),
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return downloads, nil
}

func get(d *Download) (err error) {
	resp, err := http.Get(d.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.OpenFile(d.Filename, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	downloads, err := collect(URL)
	if err != nil {
		panic(err)
	}
	for _, d := range downloads {
		fmt.Println(d.Filename)
		err = get(d)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Done; %d files downloaded.\n", len(downloads))
}
