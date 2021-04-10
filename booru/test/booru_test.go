package booru

import (
	"fmt"
	"testing"

	"github.com/jordanjohnston/ayamego/booru"
)

func TestSearch(t *testing.T) {
	tags := "thighhighs, seifuku"
	found, result := booru.Search(tags)
	if found {
		fmt.Println(result)
	} else {
		t.Error("Didn't find any results")
	}
}
