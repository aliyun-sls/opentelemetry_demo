package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type HttpClient struct {
	Host string
}

func (m *HttpClient) Post(ctx context.Context, path string, data interface{}, res *Result) (err error) {
	m.Host = "http://" + m.Host
	b, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, m.Host+path, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = request(ctx, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (m *HttpClient) Get(ctx context.Context, path string, dataMap map[string]interface{}, res *Result) (err error) {
	m.Host = "http://" + m.Host
	values := url.Values{}
	for k, v := range dataMap {
		switch str := v.(type) {
		case string:
			values.Add(k, str)
		case int:
			values.Add(k, strconv.Itoa(str))
		case uint:
			values.Add(k, strconv.Itoa(int(str)))
		case float64:
			values.Add(k, strconv.FormatFloat(str, 'f', 2, 64))
		}
	}
	mockURI := url.URL{
		Path:     path,
		RawQuery: values.Encode(),
	}
	req, err := http.NewRequest(http.MethodGet, m.Host+mockURI.RequestURI(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = request(ctx, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (m *HttpClient) Delete(ctx context.Context, path string, data map[string]interface{}, res *Result) (err error) {
	m.Host = "http://" + m.Host
	values := url.Values{}
	for k, v := range data {
		values.Add(k, v.(string))
	}
	mockURI := url.URL{
		Path:     path,
		RawQuery: values.Encode(),
	}
	req, err := http.NewRequest(http.MethodDelete, m.Host+mockURI.RequestURI(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = request(ctx, req, res)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func request(ctx context.Context, req *http.Request, res *Result) (err error) {
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Transport = NewOTELHTTPTransport(client.Transport)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, res)
	if err != nil {
		return
	}
	return
}
