package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test1(t *testing.T) {
	aa := []string{"a", "c"}

	for _, s := range aa {
		t.Log(s)
	}
}

func Test2(t *testing.T) {
	t.Log(nextUrl("https://docs.ceph.com/docs/master/", "genindex"))
}

func eachElem(idx int, sel *goquery.Selection) {
	s, _ := sel.Attr("href")
	fmt.Println(s)
}

func Test3(t *testing.T) {
	f, _ := os.Open("index.html")
	data, _ := ioutil.ReadAll(f)
	f.Close()
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(data))
	doc.Find("a,link,script").Each(eachElem)
}

func Test4(t *testing.T) {
	d, f := parseReqUrl("https://test/")
	t.Log(d, f)

	d, f = parseReqUrl("https://test/aa")
	t.Log(d, f)
}
