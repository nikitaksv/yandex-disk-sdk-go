package yadisk

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Get meta-information about a public file or directory.
func (yad *yandexDisk) GetPublicResource(publicKey string, fields []string, limit int, offset int, path string, previewCrop bool, previewSize string, sort string) (r *PublicResource, e error) {
	values := url.Values{}
	values.Add("public_key", publicKey)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("offset", strconv.Itoa(offset))
	values.Add("path", path)
	values.Add("preview_crop", strconv.FormatBool(previewCrop))
	values.Add("preview_size", previewSize)
	values.Add("sort", sort)

	req, e := yad.client.request(http.MethodGet, "/disk/public/resources?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	r = new(PublicResource)
	_, e = yad.client.getResponse(req, &r)
	if e != nil {
		return nil, e
	}
	return
}

// Get a link to download a public resource.
func (yad *yandexDisk) GetPublicResourceDownloadLink(publicKey string, fields []string, path string) (l *Link, e error) {
	values := url.Values{}
	values.Add("public_key", publicKey)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("path", path)

	req, e := yad.client.request(http.MethodGet, "/disk/public/resources/download?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(Link)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}

// Save the public resource to the Downloads folder.
//
// If saving occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with code 201 and a link to the created resource.
func (yad *yandexDisk) SaveToDiskPublicResource(publicKey string, fields []string, forceAsync bool, name string, path string, savePath string) (l *Link, e error) {
	values := url.Values{}
	values.Add("public_key", publicKey)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("force_async", strconv.FormatBool(forceAsync))
	values.Add("name", name)
	values.Add("path", path)
	values.Add("save_path", savePath)

	req, e := yad.client.request(http.MethodPost, "/disk/public/resources/save-to-disk?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(Link)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}
