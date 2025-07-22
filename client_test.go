package zei_test

import (
	"bytes"
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ghifari160/zei"
)

const userAgent = "Test_Zei/0.1"

func TestGet(t *testing.T) {
	testClient(t, config(), func(client *zei.Client, url string) (*http.Response, error) {
		return client.Get(url)
	})
}

func TestHead(t *testing.T) {
	testClient(t, config(), func(client *zei.Client, url string) (*http.Response, error) {
		return client.Head(url)
	})
}

func TestPost(t *testing.T) {
	testClient(t, config(), func(client *zei.Client, url string) (*http.Response, error) {
		body := `{"testing":true}`
		return client.Post(url, "application/json", bytes.NewBufferString(body))
	})
}

func TestPostForm(t *testing.T) {
	testClient(t, config(), func(client *zei.Client, u string) (*http.Response, error) {
		data := make(url.Values)
		data.Set("testing", "true")
		return client.PostForm(u, data)
	})
}

func TestBasicAuth(t *testing.T) {
	const username = "zei"
	const password = "password"
	config := config()
	config.SetBasicAuth(username, password)

	testClientExt(t, config, func(client *zei.Client, url string) (*http.Response, error) {
		return client.Get(url)
	}, func(t testing.TB, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			auths := strings.SplitN(auth, ":", 2)
			if l := len(auths); l != 2 {
				t.Fatalf("Expected auths to have a length of %d, got %d instead", 2, l)
			}
			// Simple *insecure* comparison.
			// We don't care about security here since this is being executed in test.
			// In the real world, comparison should be done with the [crypto/subtle] package.
			if auths[0] != username {
				t.Logf("Expected username to be %q, got %q instead", username, auths[0])
			}
			if auths[1] != password {
				t.Logf("Expected password to be %q, got %q instead", password, auths[1])
			}
		} else {
			t.Logf("Expected header %q to not be empty", "Authorization")
			t.Fail()
		}
	})
}

func TestBearerAuth(t *testing.T) {
	const token = "bearer_token"
	config := config()
	config.SetBearerAuth(token)

	testClientExt(t, config, func(client *zei.Client, url string) (*http.Response, error) {
		return client.Get(url)
	}, func(t testing.TB, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			auths := strings.SplitN(auth, " ", 2)
			if l := len(auths); l != 2 {
				t.Fatalf("Expected auths to have a length of %d, got %d instead", 2, l)
			}
			if auths[0] != "Bearer" {
				t.Logf("Expected Bearer Authorization, got %q instead", auths[0])
				t.Fail()
			}
			// Simple *insecure* comparison.
			// We don't care about security here since this is being executed in test.
			// In the real world, comparison should be done with the [crypto/subtle] package.
			if auths[1] != token {
				t.Logf("Expected token to be %q, got %q instead", token, auths[1])
				t.Fail()
			}
		} else {
			t.Logf("Expected header %q to not be empty", "Authorization")
			t.Fail()
		}
	})
}

func config() *zei.Config {
	return &zei.Config{
		UserAgent: userAgent,
	}
}

type testerFn func(client *zei.Client, url string) (*http.Response, error)
type reqChecker func(t testing.TB, r *http.Request)

func testClient(t testing.TB, conf *zei.Config, tester testerFn) {
	t.Helper()
	testClientExt(t, conf, tester, nil)
}

func testClientExt(t testing.TB, conf *zei.Config, tester testerFn, checker reqChecker) {
	t.Helper()

	client := newClient(t, conf)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != userAgent {
			t.Logf("Expected User-Agent to be %q, got %q instead", userAgent, ua)
			t.Fail()
		}

		if r.Method == http.MethodPost {
			contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				t.Fatalf("Expected no error, got %v instead", err)
			}

			switch contentType {
			case "application/x-www-form-urlencoded":
				if err = r.ParseForm(); err != nil {
					t.Fatalf("Expected no error, got %v instead", err)
				}
				if pv := r.PostFormValue("testing"); pv != "true" {
					t.Logf("Expected %q to be %q, got %q instead", "testing", "true", pv)
					t.Fail()
				}

			case "application/json":
				type exp struct {
					Testing bool `json:"testing"`
				}
				var data exp
				err = json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					t.Fatalf("Expected no error, got %v instead", err)
				}

				if !data.Testing {
					t.Logf("Expected %q to be %t, got %t instead", "testing", true, data.Testing)
				}
			}
		}

		if checker != nil {
			checker(t, r)
		}

		w.WriteHeader(http.StatusOK)
	}))

	_, err := tester(client, srv.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v instead", err)
	}
}

func newClient(t testing.TB, conf *zei.Config) *zei.Client {
	t.Helper()

	return zei.New(conf)
}
