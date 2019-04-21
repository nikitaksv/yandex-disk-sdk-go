package yadisk

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Empty trash.
//
// If the deletion occurs asynchronously, it will return a response with status 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with status 204 and an empty body.
//
// If the path parameter is not specified or points to the root of the Recycle Bin,
// the recycle bin will be completely cleared, otherwise only the resource pointed to by the path will be deleted from the Recycle Bin.
func (yad *yandexDisk) ClearTrash(fields []string, forceAsync bool, path string) (l *Link, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))
	values.Add("force_async", strconv.FormatBool(forceAsync))
	values.Add("path", path)

	req, e := yad.client.request(http.MethodDelete, "/disk/trash/resources?"+values.Encode(), nil)
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

// Get the contents of the Trash.
func (yad *yandexDisk) GetTrashResource(path string, fields []string, limit int, offset int, previewCrop bool, previewSize string, sort string) (r *TrashResource, e error) {
	req, e := yad.getResource("trash/", path, fields, limit, offset, previewCrop, previewSize, sort)
	if e != nil {
		return nil, e
	}

	r = new(TrashResource)
	_, e = yad.client.getResponse(req, &r)
	if e != nil {
		return nil, e
	}
	return
}

// Recover Resource from Trash.
//
// If recovery is asynchronous, it will return a response with code 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with code 201 and a link to the created resource.
func (yad *yandexDisk) RestoreFromTrash(path string, fields []string, forceAsync bool, name string, overwrite bool) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("force_async", strconv.FormatBool(forceAsync))
	values.Add("name", name)
	values.Add("overwrite", strconv.FormatBool(overwrite))

	req, e := yad.client.request(http.MethodPut, "/disk/trash/resources/restore?"+values.Encode(), nil)
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
