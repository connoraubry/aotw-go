package server

import (
	"fmt"
	"testing"
)

func TestParseSubmission(t *testing.T) {

	title := "Midnight Marauders"
	artist := "A Tribe Called Quest"

	queryString := fmt.Sprintf("%s (%s)", title, artist)
	albumInfo, _ := ParseSubmission(queryString)

	if albumInfo.Title != title {
		t.Errorf("Got %v. Wanted %v", albumInfo.Title, title)
	}

	if albumInfo.Artist != artist {
		t.Errorf("Got %v. Wanted %v", albumInfo.Artist, artist)
	}
}

func TestParseSubmission2(t *testing.T) {
	title := "Red (Taylor's Version) (From the Vault)"
	artist := "Taylor Swift"

	queryString := fmt.Sprintf("%s (%s)", title, artist)
	albumInfo, _ := ParseSubmission(queryString)

	if albumInfo.Title != title {
		t.Errorf("Got %v. Wanted %v", albumInfo.Title, title)
	}

	if albumInfo.Artist != artist {
		t.Errorf("Got %v. Wanted %v", albumInfo.Artist, artist)
	}

}
