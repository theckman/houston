// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, you can
// obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2017 Tim Heckman

package twilio

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/theckman/houston/twilio/util"
)

func testClient(addr string) *Client {
	return &Client{
		SID:        "x",
		Secret:     "y",
		HTTPClient: util.DefaultClient(),
		BaseURL:    fmt.Sprintf("http://%s", addr),
	}
}

func setUpTestHTTPServer(h http.HandlerFunc) (net.Listener, *http.Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		return nil, nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h)

	server := &http.Server{
		Addr:         listener.Addr().String(),
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	server.SetKeepAlivesEnabled(false)

	go func() { server.Serve(listener) }()

	return listener, server, nil
}

func TestNew(t *testing.T) {
	_, err := New("", "y")

	if err.Error() != "sid cannot be zero length" {
		t.Errorf("New() failed to validate sid to ensure it is valid")
	}

	_, err = New("x", "")

	if err.Error() != "secret cannot be zero length" {
		t.Errorf("New() failed to validate secret to ensure it is valid")
	}

	client, err := New("x", "y")

	if err != nil || client == nil {
		t.Errorf("New(\"x\", \"y\") = %q, %q; want *Client, <nil>", client, err)
	}

	if client.SID != "x" {
		t.Errorf("client.SID = %q; want %q", client.SID, "x")
	}

	if client.Secret != "y" {
		t.Errorf("client.Secret = %q; want %q", client.Secret, "y")
	}

	if client.BaseURL != TwilioAPIBase {
		t.Errorf("client.BaseURL = %q; want %q", client.BaseURL, TwilioAPIBase)
	}
}

func Test_formatResource(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"resource", "/resource"},
		{"/resource/rsrc", "/resource/rsrc"},
		{"resource/rsrc", "/resource/rsrc"},
	}

	for _, tt := range tests {
		if rsrc := formatResource(tt.in); rsrc != tt.out {
			t.Errorf("formatResource(%q) = %q; want %q", tt.in, rsrc, tt.out)
		}
	}
}

func Test_formatValues(t *testing.T) {
	v1, v2, v3 := url.Values{}, url.Values{}, url.Values{}

	v1.Set("testValue", "test")

	v2.Set("testValue", "test")
	v2.Set("otherValue", "wat")

	v3.Add("testValues", "test0")
	v3.Add("testValues", "test1")

	tests := []struct {
		in  url.Values
		out string
	}{
		{url.Values{}, ""},
		{v1, "?testValue=test"},
		{v2, "?otherValue=wat&testValue=test"},
		{v3, "?testValues=test0&testValues=test1"},
	}

	for _, tt := range tests {
		if vals := formatValues(tt.in); vals != tt.out {
			t.Errorf("formatResource(%q) = %q; want %q", tt.in, vals, tt.out)
		}
	}
}

func Test_newRequest(t *testing.T) {
	client := testClient("127.0.0.1:8080")
	v := url.Values{}
	v.Set("testQuery", "set")

	req, err := newRequest(client, "GET", "/q", v)

	if err != nil {
		t.Fatalf("newRequest(client, \"GET\", \"/q\", v) = <nil>, %s; want *http.Request, <nil>", err.Error())
	}

	if req.Method != "GET" {
		t.Errorf("req.Method = %q; want %q", req.Method, "GET")
	}

	uaStr := fmt.Sprintf(
		"Twilio/%s (github.com/theckman/houston/twilio) Go-http-client/1.1 (%s; %s %s)",
		Version, runtime.Version(), runtime.GOOS, runtime.GOARCH,
	)

	if ua := req.Header.Get("User-Agent"); ua != uaStr {
		t.Errorf("req.Header.Get(\"User-Agent\") = %q; want %q", ua, uaStr)
	}

	if user, pass, ok := req.BasicAuth(); !ok || (user != "x" && pass != "y") {
		t.Errorf("req.BasicAuth = %q, %q, %q; want \"x\", \"y\", true ", user, pass, ok)
	}
}

func TestClient_get(t *testing.T) {
	l, s, err := setUpTestHTTPServer(func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()

		// if there is no Basic Authentication provided
		// or if the User/Pass combination is not x:y
		// return a 401
		if user, pass, ok := r.BasicAuth(); !ok || (user != "x" && pass != "y") {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "go away")
			return
		}

		// if there are no query params
		// and the request was to /
		// this is expected and OK
		if len(vals) == 0 && r.URL.Path == "/x.json" {
			fmt.Fprint(w, "imok")
			return
		}

		if len(vals) == 1 && r.URL.Path == "/x/q.json" {
			if len(vals["testQuery"]) == 1 {
				fmt.Fprintf(w, "imok:%s", vals["testQuery"][0])
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "nook")
	})

	if err != nil {
		t.Fatalf("setUpTestHTTPServer() = %s; want <nil>", err.Error())
	}

	defer func() {
		s.Close()
		l.Close()
	}()

	//
	// Test GET with no parameters
	//
	client := testClient(l.Addr().String())

	resp1, err := client.get("", nil)

	if err != nil {
		t.Fatalf("client.get(\"\", nil) = <nil>, %s; want <nil>", err.Error())
	}

	defer resp1.Body.Close()

	if resp1.StatusCode != 200 {
		t.Fatalf("resp1.StatusCode = %d; want 200", resp1.StatusCode)
	}

	body, err := ioutil.ReadAll(resp1.Body)

	if err != nil {
		t.Fatalf("ioutil.ReadAll(%q) = %q, %s; want []byte(\"imok\"), <nil>", resp1.Body, body, err.Error())
	}

	if bodyStr := string(body); bodyStr != "imok" {
		t.Fatalf("string(body) = %s; want \"imok\"", bodyStr)
	}

	//
	// Test GET With Query Parameters
	//
	v := url.Values{}
	v.Set("testQuery", "set")

	resp2, err := client.get("/q", v)

	if err != nil {
		t.Fatalf("client.get(\"/q\", %q) = <nil>, %s; want <nil>", v, err.Error())
	}

	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		t.Fatalf("resp2.StatusCode = %d; want 200", resp2.StatusCode)
	}

	body, err = ioutil.ReadAll(resp2.Body)

	if err != nil {
		t.Fatalf("ioutil.ReadAll(%q) = %q, %s; want []byte(\"imok:set\"), <nil>", resp2.Body, body, err.Error())
	}

	if bodyStr := string(body); bodyStr != "imok:set" {
		t.Fatalf("string(body) = %s; want \"imok:set\"", bodyStr)
	}
}
