package yadisk

import (
	"context"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var (
	// TestData
	// testUploadFilePath = "testdata/upload.txt"
	// Token
	testValidToken = Token{
		AccessToken: os.Getenv("YANDEX_TOKEN"),
	}
	testInvalidToken = Token{
		AccessToken: "AQA0AA00qEYz00WXA7olo",
	}
	// Error
	//testUnAuthError = Error{
	//	Message:     "Не авторизован.",
	//	Description: "Unauthorized",
	//	ErrorID:     "UnauthorizedError",
	//}
	// Struct

	// Disk
	testYaDisk, _                 = NewYaDisk(context.Background(), nil, &testValidToken)
	testYaDiskWithInvalidToken, _ = NewYaDisk(context.Background(), nil, &testInvalidToken)
)

func TestNewYaDisk(t *testing.T) {
	type args struct {
		ctx    context.Context
		token  *Token
		client *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    YaDisk
		wantErr bool
	}{
		{"success_test", args{context.Background(), &testValidToken, http.DefaultClient}, testYaDisk, false},
		{"error_test", args{context.Background(), nil, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewYaDisk(tt.args.ctx, tt.args.client, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewYaDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewYaDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_yandexDisk_GetDisk(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		yaDisk  YaDisk
		wantD   *Disk
		wantErr bool
	}{
		{"success_test", args{[]string{"is_paid"}}, testYaDisk, &Disk{IsPaid: false}, false},
		{"error_test", args{[]string{}}, testYaDiskWithInvalidToken, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := tt.yaDisk.GetDisk(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("yandexDisk.GetDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("yandexDisk.GetDisk() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func Test_yandexDisk_PerformUpload(t *testing.T) {
	type fields struct {
		Token  *Token
		client *client
	}
	type args struct {
		ur   *ResourceUploadLink
		data io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantPu  *PerformUpload
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yad := &yandexDisk{
				Token:  tt.fields.Token,
				client: tt.fields.client,
			}
			gotPu, err := yad.PerformUpload(tt.args.ur, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("yandexDisk.PerformUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPu, tt.wantPu) {
				t.Errorf("yandexDisk.PerformUpload() = %v, want %v", gotPu, tt.wantPu)
			}
		})
	}
}
