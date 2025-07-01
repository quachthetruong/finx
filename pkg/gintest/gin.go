package gintest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func GetTestContext() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	return ctx, engine, recorder
}

func ExtractBodyAsString(body io.ReadCloser) string {
	bodyStr, _ := io.ReadAll(body)
	return *(*string)(unsafe.Pointer(&bodyStr))
}

func ExtractBody(body io.ReadCloser) []byte {
	bodyBytes, _ := io.ReadAll(body)
	return bodyBytes
}

func MustMakeRequest(method, path string, body any) *http.Request {
	payload, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept-Language", "en")
	return req
}
