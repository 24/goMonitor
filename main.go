package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/axgle/mahonia"
	"regexp"
	"strings"
	"time"
	"log"
	"os"
)

const url = "https://www.zuanke8.com/zuixin.php"

func main() {
	argCount := len(os.Args[1:])
	if argCount > 0 {
		cmd()
	} else {
		http.HandleFunc("/zuanke8", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
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

func parseData(data string) []map[string]string {
	reg := regexp.MustCompile(`<a href="https://www.zuanke8.com/thread.*</a>`)
	result := reg.FindAllStringSubmatch(data, 1000)
	retArr := make([]map[string]string, 0)
	for _, item := range result {
		titleReg := regexp.MustCompile(`title.*target`)
		title := titleReg.FindString(item[0])
		urlReg := regexp.MustCompile(`https://.*html`)
		url := urlReg.FindString(item[0])
		title = strings.TrimLeft(title, `title="`)
		title = strings.TrimRight(title, `"  target`)
		match, _ := regexp.MatchString("(速度|水|快|好价|还款)", title)
		if match {
			retItem := make(map[string]string)
			retItem["title"] = title
			retItem["url"] = url
			retArr = append(retArr, retItem)
		}
	}

	return retArr
}

func cmd() {
	for {
		data, statusCode := Get(url)
		if statusCode != 200 {
			fmt.Println(statusCode)
		}
		retArr := parseData(data)
		for _, Item := range retArr {
			fmt.Printf("%s\n%s\n\n", Item["title"], Item["url"])
		}
		println()
		println()
		time.Sleep(5 * time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, statusCode := Get(url)
	if statusCode != 200 {
		fmt.Println(statusCode)
	}
	retArr := parseData(data)
	fmt.Fprintf(w, "%v", `<!DOCTYPE HTML><html> <meta http-equiv="refresh" content="5"><body>`)
	for _, item := range retArr {
		fmt.Fprintf(w, `<a href="%s" target="_blank">%s</a><br> `, item["url"], item["title"])
	}
	fmt.Fprintf(w, "%v", `</body></html>`)
}
