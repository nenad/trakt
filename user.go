package trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// TODO Add options for methods for sorting etc
func (t *Client) WatchlistSeasons() (seasons []MetadataSeason, err error) {
	err = t.get(fmt.Sprintf("https://api.trakt.tv/sync/watchlist/seasons"), &seasons)
	return seasons, err
}

// TODO Add options for methods for sorting etc
func (t *Client) WatchlistShows() (shows []MetadataShow, err error) {
	err = t.get(fmt.Sprintf("https://api.trakt.tv/sync/watchlist/shows"), &shows)
	return shows, err
}

func (t *Client) RemoveFromWatchlist(m FullMetadata) error {
	return t.post(fmt.Sprintf("https://api.trakt.tv/sync/watchlist/remove"), m, nil)
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

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)
	}

	return err
}

func (t *Client) post(url string, body interface{}, response interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error while removing from watchlist: %s", string(body))
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(&response)
	}

	return err
}
