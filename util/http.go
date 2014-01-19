package babelcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
)

const (
	HTTP_DEBUG = false
)

type HttpError struct {
	NestedError  error
	ResponseBody []byte
}

func (e *HttpError) Error() string {
	return e.NestedError.Error()
}

// fetch a json response from provided url and unmarshal into the provided r
func HttpGetJson(url string, r interface{}) *HttpError {
	bytes, err := HttpDurableGet(url, 10)
	if err != nil {
		return &HttpError{err, bytes}
	}

	if err = json.Unmarshal(bytes, &r); err != nil {
		return &HttpError{err, bytes}
	}

	return nil
}

// attempt an HTTP GET, retrying up to n times
func HttpDurableGet(url string, times int) ([]byte, error) {
	var body []byte

	for i := 0; i < times; i++ {
		resp, err := http.Get(url)

		if HTTP_DEBUG {
			if err != nil {
				log.Println("HTTP Get failed: " + err.Error())
			}

			bytes, _ := httputil.DumpResponse(resp, false)
			fmt.Println(string(bytes))
		}

		switch err.(type) {
		case *urlpkg.Error:
			return []byte{}, err
		}

		if err != nil && i == times {
			return []byte{}, err
		} else if err != nil {
			continue
		}

		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil && i == times {
			return []byte{}, err
		} else if err != nil {
			continue
		} else {
			break
		}
	}

	return body, nil
}
