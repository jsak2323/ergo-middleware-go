package ergo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/btcid/ergo-middleware-go/cmd/config"
)

type (
	Http struct {
		Url    string
		Method string
	}

	Resp struct {
		Err
		BodyResp *interface{}
	}

	// Error struct {
	// 	Error  int    `json:"error"`
	// 	Reason string `json:"reason"`
	// 	Detail string `json:"detail"`
	// }
)

// Client setting
var c = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 2 * time.Minute,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 2 * time.Minute,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func (h *Http) HttpCall(ctx context.Context, path string, m *Resp) (e error) {

	var body io.Reader

	if h.Method != "GET" {
		jsonBody, e := json.Marshal(&m)
		if e != nil {
			return e
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(h.Method, config.CONF.NodeJsonHtppUrl+path, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api_key", config.CONF.NodeJsonHtppApiKey)

	res, e := c.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	// Response not used yet
	resp, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return e
	}

	// Success is indicated with 2xx status codes:
	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		e = json.Unmarshal(resp, &m)
		log.Println("Non-OK HTTP status:", res.StatusCode)
		return errors.New(string(resp))
	}

	// e = json.Unmarshal(resp, &respStruct)

	return nil
}
