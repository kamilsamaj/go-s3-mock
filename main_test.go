package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io/ioutil"
	"strconv"
	"testing"
)


type mockGetObjectAPI struct {
	t *testing.T
}

func (m mockGetObjectAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	m.t.Helper()
	if params.Bucket == nil {
		m.t.Fatal("expect bucket to not be nil")
	}
	if e, a := "fooBucket", *params.Bucket; e != a {
		m.t.Errorf("expect %v, got %v", e, a)
	}
	if params.Key == nil {
		m.t.Fatal("expect key to not be nil")
	}
	if e, a := "barKey", *params.Key; e != a {
		m.t.Errorf("expect %v, got %v", e, a)
	}

	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("this is the body foo bar baz"))),
	}, nil
}

func TestGetObjectFromS3(t *testing.T) {
	cases := []struct {
		client S3GetObjectAPI
		bucket string
		key	string
		expect []byte
	}{
		{
			client: mockGetObjectAPI{t},
			bucket: "fooBucket",
			key:	"barKey",
			expect: []byte("this is the body foo bar baz"),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := GetObjectFromS3(ctx, tt.client, tt.bucket, tt.key)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; bytes.Compare(e, a) != 0 {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
