package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func eventSourceHandler(lines ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		for _, l := range lines {
			fmt.Fprintln(w, l)
		}
	})
}

func TestConnectToMarathonEventStream(t *testing.T) {
	for _, test := range []struct {
		desc     string
		user     string
		password string
		handler  http.Handler
		count    int
		payloads [][]byte
	}{
		{
			desc: "not-found",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
		},
		{
			desc:    "single-payload",
			handler: eventSourceHandler("data: payload"),
			count:   1,
			payloads: [][]byte{
				[]byte("payload\n"),
			},
		},
		{
			desc:    "heartbeat",
			handler: eventSourceHandler("", "", ""),
			count:   0,
		},
		{
			desc:    "event line",
			handler: eventSourceHandler("event: eventName"),
			count:   0,
		},
		{
			desc:    "mixed content",
			handler: eventSourceHandler("", "event: eventName", "data: payload", "", "", "event: eventName2", "data: payload2", ""),
			count:   2,
			payloads: [][]byte{
				[]byte("payload\n"),
				[]byte("payload2\n"),
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			ch := connectToMarathonEventStream(server.URL, test.user, test.password)

			count := 0
			for p := range ch {
				count++

				if count > len(test.payloads) {
					t.Errorf("got payload %s, but wanted none", p)
					continue
				}

				expected := test.payloads[count-1]
				if !reflect.DeepEqual(p, expected) {
					t.Errorf("got payload %s, wanted %s", p, expected)
				}
			}

			if count != test.count {
				t.Errorf("got %d payloads, wanted %d", count, test.count)
			}
		})
	}

	t.Run("url-error", func(t *testing.T) {
		ch := connectToMarathonEventStream("unknown://", "", "")

		count := 0
		for _ = range ch {
			count++
		}

		if count != 0 {
			t.Errorf("got %d payloads, wanted none", count)
		}
	})

	t.Run("request-error", func(t *testing.T) {
		server := httptest.NewServer(nil)
		server.Close()

		ch := connectToMarathonEventStream(server.URL, "", "")

		count := 0
		for _ = range ch {
			count++
		}

		if count != 0 {
			t.Errorf("got %d payloads, wanted none", count)
		}
	})
}
