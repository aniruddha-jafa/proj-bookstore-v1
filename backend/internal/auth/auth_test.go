package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		headers http.Header
		want    string
	}{
		{
			headers: http.Header{
				"Authorization": {"Bearer abcde"},
				"Content-Type":  {"text"},
			},
			want: "abcde",
		},
		{
			headers: http.Header{
				"Authorization": {"Bearer ", " ", "abcde"},
				"Content-Type":  {"text"},
			},
			want: "abcde",
		},
	}
	for _, tt := range tests {
		got, err := GetBearerToken(tt.headers)
		if err != nil {
			t.Errorf("Failed to get bearer token: %v", err)
		}
		if got != tt.want {
			t.Errorf("want: %s, got: %s", tt.want, got)
		}
	}
}
