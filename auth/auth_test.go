package auth

import (
	"encoding/hex"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestBearerToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if _, err := BearerToken(req); err != ErrNoBearerToken {
		t.Errorf("BearerToken() error = %v, wantErr %v", err, ErrNoBearerToken)
	}
	req.Header.Set("Authorization", "any")
	if _, err := BearerToken(req); err != ErrNoBearerToken {
		t.Errorf("BearerToken() error = %v, wantErr %v", err, ErrNoBearerToken)
	}
	req.Header.Set("Authorization", "Bearer 0g")
	if _, err := BearerToken(req); err != hex.InvalidByteError('g') {
		t.Errorf("BearerToken() error = %v, wantErr %v", err, hex.InvalidByteError('g'))
	}
	req.Header.Set("Authorization", "Bearer 0001020304050607")
	got, err := BearerToken(req)
	if err != nil {
		t.Errorf("BearerToken() error = %v, wantErr %v", err, nil)
	}
	want := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BearerToken() got = %v, want %v", got, want)
	}
}
