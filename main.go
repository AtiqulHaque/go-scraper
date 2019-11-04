package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) <= 1 {
		log.Printf("not enough arguments provided")
		os.Exit(1)
	}

	url := os.Args[1]
	url = "https://" + url // We "parse" our URL.

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("while creating a new request for %q: %v", url, err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("while making request to %q: %v", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("response status is not %v", http.StatusOK)
		os.Exit(1)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("while parsing the response body from %q: %v", url, err)
		os.Exit(1)
	}

	links := findLinks(nil, doc)

	for _, l := range links {
		fl := formatURL(url, l) // Our formatted link
		req, err := http.NewRequest(http.MethodGet, fl, nil)
		if err != nil {
			log.Printf("while creating a new request for %q: %v", fl, err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("while making request to %q: %v", fl, err)
			continue
		}
		defer resp.Body.Close()

		fmt.Println(resp.StatusCode, fl)
	}
}

func findLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = findLinks(links, c)
	}

	return links
}

func formatURL(base, url string) string {
	base = strings.TrimSuffix(base, "/")

	switch {
	case strings.HasPrefix(url, "/"):
		return base + url
	case strings.HasPrefix(url, "#"):
		return base + "/" + url
	case !strings.HasPrefix(url, "http"):
		return "https://" + url
	default:
		return url
	}
}
