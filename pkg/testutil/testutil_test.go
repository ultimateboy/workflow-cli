package testutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/arschles/assert"
	deis "github.com/deis/controller-sdk-go"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestNewTestServerAndClient(t *testing.T) {
	t.Parallel()

	cf, server, err := NewTestServerAndClient()
	assert.NoErr(t, err)

	if _, err := os.Stat(cf); os.IsNotExist(err) {
		t.Fatal(err)
	}

	if server.Server.URL == "" {
		t.Fatal(errors.New("No server URL."))
	}

	server.Close()
}

// TestStripProgress ensures StripProgress strips what is expected.
func TestStripProgress(t *testing.T) {
	t.Parallel()

	testInput := "Lorem ipsum dolar sit amet"
	expectedOutput := "Lorem ipsum dolar sit amet"

	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")

	testInput = "Lorem ipsum dolar sit amet...\b\b\b"
	assert.Equal(t, StripProgress(testInput), expectedOutput, "output")
}

// TestAssertBody ensures AssertBody correctly marshals into the interface.
func TestAssertBody(t *testing.T) {
	t.Parallel()

	b := nopCloser{bytes.NewBufferString(`{"data":{"lorem":"ipsum"},"dolar":["sit","amet"]}`)}

	sampleRequest := http.Request{
		Body: b,
	}

	expected := map[string]interface{}{
		"data": map[string]interface{}{
			"lorem": "ipsum",
		},
		"dolar": []string{
			"sit",
			"amet",
		},
	}

	AssertBody(t, expected, &sampleRequest)
}

func TestSetHeaders(t *testing.T) {
	t.Parallel()

	_, server, err := NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	server.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		SetHeaders(w)
		version := w.Header().Get("DEIS_API_VERSION")
		assert.Equal(t, version, deis.APIVersion, "version")
	})
}
