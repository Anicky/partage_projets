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

func TestPostComment(testing *testing.T) {
	router := InitTest()

	comment := map[string]interface{}{
		"project_id": 1,
		"content":    "Test comment",
	}

	data, err := json.Marshal(comment)
	if err != nil {
		log.Fatal("Unable to marshal data: ", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/comments/", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Unable to create request: ", err)
	}

	request.Header.Set("Content-Type", "application/json")

	AuthenticateUser(request)

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(testing, http.StatusCreated, response.Code)

	body := response.Body.String()

	assert.Contains(testing, body, "Test comment")
}
