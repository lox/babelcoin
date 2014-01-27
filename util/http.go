package babelcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
	"os"
	"time"
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
		timer := time.Now()
		resp, err := http.Get(url)

		// 5xx responses don't count as errors
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			err = errors.New(fmt.Sprintf("Server returned %s", resp.Status))
		}

		if os.Getenv("HTTP_DEBUG") != "" {
			log.Printf("Fetching %s", url)
			if err != nil {
				log.Println("HTTP Get failed: " + err.Error())
			}

			bytes, _ := httputil.DumpResponse(resp, os.Getenv("HTTP_DEBUG") == "2")
			fmt.Println(string(bytes))

			log.Printf("Loaded %d bytes of request in %s", len(bytes), time.Now().Sub(timer))
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

		if os.Getenv("HTTP_DEBUG") != "" {
			log.Printf("Starting reading body")
		}

		readTimer := time.Now()
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if os.Getenv("HTTP_DEBUG") != "" {
			log.Printf("Finished reading response body in %s",
				time.Now().Sub(readTimer))
		}

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
