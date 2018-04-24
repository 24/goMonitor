package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"time"
	"log"
	"os"
	"encoding/json"
	"regexp"
	"./src/userUtil"
)

func main() {
	argCount := len(os.Args[1:])
	if argCount > 0 {
		cmd()
	} else {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
	}
}

func zuanke8RawData() (content string, statusCode int) {
	zuanke8url := "http://www.zuanke8.com/api/mobile/index.php?sessionid=&version=4.1&zstamp=" + "1524549400" + "&module=zuixin&sign=" + "16f179de310d945c7f17d91d3e97b8e4"
	req, err := http.NewRequest("POST", zuanke8url, strings.NewReader("mod=all"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		statusCode = -1
		return
	}
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		statusCode = -2
		return
	}
	statusCode = response.StatusCode
	content = userUtil.ConvertToString(string(data), "gbk", "utf-8")
	return content, statusCode
}

func parseZuanke8Data(data string) []map[string]string {
	type Relist []struct {
		Subject string `json:"subject"`
		Tid     string `json:"tid"`
	}
	type Data struct {
		Relist Relist
	}
	type Result struct {
		Data Data `json:"data"`
	}
	var result Result
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Fatal("json unmarsha1 failed", err)
	}
	retArr := make([]map[string]string, 0)
	for _, item := range result.Data.Relist {
		match, _ := regexp.MatchString("(速度|水|快|好价|还款)", item.Subject)
		if match {
			retItem := make(map[string]string)
			retItem["title"] = item.Subject
			retItem["zuanke8url"] = "https://www.zuanke8.com/thread-" + item.Tid + "-1-1.html"
			retArr = append(retArr, retItem)
		}
	}

	return retArr
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
	content = userUtil.ConvertToString(string(data), "gbk", "utf-8")
	return
}

func cmd() {
	for {
		data, statusCode := zuanke8RawData()
		if statusCode != 200 {
			fmt.Println(statusCode)
		}
		retArr := parseZuanke8Data(data)
		for _, Item := range retArr {
			fmt.Printf("%s\n%s\n\n", Item["title"], Item["zuanke8url"])
		}
		println()
		println()
		time.Sleep(5 * time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, statusCode := zuanke8RawData()
	if statusCode != 200 {
		fmt.Println(statusCode)
	}
	retArr := parseZuanke8Data(data)
	fmt.Fprintf(w, "%v", `<!DOCTYPE HTML><html> <meta http-equiv="refresh" content="5"><body>`)
	for _, item := range retArr {
		fmt.Fprintf(w, `<a href="%s" target="_blank">%s</a><br> `, item["zuanke8url"], item["title"])
	}
	fmt.Fprintf(w, "%v", `</body></html>`)
}
