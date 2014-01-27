package babelcoin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

type JsonRPCClient struct {
	Url, Key, Secret string
}

// generate hmac-sha512 hash, hex encoded
func (c *JsonRPCClient) sign(secret string, payload string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// returns a url encoded string to be signed an send
func (c *JsonRPCClient) encodePostData(method string, params map[string]string) string {
	nonce := time.Now().Unix()
	result := fmt.Sprintf("method=%s&nonce=%d", method, nonce)

	// params are unordered, but after method and nonce
	if len(params) > 0 {
		v := url.Values{}
		for key := range params {
			v.Add(key, params[key])
		}
		result = result + "&" + v.Encode()
	}

	return result
}

// marshal an api response into an object
func (c *JsonRPCClient) marshalResponse(resp *http.Response, v interface{}) error {
	// read the response
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll() failed: %v", err)
		return err
	}

	data := &struct {
		Success interface{}     `json:"success"`
		Return  json.RawMessage `json:"return"`
		Error   string          `json:"error"`
	}{}

	if err = json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	// some impls use strings, some ints
	var success int
	if val, ok := data.Success.(int); ok {
		success = val
	} else if val, ok := data.Success.(string); ok {
		success, err = strconv.Atoi(val)
		if err != nil {
			return err
		}
	}

	if success != 1 {
		log.Fatalf("Request failed: %v", data.Error)
		return errors.New("Request failed: " + data.Error)
	}

	if err = json.Unmarshal(data.Return, &v); err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
		return err
	}

	return nil
}

// make a call to the jsonrpc api, marshal into v
func (c *JsonRPCClient) Call(method string, v interface{}, params map[string]string) error {
	client := &http.Client{}
	postData := c.encodePostData(method, params)

	r, err := http.NewRequest("POST", c.Url, bytes.NewBufferString(postData))
	if err != nil {
		return err
	}

	r.Header.Add("Sign", c.sign(c.Secret, postData))
	r.Header.Add("Key", c.Key)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	if os.Getenv("HTTP_DEBUG") != "" {
		bytes, _ := httputil.DumpRequest(r, os.Getenv("HTTP_DEBUG") == "2")
		fmt.Println(string(bytes))
	}

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	if os.Getenv("HTTP_DEBUG") != "" {
		bytes, _ := httputil.DumpResponse(resp, os.Getenv("HTTP_DEBUG") == "2")
		fmt.Println(string(bytes))
	}

	return c.marshalResponse(resp, v)
}
