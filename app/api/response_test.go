package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type failingResponseWriter struct {
	header http.Header
	code   int
}

func (w *failingResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *failingResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
}

func (w *failingResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("forced write error")
}

func TestOKResponse(t *testing.T) {

	type sampleResponse struct {
		Message string `json:"message"`
	}

	sample := sampleResponse{Message: "Success"}

	t.Run("succesful http200 json response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		OKResponse(recorder, sample)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"message":"Success"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("json response for a given http status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ErrorResponse(recorder, http.StatusInternalServerError, "Some error occurred")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"error":"Some error occurred"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestOKResponseEncodeErrorPath(t *testing.T) {
	t.Run("encode error branch", func(t *testing.T) {
		writer := &failingResponseWriter{}
		OKResponse(writer, map[string]any{"invalid": func() {}})

		assert.Equal(t, http.StatusInternalServerError, writer.code)
	})
}

func TestErrorResponseEncodeErrorPath(t *testing.T) {
	t.Run("encode error branch", func(t *testing.T) {
		writer := &failingResponseWriter{}
		ErrorResponse(writer, http.StatusBadRequest, "error")

		assert.Equal(t, http.StatusInternalServerError, writer.code)
	})
}
