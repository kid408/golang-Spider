package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// 抓取网页源码
func getHtml(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 提取网页title
func getTitle(body *string) (string, error) {
	reg, err := regexp.Compile(`<title>.*?</title>`)
	if err != nil {
		return "", err
	}
	title := string(reg.Find([]byte(*body)))
	return title, nil
}

// 提取文章主体
func getArtical(body *string) (string, error) {
	front := `<div class="article-body">`
	end := `<div class="previous-next-links">`
	reg, err := regexp.Compile(front + `[\s\S]*` + end)
	if err != nil {
		return "", err
	}
	art := string(reg.Find([]byte(*body)))
	art = art[:len(art)-len(end)]
	return art, nil
}

// 按小节构造一页
func makeHtml(title string, body *string) string {
	ht := `<html><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>`
	ht = ht + title
	ht = ht + `<body>` + *body + `</body></html>`
	return ht
}

// 提取网页中下一页的链接
func getNextUrl(body *string) (string, bool) {
	re := `<a href=".*?" rel="next">.*?</a>`
	reg, _ := regexp.Compile(re)
	tmp := string(reg.Find([]byte(*body)))
	if tmp == "" {
		return "", false
	}
	t := strings.Index(tmp, ` rel="next">`)
	return tmp[9 : t-1], true
}

// 主过程
func mainProcess(url string, bookName string) {
	file, err := os.OpenFile(bookName+`.html`, os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("ERROR: Openfile failed, %s", err.Error())
		return
	}
	defer file.Close()
	file.Write([]byte(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/><title>` + bookName + `</title></head><body>` + "\n"))
	num := 0
	for {
		body, err := getHtml(url)
		title, _ := getTitle(&body)
		fmt.Println("Getting:", title)
		if err != nil {
			fmt.Errorf("ERROR: getHtml failed, url = %s, %s", url, err.Error())
			return
		}
		art, err := getArtical(&body)
		if err != nil {
			fmt.Errorf("ERROR: getArtical failed, body = %s, %s", body, err.Error())
			return
		}
		file.Write([]byte(art + "\n"))
		num++
		var haveNext bool = false
		url, haveNext = getNextUrl(&body)
		if haveNext == false {
			break
		}
	}
	file.Write([]byte(`</body></html>`))
	fmt.Printf("-------------End,共抓取了%d页---------------", num)
}

func main() {

	inUrl := `http://www.runoob.com/mongodb/mongodb-tutorial.html`
	bookName := `MongoDB`
	fmt.Println(`---------------start to fetch----------------`)
	mainProcess(inUrl, bookName)
}
