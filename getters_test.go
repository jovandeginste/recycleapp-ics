package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestURLForLanguage(t *testing.T) {
	org := &Organization{
		URL: map[string]string{
			"nl": "http://example.nl",
			"fr": "http://example.fr",
		},
	}

	if got := org.URLForLanguage("nl"); got != "http://example.nl" {
		t.Errorf("expected http://example.nl, got %s", got)
	}

	if got := org.URLForLanguage("en"); got != "???" {
		t.Errorf("expected ???, got %s", got)
	}
}

func TestGetOrganization(t *testing.T) {
	origTransport := myClient.Transport
	defer func() { myClient.Transport = origTransport }()

	t.Run("success", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.Header.Get("x-consumer") != "recycleapp.be" {
					t.Errorf("expected x-consumer header to be set")
				}

				respJSON := `{"name": "FostPlus", "url": {"nl": "http://fost.nl"}, "description": {"nl": "Fost"}}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		org, err := getOrganization("12345")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if org.Name != "FostPlus" {
			t.Errorf("expected FostPlus, got %s", org.Name)
		}
	})

	t.Run("empty response name", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"name": ""}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getOrganization("12345")
		if !errors.Is(err, ErrOrganizationNoResult) {
			t.Errorf("expected ErrOrganizationNoResult, got %v", err)
		}
	})

	t.Run("http error", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		_, err := getOrganization("12345")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetStreetID(t *testing.T) {
	origTransport := myClient.Transport
	defer func() { myClient.Transport = origTransport }()

	t.Run("success", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [{"id": "street-999"}], "total": 1}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		id, err := getStreetID("12345", "Main Street")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "street-999" {
			t.Errorf("expected street-999, got %s", id)
		}
	})

	t.Run("no items", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [], "total": 0}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getStreetID("12345", "Main Street")
		if !errors.Is(err, ErrStreetNoResult) {
			t.Errorf("expected ErrStreetNoResult, got %v", err)
		}
	})

	t.Run("empty item id", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [{"id": ""}], "total": 1}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getStreetID("12345", "Main Street")
		if !errors.Is(err, ErrStreetNoResult) {
			t.Errorf("expected ErrStreetNoResult, got %v", err)
		}
	})
}

func TestGetZipcodeID(t *testing.T) {
	origTransport := myClient.Transport
	defer func() { myClient.Transport = origTransport }()

	t.Run("success", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [{"id": "zip-777"}], "total": 1}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		id, err := getZipcodeID(1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if id != "zip-777" {
			t.Errorf("expected zip-777, got %s", id)
		}
	})

	t.Run("no items", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [], "total": 0}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getZipcodeID(1000)
		if !errors.Is(err, ErrZipcodeNoResult) {
			t.Errorf("expected ErrZipcodeNoResult, got %v", err)
		}
	})

	t.Run("empty item id", func(t *testing.T) {
		myClient.Transport = &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				respJSON := `{"items": [{"id": ""}], "total": 1}`

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
					Header:     make(http.Header),
				}, nil
			},
		}

		_, err := getZipcodeID(1000)
		if !errors.Is(err, ErrZipcodeNoResult) {
			t.Errorf("expected ErrZipcodeNoResult, got %v", err)
		}
	})
}
