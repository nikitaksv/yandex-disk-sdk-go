package yadisk

import (
	"net/http"
	"net/url"
	"strings"
)

// Get the status of an asynchronous operation.
func (yad *yandexDisk) GetOperationStatus(operationID string, fields []string) (s *OperationStatus, e error) {
	values := url.Values{}
	values.Add("fields", strings.Join(fields, ","))

	req, e := yad.client.request(http.MethodGet, "/disk/operations/"+operationID+"?"+values.Encode(), nil)
	if e != nil {
		return
	}

	s = new(OperationStatus)
	_, e = yad.client.getResponse(req, &s)
	if e != nil {
		return nil, e
	}
	return
}
