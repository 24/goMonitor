package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/axgle/mahonia"
	"regexp"
	"strings"
	"time"
)

const url = "https://www.zuanke8.com/zuixin.php"
const delimiter  = "-------------分割线------------"

func main() {
	for {
		data, statusCode := Get(url)
		if statusCode != 200 {
			fmt.Println(statusCode)
		}
		parseData(data)
		println(delimiter)
		time.Sleep(5*time.Second)
	}
}

func Get(url string) (content string, statusCode int) {
	response, err := http.Get(url)
	if err != nil {
		statusCode = -1
		return
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		statusCode = -2
		return
	}
	statusCode = response.StatusCode
	content = ConvertToString(string(data), "gbk", "utf-8")
	return
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func parseData(data string) {
	reg := regexp.MustCompile(`<a href="https://www.zuanke8.com/thread.*</a>`)
	result := reg.FindAllStringSubmatch(data, 1000)
	for _, item := range result {
		titleReg := regexp.MustCompile(`title.*target`)
		title := titleReg.FindString(item[0])
		urlReg := regexp.MustCompile(`https://.*html`)
		url := urlReg.FindString(item[0])
		title = strings.TrimLeft(title, `title="`)
		title = strings.TrimRight(title, `"  target`)
		match, _ := regexp.MatchString("(速度|水|快|好价|还款)", title)
		if match {
			fmt.Println(title)
			fmt.Println(url)
			fmt.Println()
		}
	}
}
