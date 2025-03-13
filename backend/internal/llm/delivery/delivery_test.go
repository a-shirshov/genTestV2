package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"gentest/internal/llm/delivery"
	mockLLM "gentest/internal/llm/mock"
	"gentest/internal/llm/models"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateTest_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockLLM.NewMockUsecase(ctrl)
	delivery := delivery.NewDelivery(mockUsecase)

	router := gin.Default()
	router.POST("/generate-test", delivery.GenerateFirstTest)

	// Invalid JSON
	reqBody := bytes.NewBufferString(`{"prompt": "test"`) // Missing closing brace
	req, _ := http.NewRequest("POST", "/generate-test", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

// Custom response writer to support CloseNotify
type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func TestGenerateTest_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockLLM.NewMockUsecase(ctrl)
	delivery := delivery.NewDelivery(mockUsecase)

	router := gin.Default()
	router.POST("/generate-test", delivery.GenerateFirstTest)

	// Mock Ollama response
	ollamaResponses := []models.OllamaResponse{
		{Response: "Hello", Done: false},
		{Response: " world", Done: false},
		{Response: "!", Done: true},
	}

	var responseBody bytes.Buffer
	for _, resp := range ollamaResponses {
		jsonData, _ := json.Marshal(resp)
		responseBody.Write(jsonData)
		responseBody.WriteString("\n") // Ensure newline between JSON objects
	}

	mockResp := &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&responseBody),
		},
	}

	mockUsecase.EXPECT().
		GenerateResponse(gomock.Any()).
		Return(mockResp, nil)

	reqBody := bytes.NewBufferString(`{"prompt": "Hello"}`)
	req, _ := http.NewRequest("POST", "/generate-test", reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Use custom response writer
	w := &closeNotifyingRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closed:           make(chan bool, 1),
	}
	
	router.ServeHTTP(w, req)

	// Decode chunked encoding (if still needed)
	decoded, err := io.ReadAll(httputil.NewChunkedReader(w.Body))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello world!", string(decoded))
}

func TestGenerateTest_OllamaFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockLLM.NewMockUsecase(ctrl)
	delivery := delivery.NewDelivery(mockUsecase)

	router := gin.Default()
	router.POST("/generate-test", delivery.GenerateFirstTest)

	// Set expectations to return error
	mockUsecase.EXPECT().
		GenerateResponse(gomock.Any()).
		Return(nil, errors.New("Ollama API failed"))

	reqBody := bytes.NewBufferString(`{"prompt": "Hello"}`)
	req, _ := http.NewRequest("POST", "/generate-test", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get response from Ollama")
}