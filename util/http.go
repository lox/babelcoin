package babelcoin

import (
	"io/ioutil"
    "encoding/json"
    "net/http"
)

func HttpGetJson(url string, m interface{}) (error) {
    resp, error := http.Get(url)
    if error != nil {
    	return error
    }

    defer resp.Body.Close()
    bytes, error := ioutil.ReadAll(resp.Body)
    if error != nil {
    	return error
    }

    if error = json.Unmarshal(bytes, &m); error != nil {
    	return error
    }

    return nil
}