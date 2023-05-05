package dex

import (
	"context"
	"testing"
)

func Test_xx(t *testing.T) {

	conn, err := NewConn()
	if err != nil {
		t.Fatal(err)
	}

	c := &Client{
		ID:           "test",
		Secret:       "test",
		RedirectURIs: make([]string, 0),
	}
	ret, err := conn.Create(context.TODO(), c)
	if err != nil {
		t.Fatalf("1ST ERROR: %v", err.Error())
	}
	t.Log(ret)

	ret, err = conn.Create(context.TODO(), c)
	if err != nil {
		t.Fatalf("2ND ERROR: %v", err.Error())
	}
	t.Log(ret)
}
