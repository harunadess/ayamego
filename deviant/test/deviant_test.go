package booru

import (
	"fmt"
	"testing"

	"github.com/jordanjohnston/ayamego/deviant"
)

func TestAuth(t *testing.T) {
	deviant.Auth()
}

func TestSearch(t *testing.T) {
	searchTerms := "anime"
	found, result := deviant.Search(searchTerms)
	if found {
		fmt.Println(result)
	} else {
		t.Error("Didn't find any results")
	}
}
