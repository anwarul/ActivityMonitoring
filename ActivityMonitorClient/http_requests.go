package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var httpClient *http.Client
var lastHttpResponse *http.Response

//-------------------------------
//exported functions
//-------------------------------

func MakeGet(URL string) string {
	resp, err := http.Get(URL)
	lastHttpResponse = resp
	if err == nil {
		return readData(lastHttpResponse.Body)
		lastHttpResponse.Body.Close()
	}
	return err.Error()
}

func MakeGetWithArgs(URL string, args map[string]string) string {
	URL += "?"
	for key, val := range args {
		URL += fmt.Sprintf("%s=%s&", key, val)
	}
	resp, err := httpClient.Get(URL)
	lastHttpResponse = resp
	if err == nil {
		return readData(lastHttpResponse.Body)
		lastHttpResponse.Body.Close()
	}
	return err.Error()
}

func MakeGetWithHeaders(URL string, headers map[string]string) string {
	if httpClient == nil {
		initClient()
	}
	request, _ := http.NewRequest("GET", URL, nil)
	addHeaders(&request.Header, headers)
	resp, err := httpClient.Do(request)
	lastHttpResponse = resp
	defer lastHttpResponse.Body.Close()
	if err == nil {
		return readData(lastHttpResponse.Body)
	}
	return err.Error()
}

func MakePost(URL string, args map[string]string) string {
	fields := url.Values{}
	for key, value := range args {
		fields.Add(key, value)
	}
	resp, err := http.PostForm(URL, fields)
	lastHttpResponse = resp
	defer lastHttpResponse.Body.Close()
	if err == nil {
		return readData(lastHttpResponse.Body)
	}
	return err.Error()
}

func MakePostWithHeaders(URL string, args map[string]string, headers map[string]string) string {
	if httpClient == nil {
		initClient()
	}
	request, _ := http.NewRequest("POST", URL, nil)
	fields := url.Values{}
	for key, value := range args {
		fields.Add(key, value)
	}
	request.PostForm = fields
	addHeaders(&request.Header, headers)
	resp, err := httpClient.Do(request)
	lastHttpResponse = resp
	defer lastHttpResponse.Body.Close()
	if err == nil {
		return readData(lastHttpResponse.Body)
	}
	return err.Error()
}

func GetLastQueryHeader(key string) string {
	if lastHttpResponse != nil && lastHttpResponse.Header != nil {
		return lastHttpResponse.Header.Get(key)
	}
	return ""
}

//-------------------------------
//private functions
//-------------------------------

func readData(source io.ReadCloser) string {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func initClient() {
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
}

func addHeaders(header *http.Header, headers map[string]string) {
	for key, value := range headers {
		header.Add(key, value)
	}
}

func removeHeaders(header *http.Header, headers map[string]string) {
	for key, _ := range headers {
		header.Del(key)
	}
}
