package openWeather

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myFakeService func(*http.Request) (*http.Response, error)

func (s myFakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

func TestGetTodayHolidayByCountry(t *testing.T) {
	client := &http.Client{
		Transport: myFakeService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(`{"coord": {"lon": 33,"lat": 22}}`)),
			}, nil
		}),
	}

	repo := &factRepository{
		client: client,
	}

	got := repo.GetWeatherFromApi()

	assert.Equal(t, got.Coord.Lon, 33.0)
	assert.Equal(t, got.Coord.Lat, 22.0)
}
