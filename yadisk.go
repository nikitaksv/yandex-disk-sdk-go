package yadisk

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
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
func (yad *yandexDisk) PerformUpload(ur *ResourceUploadLink, data io.Reader) (pu *PerformUpload, e error) {
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
	return nil, nil
}
