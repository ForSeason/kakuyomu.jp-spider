//
// main.go
// Copyright (C) 2020 forseason <me@forseason.vip>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	baseUrl               string
	novelUrl              string
	novelName             string
	channelNum            int
	ifNumbersInsteadTitle bool
)

func main() {
	parseParams()
	novelName = getNovelTitle()
	createNovelDirectory()
	chapterList := getChapterList()
	getNovelContent(chapterList)
}

func parseParams() {
	flag.IntVar(&channelNum, "j", 5, "numbers of goroutines used to download")
	flag.BoolVar(&ifNumbersInsteadTitle, "n", false, "use number instead of title as filename")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalf("wrong argument")
	}
	baseUrl = "https://kakuyomu.jp"
	novelUrl = flag.Arg(1)
}

func getNovelTitle() string {
	client := &http.Client{}
	resp, err := client.Get(novelUrl)
	if err != nil {
		log.Fatalf("failed to reach: %s", novelUrl)
	}
	html, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r, _ := regexp.Compile(`<h1 id="workTitle"><a href=".*?">(.*?)</a></h1>`)
	match := r.FindAllStringSubmatch(string(html), -1)
	for _, v := range match {
		return v[1]
	}
	return "null"
}

func createNovelDirectory() {
	err := os.Mkdir(novelName, 0777)
	if err != nil {
		log.Fatalf("failed to create directory:%s", err.Error())
	}
}

func getChapterList() []([]string) {
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
	count := 0
	r, _ := regexp.Compile(`<p id="p\d+">(.*?)</p>`)
	for _, v := range chapterList {
		ch <- count
		go func(v []string, ch chan int, r *regexp.Regexp) {
			client := &http.Client{}
			resp, err := client.Get(v[1])
			if err != nil {
				log.Fatalf("failed to reach: %s", v[1])
			}
			html, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			chapterNum := <-ch
			match := r.FindAllStringSubmatch(string(html), -1)
			var fileName string
			if ifNumbersInsteadTitle {
				fileName = strconv.Itoa(chapterNum)
			} else {
				fileName = v[0]
			}
			file, err := os.OpenFile(novelName+"/"+fileName+".txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			defer file.Close()
			for _, v := range match {
				text := fuckRuby(v[1])
				io.WriteString(file, text+"\r\n")
			}
			log.Printf("chapter %s download finished.", v[0])
		}(v, ch, r)
		count++
	}
}

func fuckRuby(text string) string {
	r, _ := regexp.Compile(`<ruby><rb>(.*?)</rb><rp>（</rp><rt>(.*?)</rt><rp>）</rp></ruby>`)
	return r.ReplaceAllString(text, `$1($2)`)
}
