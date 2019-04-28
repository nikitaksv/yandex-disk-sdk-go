package yadisk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// httpClient for send request to Yandex.Disk API
type client struct {
	httpClient *http.Client
	token      *Token
	baseURL    *url.URL
	ctx        context.Context
}

// Construct httpClient
func newClient(ctx context.Context, token *Token, baseURL string, version int, httpClient *http.Client) (*client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	base, e := url.Parse(baseURL + fmt.Sprintf("/v%d", version))
	if e != nil {
		return nil, e
	}

	c := &client{httpClient: httpClient, token: token, baseURL: base, ctx: ctx}
	return c, nil
}

func (c *client) setRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "OAuth "+c.token.AccessToken)
}

func (c *client) request(method string, pathURL string, body io.Reader) (*http.Request, error) {
	rel, e := url.Parse(c.baseURL.Path + pathURL)
	if e != nil {
		return nil, e
	}

	fullURL := c.baseURL.ResolveReference(rel)

	req, e := http.NewRequest(method, fullURL.String(), body)
	if e != nil {
		return nil, e
	}

	c.setRequestHeaders(req)

	return req, nil
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	resp, e := c.httpClient.Do(req.WithContext(c.ctx))
	if e != nil {
		select {
		case <-c.ctx.Done():
			return nil, c.ctx.Err()
		default:
		}
		return nil, e
	}

	return resp, e
}

func (c *client) getResponse(req *http.Request, obj interface{}) (i *responseInfo, e error) {
	resp, e := c.do(req)
	if e != nil {
		return
	}
	defer bodyClose(resp.Body)
	i = new(responseInfo)
	i.setResponseInfo(resp.Status, resp.StatusCode)

	if e != nil {
		return
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return
	}
	if len(body) > 0 {
		err := new(Error)
		e = json.Unmarshal(body, &err)
		if e != nil {
			return
		} else if (Error{}) != *err {
			return i, err
		}
		e = json.Unmarshal(body, &obj)
		if e != nil {
			return
		}
	}
	return i, nil
}

func bodyClose(closer io.Closer) {
	e := closer.Close()
	if e != nil {
		panic(e.Error())
	}
}

func getRange(start, end, total int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", start, end, total)
}

func requestWithRange(ur *ResourceUploadLink, data []byte, partSize, contentLength int64, portions int) ([]*http.Request, error) {
	portionSize := partSize
	startSize := int64(0)
	reqs := make([]*http.Request, portions)
	for i := 0; i < portions; i++ {
		var dataSize []byte
		if i == portions-1 {
			portionSize = contentLength
			dataSize = data[startSize:contentLength]
		} else {
			dataSize = data[startSize:portionSize]
		}
		req, e := http.NewRequest(ur.Method, ur.Href, bytes.NewReader(dataSize))
		if e != nil {
			return nil, e
		}
		req.Header.Set("Content-Range", getRange(startSize, portionSize-1, contentLength))
		reqs[i] = req
		startSize = portionSize
		portionSize += partSize
	}
	return reqs, nil
}
