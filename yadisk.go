package yadisk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	BaseURL           string = "https://cloud-api.yandex.net"
	MaxFileUploadSize int64  = 1e10
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

// This custom method to upload data by link.
func (yad *yandexDisk) PerformUpload(ur *ResourceUploadLink, data *bytes.Buffer) (pu *PerformUpload, e error) {
	req, e := http.NewRequest(ur.Method, ur.Href, data)
	if e != nil {
		return
	}
	return yad.performUpload(req)
}

// This custom method to upload data by link.
//
// portions - the number of portions to upload the file. data len / portions
//
// partSize - if partSize > 1e10 then partsSize = 1e10 (upload file max size = 1e10)
func (yad *yandexDisk) PerformPartialUpload(ur *ResourceUploadLink, data *bytes.Buffer, partSize int64) (pu *PerformUpload, e error) {
	contentLength := int64(data.Len())
	var wg sync.WaitGroup
	if partSize > contentLength {
		return nil, fmt.Errorf("partSize can not be more than data.Len()")
	}
	if partSize > MaxFileUploadSize {
		log.Printf("partSize %v > MaxFileUploadSize %v. change value partSize on %v", partSize, MaxFileUploadSize, MaxFileUploadSize)
		partSize = MaxFileUploadSize
	}
	portions := int(contentLength / partSize)
	resultsChan := make(chan *performPartialUploadResult, portions)
	defer close(resultsChan)

	reqs, e := requestWithRange(ur, data.Bytes(), partSize, contentLength, portions)
	if e != nil {
		return nil, e
	}

	for _, req := range reqs {
		wg.Add(1)
		go func(r *http.Request) {
			defer wg.Done()
			pu, err := yad.performUpload(r)
			res := &performPartialUploadResult{pu, err}
			resultsChan <- res
		}(req)
		wg.Wait()
	}
	var results []performPartialUploadResult
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

func (yad *yandexDisk) performUpload(req *http.Request) (pu *PerformUpload, e error) {
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
