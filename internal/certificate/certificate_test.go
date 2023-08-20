package certificate

import (
	"bytes"
	"encoding/json"
	"gcipher/internal/server/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCertificateRequest(t *testing.T) {
	// Prepare a sample request
	requestData := api.RequestData{
		Applicant: "John Doe",
		CSR:       "base64_encoded_csr_here",
		Lifetime:  365,
		Type:      "server",
	}
	authData := api.Auth{
		Username: "testuser",
		Password: "testpassword",
	}
	request := api.Request{
		Data: requestData,
		Auth: authData,
	}

	// Serialize the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal("Failed to marshal JSON request:", err)
	}

	// Create a mock HTTP request
	req, err := http.NewRequest("POST", "/certificate/request", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal("Failed to create HTTP request:", err)
	}

	// Create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	HandleCertificateRequest(rr, req)

	// Check the response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	// Decode the response
	var response api.Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if !response.Success {
		t.Error("Expected success response")
	}
}
