package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestValidator(t *testing.T) {

	t.Run("run 1", func(t *testing.T) {
		body := `{
			"projectName" : "test-project",
			"development" : {
				"values" : {
					"image" : {
						"tag" : "test-project-dev-main"
					}
				}
			}
		}`
		client := mockHttpClient{
			InputBody: body,
			InputCode: 200,
		}
		if bVal := validator(&client, "test-token", "test-project", "test-project-dev-main", "development"); !bVal {
			t.Errorf("should return true")
		}
	})

	t.Run("run 2", func(t *testing.T) {
		body := `{
			"projectName" : "test-project",
			"production" : {
				"values" : {
					"image" : {
						"tag" : "test-project-dev-main"
					}
				}
			}
		}`
		client := mockHttpClient{
			InputBody: body,
			InputCode: 200,
		}
		if bVal := validator(&client, "test-token", "test-project", "test-project-dev-main", "production"); !bVal {
			t.Errorf("should return true")
		}
	})

	t.Run("run 3", func(t *testing.T) {
		body := `{
			"projectName" : "test-project",
			"development" : {
				"values" : {
					"image" : {
						"tag" : "test-project-dev-main"
					}
				}
			}
		}`
		client := mockHttpClient{
			InputBody: body,
			InputCode: 404,
		}
		if bVal := validator(&client, "test-token", "test-project", "test-project-dev-main", "development"); bVal {
			t.Errorf("should return false")
		}

	})
	t.Run("run 4", func(t *testing.T) {
		body := ``
		client := mockHttpClient{
			InputBody:  body,
			InputCode:  404,
			InputError: errors.New("intentional"),
		}
		if bVal := validator(&client, "test-token", "test-project", "test-project-dev-main", "development"); bVal {
			t.Errorf("should return false")
		}

	})

}

type mockHttpClient struct {
	InputBody  string
	InputCode  int
	InputError error
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	res := new(http.Response)
	res.Body = io.NopCloser(strings.NewReader(m.InputBody))
	res.StatusCode = m.InputCode
	return res, m.InputError
}
