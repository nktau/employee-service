package httplayer

import (
	"github.com/employee-service/internal/applayer"
	"github.com/employee-service/internal/storagelayer"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestGetUrlByShorter(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		url         string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
				url:         "google.com",
			},
			path: "/123",
		},
	}

	storage := storagelayer.New()
	app := applayer.New(storage)
	api := New(app)
	ts := httptest.NewServer(api.mux)

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.R().Get(ts.URL + tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.url, result.String())
		})
	}
}
