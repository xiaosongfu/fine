package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const host = "http://novel.mikepeng.cn"
const startUrl = "/owllook_content?url=http://www.suimeng.la/files/article/html/6/6573/2719256.html&name=%E7%AC%AC%E5%85%AD%E4%B9%9D%E4%B8%80%E7%AB%A0%20%E4%B8%87%E5%8D%83%E6%9F%94%E6%83%85%EF%BC%88%E5%85%A8%E5%89%A7%E7%BB%88%EF%BC%81%EF%BC%89&chapter_url=http://www.suimeng.la/files/article/html/6/6573/&novels_name=%E6%9E%81%E5%93%81%E5%AE%B6%E4%B8%81"

var startIndex = 1

const summaryFileName = "SUMMARY.md"

func main() {
	handlerFunction(host + startUrl)
}

func handlerFunction(webUrl string) {
	fmt.Println("**********************")
	fmt.Println("url: " + webUrl)

	// 请求网页内容
	res, err := http.Get(webUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// 装载网页内容
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 获取标题
	title := doc.Find("#content_name").First().Text()
	fmt.Println("title: " + title)

	// 获取内容
	content, err := doc.Find("#ccontent").First().Html()
	if err != nil {
		panic(err)
	}
	content = strings.Replace(content, "<br/>", "\n", -1)

	contentSplit1 := strings.Split(content, "\n")
	var contentSplit2 []string
	for _, line := range contentSplit1 {
		contentSplit2 = append(contentSplit2, strings.TrimSpace(line))
	}

	content = strings.Join(contentSplit2, "\n")
	//fmt.Println(content)

	// 开始写入文件
	var fileName = strconv.Itoa(startIndex) + ".md"
	_, err = writeSummary(title, fileName) // 写入目录
	if err != nil {
		panic(err)
	}
	_, err = writeContent(fileName, content) // 写入内容
	if err != nil {
		panic(err)
	}

	// 获取下一页链接
	var nextPageUrl string
	var exist bool
	doc.Find(".btn-default").Each(func(i int, selection *goquery.Selection) {
		if strings.TrimSpace(selection.Text()) == "下一页" {
			nextPageUrl, exist = selection.Attr("href")
			if !exist {
				panic(err)
			}
		}
	})
	//fmt.Println(nextPageUrl)

	fmt.Println("**********************")
	fmt.Println("")

	time.Sleep(3000)

	handlerFunction(host + nextPageUrl)
}

func writeContent(fileName, content string) (len int, err error) {
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return
	}

	return file.WriteString(content)
}

func writeSummary(title, filePath string) (len int, err error) {
	file, err := os.OpenFile(summaryFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}

	line := "* [" + title + "](" + filePath + ")\n"

	return file.WriteString(line)
}
