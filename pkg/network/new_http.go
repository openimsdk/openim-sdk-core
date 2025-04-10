package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"unsafe"

	"github.com/openimsdk/tools/errs"
)

type HttpCli struct {
	httpClient  *http.Client
	httpRequest *http.Request
	Error       error
}

func newHttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func PostWithTimeOut(url string, data interface{}, token string, timeout time.Duration) (content []byte, err error) {
	return Post(url).BodyWithJson(data).SetTimeOut(timeout).SetHeader("token", token).ToBytes()
}

func Get(url string) *HttpCli {
	request, err := http.NewRequest("GET", url, nil)
	return &HttpCli{
		httpClient:  newHttpClient(),
		httpRequest: request,
		Error:       err,
	}
}

func Post(url string) *HttpCli {
	request, err := http.NewRequest("POST", url, nil)
	return &HttpCli{
		httpClient:  newHttpClient(),
		httpRequest: request,
		Error:       errs.WrapMsg(err, "newRequest failed, url"),
	}
}

func (c *HttpCli) SetTimeOut(timeout time.Duration) *HttpCli {
	c.httpClient.Timeout = timeout
	return c
}

func (c *HttpCli) SetHeader(key, value string) *HttpCli {
	c.httpRequest.Header.Set(key, value)
	return c
}

func (c *HttpCli) BodyWithJson(obj interface{}) *HttpCli {
	if c.Error != nil {
		return c
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		c.Error = errs.WrapMsg(err, "marshal failed, url")
		return c
	}
	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	c.httpRequest.Header.Set("Content-Type", "application/json")
	return c
}

func (c *HttpCli) BodyWithBytes(buf []byte) *HttpCli {
	if c.Error != nil {
		return c
	}

	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	return c
}

func (c *HttpCli) BodyWithForm(form map[string]string) *HttpCli {
	if c.Error != nil {
		return c
	}

	var value url.Values = make(map[string][]string, len(form))
	for k, v := range form {
		value.Add(k, v)
	}
	buf := Str2bytes(value.Encode())

	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	c.httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (c *HttpCli) ToBytes() (content []byte, err error) {
	if c.Error != nil {
		return nil, c.Error
	}

	resp, err := c.httpClient.Do(c.httpRequest)
	if err != nil {
		return nil, errs.WrapMsg(err, "client.Do failed, url")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.WrapMsg(errors.New(resp.Status), "status code failed ")
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapMsg(err, "ioutil.ReadAll failed, url")
	}

	return buf, nil
}

func (c *HttpCli) ToJson(obj interface{}) error {
	if c.Error != nil {
		return c.Error
	}

	resp, err := c.httpClient.Do(c.httpRequest)
	if err != nil {
		return errs.WrapMsg(err, "client.Do failed, url")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.WrapMsg(errors.New(resp.Status), "status code failed ")
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errs.WrapMsg(err, "ioutil.ReadAll failed, url")
	}
	err = json.Unmarshal(buf, obj)
	if err != nil {
		return errs.WrapMsg(err, "marshal failed, url")
	}
	return nil
}
