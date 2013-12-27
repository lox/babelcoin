package btce

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
	"net/url"
	"strconv"
	"time"
)

// generate hmac-sha512 hash, hex encoded
func sign(secret string, payload string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// returns a url encoded string to be signed an send
func (b *Driver) encodePostData(method string, params map[string]string) string {
	// for some reason, btce barfs on larger nonces
	nonce := time.Now().UnixNano() / 1000000000
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

// marshal an api response into an object, with a form for success and
func (b *Driver) marshalResponse(resp *http.Response, v interface{}) error {
	// read the response
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll() failed: %v", err)
		return err
	}

	// sometimes, btc-e returns a non-json error
	if string(bytes) == "invalid POST data" {
		log.Fatalf("Invalid post data")
		return errors.New("Request failed: invalid post data")
	}

	data := &struct {
		Success int             `json:"success"`
		Return  json.RawMessage `json:"return"`
		Error   string          `json:"error"`
	}{}

	if err = json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	if data.Success != 1 {
		log.Fatalf("Request failed: %v", data.Error)
		return errors.New("Request failed: " + data.Error)
	}

	if err = json.Unmarshal(data.Return, &v); err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
		return err
	}

	return nil
}

// make a call to the btc-e api, marshal into v
func (d *Driver) privateApiCall(method string, v interface{}, params map[string]string) error {
	client := &http.Client{}
	postData := d.encodePostData(method, params)

	r, err := http.NewRequest("POST", d.privateApi, bytes.NewBufferString(postData))
	if err != nil {
		return err
	}

	r.Header.Add("Sign", sign(d.config["secret"].(string), postData))
	r.Header.Add("Key", d.config["key"].(string))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	return d.marshalResponse(resp, v)
}
