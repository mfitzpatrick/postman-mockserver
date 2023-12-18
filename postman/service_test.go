package postman

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestRouter_URLNotFound(t *testing.T) {
	body := ""
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	postmanRouter(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestRouter_LogsRecordedAsJSON(t *testing.T) {
	var logBuffer strings.Builder
	log.Logger = log.Output(&logBuffer)
	defer func() { log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) }()
	loglevel := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	defer zerolog.SetGlobalLevel(loglevel)

	body := `{"item": "test"}`
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	postmanRouter(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, logBuffer.String(), "request")
	logLines := strings.Split(logBuffer.String(), "\n")
	found := false
	for i, line := range logLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		js := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &js)
		assert.Nilf(t, err, "log line %d should be valid json: %s", i, line)
		if js["message"] == "request" {
			found = true
			assert.Equalf(t, "post/test", js["path"], "JSON %s", line)
			assert.Equalf(t, `{"item": "test"}`, js["body"], "JSON %s", line)
		}
	}
	assert.Truef(t, found, "log line should contain request: %s", logBuffer.String())
}
