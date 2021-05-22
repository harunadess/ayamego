package imageresults

// ImageResults contains the URL and thumbnail of an image
type ImageResults struct {
	ImageURL  string
	Thumbnail string
}

// SearchResults contains the search results for the image
type SearchResults struct {
	Title  string
	Images ImageResults
	Tags   string
}
