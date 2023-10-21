package client

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_sigV1_generate(t *testing.T) {
	type fields struct {
		nowStr          string
		secret          string
		apiKey          string
		literal         string
		saasID          string
		httpMethod      string
		httpURL         string
		httpRequestBody io.ReadCloser
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "success_GET",
			fields: fields{
				nowStr:          "202109011234",
				secret:          "secret",
				apiKey:          "apiKey",
				literal:         "SAASUSSIGV1",
				saasID:          "saasID",
				httpMethod:      "GET",
				httpURL:         "https://api.saasus.io/v1/auth/userinfo",
				httpRequestBody: nil,
			},
			want:    "SAASUSSIGV1 Sig=a041e74fb58621f9e2e25453573b6009a93f457a33df970a7d8a5637a593708b, SaaSID=saasID, APIKey=apiKey",
			wantErr: false,
		},
		{
			name: "success_POST",
			fields: fields{
				nowStr:          "202109011234",
				secret:          "secret",
				apiKey:          "apiKey",
				literal:         "SAASUSSIGV1",
				saasID:          "saasID",
				httpMethod:      "POST",
				httpURL:         "https://api.saasus.io/v1/auth/userinfo",
				httpRequestBody: io.NopCloser(bytes.NewReader([]byte("hello world"))),
			},
			want:    "SAASUSSIGV1 Sig=a23ee22c670a2e7756a7bad8bc2847196783e1e2576b0d999f6d1b6ccbf5e3fa, SaaSID=saasID, APIKey=apiKey",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sigV1{
				nowStr:          tt.fields.nowStr,
				secret:          tt.fields.secret,
				apiKey:          tt.fields.apiKey,
				literal:         tt.fields.literal,
				saasID:          tt.fields.saasID,
				httpMethod:      tt.fields.httpMethod,
				httpURL:         tt.fields.httpURL,
				httpRequestBody: tt.fields.httpRequestBody,
			}
			got, err := s.generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("sigV1.generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
