package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const AuthCodeUrl = "https://api.trakt.tv/oauth/device/code"
const AuthCodeTokenUrl = "https://api.trakt.tv/oauth/device/token"
const AuthTokenUrl = "https://api.trakt.tv/oauth/token"

type Code struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func (auth *AuthClient) Code() (code Code, err error) {
	req, err := http.NewRequest(
		"POST",
		AuthCodeUrl,
		bytes.NewBufferString(fmt.Sprintf(`{"client_id":"%s"}`, auth.ClientID)),
	)
	req.Header.Add("Content-Type", "application/json")

	resp, err := auth.HttpClient.Do(req)
	if err != nil {
		return code, err
	}

	if resp.StatusCode != http.StatusOK {
		return code, fmt.Errorf("could not get code, status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&code); err != nil {
		return code, err
	}

	return code, nil
}

func (auth *AuthClient) TryToken(code Code) (token Token, err error) {
	params := map[string]string{
		"client_id":     auth.ClientID,
		"client_secret": auth.ClientSecret,
		"code":          code.DeviceCode,
	}

	body, _ := json.Marshal(params)

	req, err := http.NewRequest(
		"POST",
		AuthCodeTokenUrl,
		bytes.NewBuffer(body),
	)
	req.Header.Add("Content-Type", "application/json")

	resp, err := auth.HttpClient.Do(req)
	if err != nil {
		return token, err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO Explain status code
		return token, fmt.Errorf("could not get Token, status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return token, err
	}

	return token, nil
}

type Result struct {
	Token Token
	Err   error
}

func (auth *AuthClient) PollToken(ctx context.Context, code Code) chan Result {
	done := make(chan Result)
	go func() {
		timer := time.NewTicker(time.Second * time.Duration(code.Interval))
		ticks := code.ExpiresIn/code.Interval - 1

		fmt.Printf("Ticks: %d\n", ticks)

	loop:
		for {
			select {
			case <-timer.C:
				t, err := auth.TryToken(code)
				if err == nil {
					done <- Result{t, err}
					timer.Stop()
					break
				}

				ticks--
				if ticks == 0 {
					timer.Stop()
					done <- Result{t, fmt.Errorf("poller: token timeout expired")}
					break
				}

			case <-ctx.Done():
				timer.Stop()
				done <- Result{Token{}, fmt.Errorf("poller: polling cancelled")}
				break loop
			}
		}
	}()

	return done
}
