package models

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParsePaginationQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		wantPage   int
		wantLimit  int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "default values",
			query:     "",
			wantPage:  1,
			wantLimit: 10,
			wantErr:   false,
		},
		{
			name:      "valid values",
			query:     "?page=2&limit=20",
			wantPage:  2,
			wantLimit: 20,
			wantErr:   false,
		},
		{
			name:       "invalid page (zero)",
			query:      "?page=0",
			wantErr:    true,
			wantErrMsg: "invalid page parameter",
		},
		{
			name:       "invalid page (negative)",
			query:      "?page=-1",
			wantErr:    true,
			wantErrMsg: "invalid page parameter",
		},
		{
			name:       "invalid page (not a number)",
			query:      "?page=abc",
			wantErr:    true,
			wantErrMsg: "invalid page parameter",
		},
		{
			name:       "invalid limit (zero)",
			query:      "?limit=0",
			wantErr:    true,
			wantErrMsg: "invalid limit parameter",
		},
		{
			name:       "limit too large",
			query:      "?limit=101",
			wantErr:    true,
			wantErrMsg: "limit too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/"+tt.query, nil)
			c.Request = req

			page, limit, err := ParsePaginationQuery(c)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPage, page)
				assert.Equal(t, tt.wantLimit, limit)
			}
		})
	}
}
