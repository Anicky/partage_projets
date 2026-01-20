package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjects(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/projects/", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test project 1")
	assert.Contains(testing, body, "Test comment on project 1")
	assert.Contains(testing, body, "Test project 2")
}

func TestGetProject(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodGet, "/projects/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test project 1")
	assert.Contains(testing, body, "Test comment on project 1")
	assert.NotContains(testing, body, "Test project 2")
}

func TestPostProject(testing *testing.T) {
	router := InitTest()

	project := map[string]interface{}{
		"name":        "Test project 3",
		"description": "Test description 3",
		"skills":      []string{"Go", "Testing", "SQLite"},
	}

	data, err := json.Marshal(project)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/projects/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test project 3")
	assert.Contains(testing, body, "Test description 3")
	assert.Contains(testing, body, "Testing")
}

func TestPutProject(testing *testing.T) {
	router := InitTest()

	update := map[string]interface{}{
		"name": "Updated project 1",
	}

	data, err := json.Marshal(update)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/projects/1", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Updated project")
}

func TestDeleteProject(testing *testing.T) {
	router := InitTest()

	request, err := http.NewRequest(http.MethodDelete, "/projects/1", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(request)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusOK, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Project deleted successfully.")
}

func TestLikeProject(testing *testing.T) {
	router := InitTest()

	// Test like
	requestLike, err := http.NewRequest(http.MethodPut, "/projects/1/like", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(requestLike)

	responseLike := httptest.NewRecorder()

	router.ServeHTTP(responseLike, requestLike)

	assert.Equal(testing, http.StatusOK, responseLike.Code)

	bodyLike := responseLike.Body.String()

	assert.Contains(testing, bodyLike, "Project liked successfully.")

	// Test unlike
	requestUnlike, err := http.NewRequest(http.MethodPut, "/projects/1/like", nil)
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	AuthenticateUser(requestUnlike)

	responseUnlike := httptest.NewRecorder()

	router.ServeHTTP(responseUnlike, requestUnlike)

	assert.Equal(testing, http.StatusOK, responseUnlike.Code)

	bodyUnlike := responseUnlike.Body.String()

	assert.Contains(testing, bodyUnlike, "Project unliked successfully.")
}
