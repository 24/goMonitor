package main

import (
	"./src/userUtil"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	argCount := len(os.Args[1:])
	if argCount > 0 {
		cmd()
	} else {
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(":8000", nil))
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
		match, _ := regexp.MatchString("(速度|水|快|好价|还款|话费)", item.Subject)
		if match {
			retItem := make(map[string]string)
			retItem["title"] = item.Subject
			retItem["zuanke8url"] = "https://www.zuanke8.com/thread-" + item.Tid + "-1-1.html"
			retArr = append(retArr, retItem)
		}
	}

	return retArr
}

func outZzuanke8(zk8ch chan<- string) {
	data, statusCode := zuanke8RawData()
	if statusCode != 200 {
		fmt.Println(statusCode)
	}
	retArr := parseZuanke8Data(data)
	zk8str := ""
	for _, Item := range retArr {
		zk8str += fmt.Sprintf("%s\n%s\n\n", Item["title"], Item["zuanke8url"])
	}
	zk8ch <- zk8str
	close(zk8ch)
}

func smzdmRawData() (content string, statusCode int) {
	const url = "https://api.smzdm.com/v1/ranking_list/articles?category_ids=&f=iphone&mall_ids=&offset=0&order=12&slot=12&tab=1&tab_id=47&tag_ids=&v=8.2&weixin=1"
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

func parseZdmData(data string) []map[string]string {
	type Rows []struct {
		Article_title           string `json:"article_title"`
		Article_price           string `json:"article_price"`
		Article_worthy          string `json:"article_worthy"`
		Article_unworthy        string `json:"article_unworthy"`
		Article_worthy_per_cent string `json:"article_worthy_per_cent"`
		Article_url             string `json:"article_url"`
		Article_date            string `json:"article_date"`
	}
	type Data struct {
		Rows Rows
	}
	type Result struct {
		Data Data `json:"data"`
	}
	var result Result
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Fatal("json unmarsha1 failed", err)
	}
	retArr := make([]map[string]string, 0)
	for _, item := range result.Data.Rows {
		retItem := make(map[string]string)
		retItem["title"] = item.Article_title
		retItem["price"] = item.Article_price
		retItem["worthy"] = item.Article_worthy
		retItem["unworthy"] = item.Article_unworthy
		retItem["percent"] = item.Article_worthy_per_cent
		retItem["date"] = item.Article_date
		retItem["url"] = item.Article_url
		retArr = append(retArr, retItem)
	}

	return retArr
}

func outZdm(zdmch chan<- string) {
	zdmData, zdmStatusCode := smzdmRawData()
	if zdmStatusCode != 200 {
		fmt.Println(zdmStatusCode)
	}
	zdmRetArr := parseZdmData(zdmData)
	zdmstr := ""
	for _, Item := range zdmRetArr {
		percent, _ := strconv.ParseFloat(Item["percent"], 64)
		interval := time.Now().Unix() - userUtil.Datetime2timeStamp(Item["date"])
		if percent > 0.70 && interval < 5*60*60 {
			zdmstr += fmt.Sprintf("%s\n%s\n值%s--不值%s\n%.2f%%\n%s\n\n", Item["title"], Item["price"], Item["worthy"], Item["unworthy"], 100*percent, Item["url"])
		}
	}
	zdmch <- zdmstr
	close(zdmch)
}

func cmd() {
	for {
		zk8ch := make(chan string, 1)
		zdmch := make(chan string, 1)
		go outZdm(zdmch)
		go outZzuanke8(zk8ch)
		println(<-zdmch)
		println(<-zk8ch)
		time.Sleep(5 * time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", `<!DOCTYPE HTML><html> <meta http-equiv="refresh" content="5"><body>`)
	zuanke8ch := make(chan string)
	smzdmch := make(chan string)
	go htmlZuanke8(zuanke8ch)
	go htmlZdm(smzdmch)
	fmt.Fprintf(w, "%v", `</body></html>`)
	fmt.Fprintf(w, <-smzdmch)
	fmt.Fprintf(w, <-zuanke8ch)
}

func htmlZuanke8(zuanke8ch chan<- string) {
	data, statusCode := zuanke8RawData()
	if statusCode != 200 {
		fmt.Println(statusCode)
	}
	retArr := parseZuanke8Data(data)
	zuanke8str := ""
	for _, item := range retArr {
		zuanke8str += fmt.Sprintf(`<a href="%s" target="_blank">%s</a><br> `, item["zuanke8url"], item["title"])
	}
	zuanke8ch <- zuanke8str
	close(zuanke8ch)
}

func htmlZdm(smzdmch chan<- string) {
	zdmData, zdmStatusCode := smzdmRawData()
	if zdmStatusCode != 200 {
		fmt.Println(zdmStatusCode)
	}
	zdmRetArr := parseZdmData(zdmData)
	smzdmstr := ""
	for _, zdmItem := range zdmRetArr {
		percent, _ := strconv.ParseFloat(zdmItem["percent"], 64)
		interval := time.Now().Unix() - userUtil.Datetime2timeStamp(zdmItem["date"])
		if percent > 0.70 && interval < 1*60*60 {
			smzdmstr += fmt.Sprintf(`%s<br>`, zdmItem["title"])
			smzdmstr += fmt.Sprintf(`%s<br>`, zdmItem["price"])
			smzdmstr += fmt.Sprintf(`值%s---不值%s<br>`, zdmItem["worthy"], zdmItem["unworthy"])
			smzdmstr += fmt.Sprintf(`%.2f%%<br>`, 100*percent)
			smzdmstr += fmt.Sprintf(`<a href="%s" target="_blank">%s</a><br><br>`, zdmItem["url"], zdmItem["url"])
			println()
		}
	}
	smzdmch <- smzdmstr
	close(smzdmch)
}
