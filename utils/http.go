package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type HttpGet struct {
	URI     string `json:"uri"`
	Proxy   string `json:"proxy"`
	Timeout int    `json:"timeout"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	DebugMode    bool `json:"debug mode"`
	ResponseBody []byte
	StatusCode   int
	Error        error
}

type HttpPost struct {
	URI         string `json:"uri"`
	RequestBody []byte `json:"body"`
	Proxy       string `json:"proxy"`
	Timeout     int    `json:"timeout"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	ResponseBody []byte
	StatusCode   int
	DebugMode    bool `json:"debug mode"`
	Error        error
}

func HttpPostRequest(httpPost *HttpPost, headers map[string]string) error {
	httpClient := &http.Client{}

	if httpPost.Proxy != "" {
		proxyUrl, err := url.Parse(httpPost.Proxy)
		if err == nil {
			httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
			if httpPost.DebugMode {
				log.Printf("Apply Proxy @ %s", proxyUrl)
			}
		} else {
			// log.Printf("%+v", err)
			return err
		}
	}
	if httpPost.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpPost.Timeout) * time.Second
	}
	httpPost.Request, httpPost.Error = http.NewRequest("POST", httpPost.URI, bytes.NewBuffer(httpPost.RequestBody))
	if nil != httpPost.Error {
		return httpPost.Error
	}

	for k, v := range headers {
		httpPost.Request.Header.Add(k, v)
	}

	httpPost.Response, httpPost.Error = httpClient.Do(httpPost.Request)
	if nil != httpPost.Error {
		return httpPost.Error
	}
	defer httpPost.Response.Body.Close()
	httpPost.StatusCode = httpPost.Response.StatusCode
	httpPost.ResponseBody, httpPost.Error = ioutil.ReadAll(httpPost.Response.Body)
	if nil != httpPost.Error {
		return httpPost.Error
	}

	return nil
}

func HttpGetRequest(httpGet *HttpGet) error {
	httpClient := &http.Client{}

	if httpGet.Proxy != "" {
		proxyUrl, err := url.Parse(httpGet.Proxy)
		if err == nil {
			httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
			if httpGet.DebugMode {
				log.Printf("Apply Proxy @ %s", proxyUrl)
			}
		} else {
			// log.Printf("%+v", err)
			return err
		}
	}
	if httpGet.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpGet.Timeout) * time.Second
	}

	httpGet.Request, httpGet.Error = http.NewRequest("GET", httpGet.URI, nil)
	if nil != httpGet.Error {
		return httpGet.Error
	}
	httpGet.Request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	httpGet.Request.Header.Add("Connection", "close")

	httpGet.Response, httpGet.Error = httpClient.Do(httpGet.Request)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	defer httpGet.Response.Body.Close()
	httpGet.StatusCode = httpGet.Response.StatusCode
	httpGet.ResponseBody, httpGet.Error = ioutil.ReadAll(httpGet.Response.Body)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	return nil
}
