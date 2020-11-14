//
// main.go
// Copyright (C) 2020 forseason <me@forseason.vip>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var baseUrl string

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("wrong argument")
	}
	baseUrl = "https://kakuyomu.jp"
	novelUrl := os.Args[1]
	chapterList := getChapterList(novelUrl)
	getNovelContent(chapterList)
}

func getChapterList(novelUrl string) []([]string) {
	client := &http.Client{}
	resp, err := client.Get(novelUrl)
	if err != nil {
		log.Fatalf("failed to reach: %s", novelUrl)
	}
	html, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r, _ := regexp.Compile(`<a href=\"(.*?)\" class=\"widget-toc-episode-episodeTitle\">[\s\S]*?<span class=\"widget-toc-episode-titleLabel js-vertical-composition-item\">(.*?)</span>[\s\S]*?</a>`)
	match := r.FindAllStringSubmatch(string(html), -1)
	res := make([]([]string), len(match))
	for k, v := range match {
		res[k] = []string{v[2], baseUrl + v[1]}
	}
	return res
}

func getNovelContent(chapterList []([]string)) {
	ch := make(chan int, 5)
	r, _ := regexp.Compile(`<p id="p\d+">(.*?)</p>`)
	for _, v := range chapterList {
		ch <- 1
		go func(v []string, ch chan int, r *regexp.Regexp) {
			client := &http.Client{}
			resp, err := client.Get(v[1])
			if err != nil {
				log.Fatalf("failed to reach: %s", v[1])
			}
			html, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			match := r.FindAllStringSubmatch(string(html), -1)
			file, err := os.OpenFile(v[0]+".txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			defer file.Close()
			for _, v := range match {
				io.WriteString(file, v[1]+"\r\n")
			}
			log.Printf("chapter %s download finished.", v[0])
			<-ch
		}(v, ch, r)
	}
}
