package gyazo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient("DUMMY_ACCESS_TOKEN")
	client.DefaultEndpoint = server.URL
}

func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("DUMMY_ACCESS_TOKEN")
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if actual, expected := c.DefaultEndpoint, defaultEndpoint; actual != expected {
		t.Errorf("NewClient DefaultEndpoint is %v, want %v", actual, expected)
	}

	if actual, expected := c.UploadEndpoint, uploadEndpoint; actual != expected {
		t.Errorf("NewClient UploadEndpoint is %v, want %v", actual, expected)
	}
}

func TestNewClientWithEmptyAccessToken(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Error("Expected error to be returned")
	}
}

func TestList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")

		// This example is response body that from https://gyazo.com/api/docs/image.
		fmt.Fprintln(w, `[
			{
				"image_id": "8980c52421e452ac3355ca3e5cfe7a0c",
				"permalink_url": "http://gyazo.com/8980c52421e452ac3355ca3e5cfe7a0c",
				"thumb_url": "https://i.gyazo.com/thumb/afaiefnaf.png",
				"url": "https://i.gyazo.com/8980c52421e452ac3355ca3e5cfe7a0c.png",
				"type": "png",
				"star": "true",
				"created_at": "2014-05-21 14:23:10+0900"
			}
		]`)
	})

	actual, err := client.List(nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	expected := &[]Image{{
		ID:           "8980c52421e452ac3355ca3e5cfe7a0c",
		PermalinkURL: "http://gyazo.com/8980c52421e452ac3355ca3e5cfe7a0c",
		ThumbURL:     "https://i.gyazo.com/thumb/afaiefnaf.png",
		URL:          "https://i.gyazo.com/8980c52421e452ac3355ca3e5cfe7a0c.png",
		Type:         "png",
		Star:         "true",
		CreatedAt:    "2014-05-21 14:23:10+0900",
	}}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("List returned %+v, want %+v", actual, expected)
	}
}
