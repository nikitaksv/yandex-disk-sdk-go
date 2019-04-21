package yadisk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Delete file or folder.
//
// By default, delete the resource in the trash.
// To delete a resource without placing it in the trash, you must specify the parameter permanently = true.
//
// If the deletion occurs asynchronously, it will return a response with status 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with status 204 and an empty body.
func (yad *yandexDisk) DeleteResource(path string, fields []string, forceAsync bool, md5 string, permanently bool) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("force_async", strconv.FormatBool(forceAsync))
	values.Add("md5", md5)
	values.Add("permanently", strconv.FormatBool(permanently))

	req, e := yad.client.request(http.MethodDelete, "/disk/resources?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(Link)
	ri, e := yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	if ri.StatusCode == 204 {
		return nil, nil
	}
	return
}

// Get meta information about a file or directory.
func (yad *yandexDisk) GetResource(path string, fields []string, limit int, offset int, previewCrop bool, previewSize string, sort string) (r *Resource, e error) {
	req, e := yad.getResource("", path, fields, limit, offset, previewCrop, previewSize, sort)
	if e != nil {
		return nil, e
	}

	r = new(Resource)
	_, e = yad.client.getResponse(req, &r)
	if e != nil {
		return nil, e
	}
	return
}

// If the path points to a directory, the response also describes the resources of that directory.
func (yad *yandexDisk) getResource(area string, path string, fields []string, limit int, offset int, previewCrop bool, previewSize string, sort string) (*http.Request, error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("offset", strconv.Itoa(offset))
	values.Add("preview_crop", strconv.FormatBool(previewCrop))
	values.Add("preview_size", previewSize)
	values.Add("sort", sort)

	r, e := yad.client.request(http.MethodGet, "/disk/"+area+"resources?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}
	return r, nil
}

// Create directory.
func (yad *yandexDisk) CreateResource(path string, fields []string) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodPut, "/disk/resources?"+values.Encode(), nil)
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

// Update User Resource Data.
func (yad *yandexDisk) UpdateResource(path string, fields []string, body *ResourcePatch) (r *Resource, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	bodyJSON, e := json.Marshal(body)
	if e != nil {
		return nil, e
	}
	req, e := yad.client.request(http.MethodPatch, "/disk/resources?"+values.Encode(), bytes.NewReader(bodyJSON))
	if e != nil {
		return nil, e
	}

	r = new(Resource)
	_, e = yad.client.getResponse(req, &r)
	if e != nil {
		return nil, e
	}
	return
}

// Create a copy of the file or folder.
//
// If copying occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with code 201 and a link to the created resource.
func (yad *yandexDisk) CopyResource(from string, path string, fields []string, forceAsync bool, overwrite bool) (l *Link, e error) {
	return yad.transportResource("copy", from, path, fields, forceAsync, overwrite)
}

// Move a file or folder.
//
// If the movement occurs asynchronously, it will return a response with code 202 and a link to the asynchronous operation.
// Otherwise, it will return a response with code 201 and a link to the created resource.
func (yad *yandexDisk) MoveResource(from string, path string, fields []string, forceAsync bool, overwrite bool) (l *Link, e error) {
	return yad.transportResource("move", from, path, fields, forceAsync, overwrite)
}

func (yad *yandexDisk) transportResource(copyMove string, from string, path string, fields []string, forceAsync bool, overwrite bool) (l *Link, e error) {
	values := url.Values{}
	values.Add("from", from)
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("force_async", strconv.FormatBool(forceAsync))
	values.Add("overwrite", strconv.FormatBool(overwrite))

	req, e := yad.client.request(http.MethodPost, "/disk/resources/"+copyMove+"?"+values.Encode(), nil)
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

// Get link to download file.
func (yad *yandexDisk) GetResourceDownloadLink(path string, fields []string) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodGet, "/disk/resources/download?"+values.Encode(), nil)
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

// Get file list sorted by name.
func (yad *yandexDisk) GetFlatFilesList(fields []string, limit int, mediaType string, offset int, previewCrop bool, previewSize string, sort string) (l *FilesResourceList, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("media_type", mediaType)
	values.Add("offset", strconv.Itoa(offset))
	values.Add("preview_crop", strconv.FormatBool(previewCrop))
	values.Add("preview_size", previewSize)
	values.Add("sort", sort)

	req, e := yad.client.request(http.MethodGet, "/disk/resources/files?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(FilesResourceList)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}

// Get a list of files ordered by download date.
func (yad *yandexDisk) GetLastUploadedFilesList(fields []string, limit int, mediaType string, previewCrop bool, previewSize string) (l *LastUploadedResourceList, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("media_type", mediaType)
	values.Add("preview_crop", strconv.FormatBool(previewCrop))
	values.Add("preview_size", previewSize)

	req, e := yad.client.request(http.MethodGet, "/disk/resources/last-uploaded?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(LastUploadedResourceList)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}

// Get a list of published resources.
//
// resourceType value: "","dir","file".
func (yad *yandexDisk) ListPublicResources(fields []string, limit int, offset int, previewCrop bool, previewSize string, resourceType string) (l *PublicResourcesList, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))
	values.Add("limit", strconv.Itoa(limit))
	values.Add("offset", strconv.Itoa(offset))
	values.Add("preview_crop", strconv.FormatBool(previewCrop))
	values.Add("preview_size", previewSize)
	values.Add("type", resourceType)

	req, e := yad.client.request(http.MethodGet, "/disk/resources/public?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(PublicResourcesList)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}

// Publish a resource.
func (yad *yandexDisk) PublishResource(path string, fields []string) (l *Link, e error) {
	return yad.pubResource("publish", path, fields)
}

// Unpublish resource.
func (yad *yandexDisk) UnpublishResource(path string, fields []string) (l *Link, e error) {
	return yad.pubResource("unpublish", path, fields)
}

func (yad *yandexDisk) pubResource(publishUnpublish string, path string, fields []string) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodPut, "/disk/resources/"+publishUnpublish+"?"+values.Encode(), nil)
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

// Upload file to Disk by URL.
//
// Download asynchronously.
//
// Therefore, in response to the request, a reference to the asynchronous operation is returned.
func (yad *yandexDisk) UploadExternalResource(path string, externalURL string, disableRedirects bool, fields []string) (l *Link, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("url", externalURL)
	values.Add("disable_redirects", strconv.FormatBool(disableRedirects))
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodPost, "/disk/resources/upload?"+values.Encode(), nil)
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

// Get file download link.
func (yad *yandexDisk) GetResourceUploadLink(path string, fields []string, overwrite bool) (l *ResourceUploadLink, e error) {
	values := url.Values{}
	values.Add("path", path)
	values.Add("fields", strings.Join(fields, ","))
	values.Add("overwrite", strconv.FormatBool(overwrite))

	req, e := yad.client.request(http.MethodGet, "/disk/resources/upload?"+values.Encode(), nil)
	if e != nil {
		return nil, e
	}

	l = new(ResourceUploadLink)
	_, e = yad.client.getResponse(req, &l)
	if e != nil {
		return nil, e
	}
	return
}
