package audio

import (
	"strings"
)

func ProcessQuery(query string) []*Song {
	if strings.Contains(query, "watch?v=") {
		return []*Song{
			NewSongByURL(query),
		}
	} else {
		return []*Song{
			NewSongBySearch(query),
		}
	}
}
