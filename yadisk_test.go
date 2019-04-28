package yadisk

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var (
	// Token
	testValidToken = Token{
		AccessToken: os.Getenv("YANDEX_TOKEN"),
	}
	testInvalidToken = Token{
		AccessToken: "AQA0AA00qEYz00WXA7olo",
	}
	// Disk
	testYaDisk, _                 = NewYaDisk(context.Background(), nil, &testValidToken)
	testYaDiskWithInvalidToken, _ = NewYaDisk(context.Background(), nil, &testInvalidToken)
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
	fileName := randStringBytes(10)
	createFile(fileName, rand.Intn(100)*1e4)
	defer removeFile(fileName)
	link, err := testYaDisk.GetResourceUploadLink("/test/"+fileName, nil, true)
	if err != nil {
		t.Errorf("yandexDisk.GetResourceUploadLink() error = %v", err)
	}
	pu, err := testYaDisk.PerformUpload(link, openFile(fileName))
	if err != nil {
		t.Errorf("testYaDisk.PerformPartialUpload() return error = %v", err)
	}

	if pu == nil {
		t.Errorf("testYaDisk.PerformPartialUpload() return nil PerformUpload = %v", err)
	}

	status, err := testYaDisk.GetOperationStatus(link.OperationID, nil)
	if err != nil {
		t.Errorf("testYaDisk.GetOperationStatus() return error = %v", err)
	}
	if status.Status != "success" {
		t.Errorf("testYaDisk.GetOperationStatus() return error = %v", err)
	}
}

func Test_yandexDisk_PerformPartialUpload(t *testing.T) {
	fileName := randStringBytes(10) + "_partial"
	createFile(fileName, rand.Intn(100)*1e4)
	defer removeFile(fileName)
	link, err := testYaDisk.GetResourceUploadLink("/test/"+fileName, nil, true)
	if err != nil {
		t.Errorf("yandexDisk.GetResourceUploadLink() error = %v", err.Error())
	}
	pu, err := testYaDisk.PerformPartialUpload(link, openFile(fileName), rand.Int63n(100)*1e3)
	if err != nil {
		t.Errorf("testYaDisk.PerformPartialUpload() return error = %v", err.Error())
	}

	if pu == nil {
		t.Errorf("testYaDisk.PerformPartialUpload() return nil PerformUpload")
	}

	status, err := testYaDisk.GetOperationStatus(link.OperationID, nil)
	if err != nil {
		t.Errorf("testYaDisk.GetOperationStatus() return error = %v", err.Error())
	}
	if status.Status != "success" {
		t.Errorf("testYaDisk.GetOperationStatus() return bad status = %v", status.Status)
	}
}

func createFile(name string, size int) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = f.WriteString(randStringBytes(size))
	if err != nil {
		panic(err)
	}
}

func removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		panic(err)
	}
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func openFile(name string) (buffer *bytes.Buffer) {
	data, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := data.Close()
		if err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(data)
	buffer = bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	for {
		var count int
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
		log.Fatal("Error Reading ", name, ": ", err)
	} else {
		err = nil
	}
	return
}
