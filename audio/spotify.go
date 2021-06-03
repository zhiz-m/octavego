package audio

import (
	"context"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func GetSpotifyPlaylistSongs(playlistID string) ([]*Song, error) {
	config := &clientcredentials.Config{
		ClientID:     "5f573c9620494bae87890c0f08a60293",
		ClientSecret: "212476d9b0f3472eaa762d90b19b0ba8",
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		return nil, err
	}
	client := spotify.Authenticator{}.NewClient(token)
	playlist, err := client.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}
	songs := make([]*Song, 0)
	for _, v := range playlist.Tracks.Tracks {
		track := v.Track
		title := track.Name
		artist := track.Artists[0].Name
		duration := track.Duration

		song := NewSongBySearch(YoutubeSearchQuery(title, artist))
		song.SetMetadata(title, artist, duration/1000)

		songs = append(songs, song)
	}
	return songs, nil
}
