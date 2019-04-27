package yadisk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL string = "https://cloud-api.yandex.net"
)

func NewYaDisk(ctx context.Context, client *http.Client, token *Token) (YaDisk, error) {
	if token == nil || token.AccessToken == "" {
		return nil, errors.New("required token")
	}
	newClient, err := newClient(ctx, token, BaseURL, 1, client)
	if err != nil {
		return nil, err
	}
	return &yandexDisk{Token: token, client: newClient}, nil
}

// Get user disk meta information.
func (yad *yandexDisk) GetDisk(fields []string) (d *Disk, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodGet, "/disk?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	d = new(Disk)
	_, e = yad.client.getResponse(req, &d)
	if e != nil {
		return nil, e
	}
	return
}

//This custom method to download data by link.
func (yad *yandexDisk) PerformUpload(ur *ResourceUploadLink, data *bytes.Buffer) (pu *PerformUpload, e error) {
	req, e := http.NewRequest(ur.Method, ur.Href, data)
	if e != nil {
		return
	}

	pu = new(PerformUpload)
	ri, e := yad.client.getResponse(req, &pu)
	if e != nil {
		return nil, e
	}
	e = pu.handleError(*ri)
	if e != nil {
		return nil, e
	}
	return pu, nil
}

type result struct {
	out *PerformUpload
	err error
}

//This custom method to download data by link.
//
// portions - the number of portions to upload the file. data len / portions
func (yad *yandexDisk) PerformPartialUpload(ur *ResourceUploadLink, data *bytes.Buffer, portions int, concurrencyLimit int) (pu *PerformUpload, e error) {
	if concurrencyLimit > portions {
		return nil, fmt.Errorf("error concurrencyLimit > portions")
	}
	contentLength := int64(data.Len())

	semaphoreChan := make(chan int, concurrencyLimit)
	resultsChan := make(chan *result, portions)
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	portion := func(req *http.Request) (pu *PerformUpload, e error) {
		pu = new(PerformUpload)
		ri, e := yad.client.getResponse(req, &pu)
		if e != nil {
			return nil, e
		}
		e = pu.handleError(*ri)
		if e != nil {
			return nil, e
		}
		return pu, nil
	}

	reqs, e := requestWithRange(ur, data.Bytes(), portions, contentLength)
	if e != nil {
		return nil, e
	}
	for _, req := range reqs {
		// timeout to protect against 500 error
		time.Sleep(200 * time.Millisecond)
		go func(r *http.Request) {
			semaphoreChan <- 1
			pu, err := portion(r)
			res := &result{pu, err}
			resultsChan <- res
			<-semaphoreChan
		}(req)
	}
	var results []result
	for {
		result := <-resultsChan
		results = append(results, *result)
		if len(results) == portions {
			break
		}
	}
	for _, res := range results {
		if res.err != nil {
			return nil, res.err
		}
		if res.out == nil {
			return nil, fmt.Errorf("error permofrm upload")
		}
	}

	pu = &PerformUpload{}
	return pu, nil
}
