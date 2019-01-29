package trakt

import (
	"net/http"
)

const ApiUrl = "https://api.trakt.tv"

type Client struct {
	Auth       *AuthClient
	HttpClient *http.Client
}

type AuthClient struct {
	ClientID     string
	ClientSecret string
	Token        *Token
	HttpClient   *http.Client
}

func NewClient(clientID, clientSecret string, token Token, httpClient *http.Client, onTokenRefresh func(t Token)) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	nonAuthClient := &http.Client{Transport: httpClient.Transport}

	source := RefreshTokenSource{
		HttpClient:   nonAuthClient,
		RefreshToken: token,
		ClientSecret: clientSecret,
		ClientID:     clientID,
	}

	httpClient.Transport = &Transport{
		Token:      token,
		Base:       httpClient.Transport,
		Source:     &source,
		ClientID:   clientID,
		OnTokenGet: onTokenRefresh,
	}

	return &Client{
		Auth: &AuthClient{
			HttpClient:   nonAuthClient,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		HttpClient: httpClient,
	}
}
