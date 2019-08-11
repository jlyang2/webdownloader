package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var (
	g_c       = colly.NewCollector()
	g_urls    []string
	g_channel = make(chan string)
	g_workMax = 20
	g_wait    = new(sync.WaitGroup)
	g_baseurl string
)

func init() {
	g_c.OnHTML("a,script,link", handleLink)
	g_c.OnResponse(response)

	flag.Parse()
}

func checkErr(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func handleOneURL(url string) {
	dir, file := parseReqUrl(url)
	f, err := os.Open(dir + "/" + file)
	if nil == err {
		f.Close()
	} else {
		err = g_c.Visit(url)
		if nil != err {
			log.Println(url, err)
		}
	}
	g_wait.Done()
}

func inStringList(li []string, s string) bool {
	for _, ss := range li {
		if s == ss {
			return true
		}
	}
	return false
}

func addSubURL() {
	for {
		url := <-g_channel
		if !inStringList(g_urls, url) {
			//fmt.Println(url)
			g_urls = append(g_urls, url)
		}
	}
}

func main() {
	g_baseurl = flag.Arg(0)
	go addSubURL()

	g_urls = append(g_urls, g_baseurl)
	nextIdx := 0
	for {
		workNum := 0
		for workNum < g_workMax {
			if nextIdx == len(g_urls) {
				break
			}
			go handleOneURL(g_urls[nextIdx])
			nextIdx += 1
			workNum += 1
			g_wait.Add(1)
		}
		g_wait.Wait()

		time.Sleep(3 * time.Second)
		if nextIdx >= len(g_urls) {
			break
		}
	}

	log.Printf("done with %d urls\n", len(g_urls))
}

func parseReqUrl(url string) (dir, file string) {
	idx := strings.LastIndex(url, "/")
	if idx+1 == len(url) {
		dir, file = url, "index.html"
	} else {
		dir, file = url[:idx], url[idx+1:]
	}

	if strings.HasPrefix(dir, g_baseurl) {
		dir = "doc/" + strings.TrimPrefix(dir, g_baseurl)
	} else {
		dir = "doc/static"
	}
	return
}

func response(r *colly.Response) {
	dir, file := parseReqUrl(r.Request.URL.String())
	err := os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	checkErr(err)
	err = ioutil.WriteFile(dir+"/"+file, r.Body, os.ModePerm)
	checkErr(err)

	//fmt.Println(string(r.Body))
	//r.Ctx.ForEach(cb)
}

func nextUrl(url, link string) string {
	idx := strings.Index(link, "#")
	if 0 == idx {
		return ""
	}

	if -1 != idx {
		link = link[:idx]
	}

	idx = strings.Index(url, "/")

	return url[:idx+2] + path.Clean(url[idx+2:]+link)
}

func handleLink(e *colly.HTMLElement) {
	var attr string

	switch e.Name {
	case "a":
		fallthrough
	case "link":
		attr = e.Attr("href")
	case "script":
		attr = e.Attr("src")
	}

	if e.Request.URL.String() == g_baseurl {
		log.Println(attr)
	}

	if strings.HasPrefix(attr, "http") {
		g_channel <- attr
		return
	}

	url := nextUrl(e.Request.URL.String(), attr)
	if "" != url {
		g_channel <- url
	}

	return
}

// func cb(k string, v interface{}) interface{} {
// 	fmt.Println(k)
// 	fmt.Println(v)

// 	return nil
// }
