package postman

import (
	"testing"

	. "github.com/dvincenz/postman-mockserver/common"
	"github.com/stretchr/testify/assert"
)

func TestGetMocks(t *testing.T) {
	mocks := getMocks([]respone{{
		Name:   "Mock Test",
		Status: "OK",
		Code:   200,
		Body:   `{"sample-body":true}`,
		// OriginalRequest.Method: "GET",
	}})
	assert.Equal(t, map[string]Mock{
		"": Mock{
			Method: "",
			Code:   200,
			Name:   "Mock Test",
			Body:   `{"sample-body":true}`,
			Header: make([]Header, 0),
		},
	}, mocks)
}
