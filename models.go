package trakt

import (
	"time"
)

type (
	Episode struct {
		Season    int         `json:"season"`
		Number    int         `json:"number"`
		Title     string      `json:"title"`
		Providers ProviderIDs `json:"ids"`
	}

	Show struct {
		Title     string      `json:"title"`
		Year      int         `json:"year"`
		Providers ProviderIDs `json:"ids"`
	}

	Movie struct {
		Title string      `json:"title"`
		Year  int         `json:"year"`
		Ids   ProviderIDs `json:"ids"`
	}

	ShowEpisode struct {
		FirstAired time.Time `json:"first_aired"`
		Episode    Episode   `json:"episode"`
		Show       Show      `json:"show"`
	}

	ProviderIDs struct {
		Trakt  int    `json:"trakt"`
		TVDB   int    `json:"tvdb"`
		IMDb   string `json:"imdb"`
		TMDB   int    `json:"tmdb"`
		TVRage int    `json:"tvrage"`
	}

	// Metadata coming from Watchlist endpoint
	MetadataMovie struct {
		Rank     int       `json:"rank"`
		ListedAt time.Time `json:"listed_at"`
		Type     string    `json:"type"`
		Movie    Movie     `json:"movie"`
	}

	MetadataShowEpisode struct {
		Rank     int       `json:"rank"`
		ListedAt time.Time `json:"listed_at"`
		Type     string    `json:"type"`
		Show     Show      `json:"show"`
		Episode  Episode   `json:"episode"`
	}

	MetadataShow struct {
		Rank     int       `json:"rank"`
		ListedAt time.Time `json:"listed_at"`
		Type     string    `json:"type"`
		Show     Show      `json:"show"`
	}
)
