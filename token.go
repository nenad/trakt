package trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// Time to accommodate for slow requests or server time inequalities
	expiryDelta = 10 * time.Second
)

// Token represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// Most users of this package should not access fields of Token
// directly. They're exported mostly for use by related packages
// implementing derivative OAuth2 flows.
type (
	Token struct {
		// AccessToken is the Token that authorizes and authenticates
		// the requests.
		AccessToken string `json:"access_token"`

		// TokenType is the type of Token.
		TokenType string `json:"token_type"`

		// RefreshToken is a Token that's used by the application
		// (as opposed to the user) to refresh the access Token
		// if it expires.
		RefreshToken string `json:"refresh_token"`

		// ExpiresIn is the expiration time of the access Token.
		ExpiresIn int64 `json:"expires_in"`

		// Scope is the authorized scope of the Token
		Scope string `json:"scope"`

		// CreatedAt is a time when the Token was generated
		CreatedAt int64 `json:"created_at"`
	}

	// TokenProvider provides a method for returning a Token
	TokenProvider interface {
		// Token returns a Token or an error.
		Token() (Token, error)
	}

	// RefreshTokenSource is TokenProvider that refreshes the Token
	RefreshTokenSource struct {
		HttpClient   *http.Client
		ClientID     string
		ClientSecret string
		RefreshToken Token
	}
)

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	expiryDate := time.Unix(t.CreatedAt, 0).Add(time.Second * time.Duration(t.ExpiresIn)).Add(-expiryDelta)
	return t != nil && t.AccessToken != "" && time.Now().Before(expiryDate)
}

func (s *RefreshTokenSource) Token() (t Token, err error) {
	body, _ := json.Marshal(map[string]string{
		"client_id":     s.ClientID,
		"client_secret": s.ClientSecret,
		"refresh_token": s.RefreshToken.RefreshToken,
		"grant_type":    "refresh_token",
		"redirect_uri":  "urn:ietf:wg:oauth:2.0:oob",
	})

	req, err := http.NewRequest("POST", AuthTokenUrl, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return t, err
	}

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return t, err
	}

	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("could not get Token, status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}

	s.RefreshToken = t

	return t, nil
}
