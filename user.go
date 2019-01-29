package trakt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (t *Client) Calendar(date string, days int) (episodes []ShowEpisode, err error) {
	err = t.get(fmt.Sprintf("https://api.trakt.tv/calendars/my/shows/%s/%d", date, days), &episodes)
	return episodes, err
}

// TODO Add options for methods for sorting etc
func (t *Client) WatchlistEpisodes() (episodes []MetadataShowEpisode, err error) {
	err = t.get(fmt.Sprintf("https://api.trakt.tv/sync/watchlist/episodes"), &episodes)
	return episodes, err
}

// TODO Add options for methods for sorting etc
func (t *Client) WatchlistMovies() (movies []MetadataMovie, err error) {
	err = t.get(fmt.Sprintf("https://api.trakt.tv/sync/watchlist/movies"), &movies)
	return movies, err
}

func (t *Client) get(url string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	return err
}
