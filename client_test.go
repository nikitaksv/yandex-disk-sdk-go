package yadisk

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	duration2 = time.Duration(2) * time.Second
)

func createClient(ctx context.Context, url string) *client {
	client, _ := newClient(ctx, &testValidToken, url, 1, http.DefaultClient)
	return client
}

func createContextWithTimeout(duration time.Duration) (ctx context.Context, cancelFunc context.CancelFunc) {
	ctx, cancelFunc = context.WithTimeout(context.Background(), duration)
	return
}

func testGetDiskRequest(testClient *client) *http.Request {
	getDiskRequest, _ := testClient.request(http.MethodGet, "/disk", nil)
	return getDiskRequest
}

func Test_newClient(t *testing.T) {

	type args struct {
		ctx        context.Context
		httpClient *http.Client
		token      *Token
		baseURL    string
		version    int
	}
	tests := []struct {
		name    string
		args    args
		want    *client
		wantErr bool
	}{
		{"success_test", args{
			ctx:        context.Background(),
			httpClient: http.DefaultClient,
			token:      &testValidToken,
			baseURL:    BaseURL,
			version:    1,
		}, createClient(context.Background(), BaseURL), false},
		{"error_test", args{
			ctx:        context.Background(),
			httpClient: http.DefaultClient,
			token:      &testValidToken,
			baseURL:    BaseURL + "\\asd:..//'",
			version:    1,
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newClient(tt.args.ctx, tt.args.token, tt.args.baseURL, tt.args.version, tt.args.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("newClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_setRequestHeaders(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		c    *client
		args args
	}{
		{"success_test", createClient(context.Background(), BaseURL), args{req: testGetDiskRequest(createClient(context.Background(), BaseURL))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.setRequestHeaders(tt.args.req)
		})
	}
}

func Test_client_request(t *testing.T) {
	type args struct {
		method  string
		pathURL string
		body    io.Reader
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    *http.Request
		wantErr bool
	}{
		{"success_test", createClient(context.Background(), BaseURL), args{
			http.MethodGet,
			"/disk",
			nil,
		}, testGetDiskRequest(createClient(context.Background(), BaseURL)), false},
		{"error_test", createClient(context.Background(), BaseURL), args{
			http.MethodGet,
			"/disk\\%$***::\\",
			nil,
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.request(tt.args.method, tt.args.pathURL, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil && !reflect.DeepEqual(got.URL, tt.want.URL) {
				t.Errorf("client.request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_do(t *testing.T) {
	ctx, cancel := createContextWithTimeout(duration2)
	defer cancel()
	testClientFail := createClient(ctx, "http://example.com:8088")
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		{"success_test", createClient(context.Background(), BaseURL), args{testGetDiskRequest(createClient(context.Background(), BaseURL))}, false},
		{"timeout_error_test", testClientFail, args{testGetDiskRequest(testClientFail)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.do(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.StatusCode != 200 {
				t.Errorf("client.do() = %v, want %v", got.StatusCode, 200)
			}
		})
	}
}

func Test_client_getResponse(t *testing.T) {
	type args struct {
		req *http.Request
		obj interface{}
	}
	requestPut := testGetDiskRequest(createClient(context.Background(), BaseURL))
	requestPut.Method = http.MethodPut

	tests := []struct {
		name    string
		c       *client
		args    args
		wantI   *responseInfo
		wantErr bool
	}{
		{"success_test", createClient(context.Background(), BaseURL), args{
			testGetDiskRequest(createClient(context.Background(), BaseURL)),
			new(Disk)},
			&responseInfo{
				"200 OK",
				200},
			false,
		},
		{"error_test", createClient(context.Background(), BaseURL), args{
			requestPut,
			new(Disk)},
			&responseInfo{
				"405 METHOD NOT ALLOWED",
				405},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, err := tt.c.getResponse(tt.args.req, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.getResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("client.getResponse() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func Test_bodyClose(t *testing.T) {
	type args struct {
		closer io.Closer
	}
	tests := []struct {
		name string
		args args
	}{
		{"success_test", args{closer: http.NoBody}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyClose(tt.args.closer)
		})
	}
}
