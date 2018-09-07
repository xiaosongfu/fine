package jpjd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const host = "http://novel.mikepeng.cn"
const startUrl = "/owllook_content?url=http://www.suimeng.la/files/article/html/6/6573/2718566.html&name=%E7%AC%AC%E4%B8%80%E7%AB%A0%20%E5%85%AC%E5%AD%90%20%E5%85%AC%E5%AD%90%EF%BC%881%EF%BC%89&chapter_url=http://www.suimeng.la/files/article/html/6/6573/&novels_name=%E6%9E%81%E5%93%81%E5%AE%B6%E4%B8%81"

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
	title = strings.TrimSpace(title)
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

	if "" == nextPageUrl {
		fmt.Println(">>>>>>>>")
		fmt.Println("已完成")
		fmt.Println(">>>>>>>>")
		return
	}

	//if startIndex > 4 {
	//	fmt.Println(">>>>>>>>")
	//	fmt.Println("已停止")
	//	fmt.Println(">>>>>>>>")
	//	return
	//}

	nextPageUrl = strings.Replace(nextPageUrl, "极品家丁", "%E6%9E%81%E5%93%81%E5%AE%B6%E4%B8%81", -1)
	startIndex++
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
