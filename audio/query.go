package audio

import (
	"errors"
	"strings"
)

func ProcessQuery(query string) ([]*Song, error) {
	if strings.Contains(query, "spotify") && strings.Contains(query, "/playlist/") {
		split := strings.Split(query, "/playlist/")
		if len(split) < 2 {
			return nil, errors.New("invalid spotify playlist URL")
		}
		playlistID := split[1]
		return GetSpotifyPlaylistSongs(playlistID)
	} else if strings.Contains(query, "watch?v=") {
		return []*Song{
			NewSongByURL(query),
		}, nil
	} else {
		return []*Song{
			NewSongBySearch(query),
		}, nil
	}
}
