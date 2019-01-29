package trakt_test

import (
	"bytes"
	"github.com/bmizerany/assert"
	"github.com/nenadstojanovikj/trakt"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestRoundTripFunc func(req *http.Request) *http.Response

func (rt TestRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req), nil
}

func NewTestClient(fn TestRoundTripFunc, token *trakt.Token) *trakt.AuthClient {
	c := &http.Client{
		Transport: fn,
	}
	return &trakt.AuthClient{HttpClient: c, ClientID: "client_id", ClientSecret: "client_secret", Token: token}
}

func TestClientAuthCode(t *testing.T) {
	authClient := NewTestClient(func(req *http.Request) *http.Response {
		reqBody, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"client_id":"client_id"}`, string(reqBody))

		respBody := `{"device_code":"1234567890", "user_code": "AAABBB123",
"verification_url": "https://trakt.tv/activate", "expires_in": 600, "interval": 5}`

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(respBody)),
		}
	}, nil)

	code, err := authClient.Code()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, trakt.Code{
		DeviceCode:      "1234567890",
		ExpiresIn:       600,
		Interval:        5,
		UserCode:        "AAABBB123",
		VerificationURL: "https://trakt.tv/activate",
	}, code)
}
