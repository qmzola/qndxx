package main

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Resp(StatusCode int, Body string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: int(StatusCode),
		Headers:	map[string]string{"Content-Type": "application/json"},
		Body: 		Body,
		IsBase64Encoded: false,
	}, nil
}

func httpHandler(method, urlVal string) (bool, string) {
	client := &http.Client{}
	var req *http.Request

	req, _ = http.NewRequest(method, urlVal, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	req.Header.Set("Referer", "http://news.cyol.com/")
	req.Header.Set("Host", "m.cyol.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	
	resp, err := client.Do(req)
 
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return true, string(b)
}

func regHandler(source, pattern_str string) (bool, string) {
	pattern := regexp.MustCompile(pattern_str)
	result := pattern.FindAllStringSubmatch(source, -1)
	if len(result) >= 1 {
		return true, result[0][1]
	} else {
		return false, ""
	}
}

func requestAndParse(method, url, pattern_str string) (bool, string) {
	if ok, respHtml := httpHandler(method, url); ok {
		if ok, res := regHandler(respHtml, pattern_str); ok {
			return true, res
		}
	}
	return false, ""
}


func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if ok, newest_url := requestAndParse("GET", "http://news.cyol.com/gb/channels/vrGlAKDl/index.html", `<li><a href=\"(.*?\.html)\"`); ok {
		if ok, title := requestAndParse("GET", newest_url, `<h1>(.*?)</h1>`); ok {
			return Resp(200, "{\"url\": \"" + newest_url + "\", \"title\": \"" + title + "\"}")
		}
		return Resp(403, "{\"detail\": \"title解析失败\"}")
	}
	return Resp(403, "{\"detail\": \"url解析失败\"}")
}

func main() {
	lambda.Start(handler)
}
