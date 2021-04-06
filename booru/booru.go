package booru

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	errors "github.com/jordanjohnston/ayamego/util/errors"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

const https string = "https://"

const domain string = "danbooru.donmai.us/"
const loginPath string = "profile.json"
const postsPath string = "posts.json?"
const largeImage string = "large_file_url"
const source string = "source"
const pixivCDN string = "i.pximg.net"
const basePixivURL string = "https://pixiv.net/en/artworks/"

var booruSecrets struct {
	Login string `json:"login"`
	Key   string `json:"key"`
}

func init() {
	fPath := parseArgs()
	readConfig(fPath)

	// fPath := "/home/pi/Documents/projects/golang/src/github.com/jordanjohnston/danbooru.json"
	// readConfig(&fPath)
}

func parseArgs() *string {
	fPath := flag.String("booru", "", "path to booru .json file")
	flag.Parse()

	if *fPath == "" {
		errors.FatalErrorHandler("parseArgs: ", fmt.Errorf("%v", "no -booru specified"))
	}

	return fPath
}

func readConfig(fPath *string) {
	const maxJSONBytes int = 256

	file, err := os.Open(*fPath)
	defer file.Close()
	errors.FatalErrorHandler("readConfig: ", err)

	data := make([]byte, maxJSONBytes)
	count, err := file.Read(data)
	errors.FatalErrorHandler("readConfig: ", err)

	err = json.Unmarshal(data[:count], &booruSecrets)
	errors.FatalErrorHandler("readConfig: ", err)
}

// Search finds images based on the search args
func Search(searchArgs string) (bool, []string) {
	args := strings.Split(searchArgs, ",")
	for i := range args {
		args[i] = strings.Trim(args[i], " ")
		args[i] = strings.Join(strings.Split(args[i], " "), "_")
	}
	found, images := searchForTags(args)

	return found, images
}

func searchForTags(tags []string) (bool, []string) {
	searchString := makeURL(postsPath)
	tagsParams := convertTagsToParams(tags)

	searchURL, err := url.Parse((searchString + tagsParams))
	errors.StandardErrorHandler("booru.searchForTags", err)
	logger.Info("searching for: ", tags)

	resp, err := http.Get(searchURL.String())
	errors.StandardErrorHandler("booru.searchForTags", err)
	defer resp.Body.Close()

	results := parseBody(resp)
	found := (len(results) > 0)
	if !found {
		return found, make([]string, 0)
	}

	randomItem := rand.Intn(len(results))
	logger.Info("Got ", len(results), " results")
	item := results[randomItem].(map[string]interface{})

	images := pluckImages(item)

	return found, images
}

func makeURL(path string) string {
	baseURL := fmt.Sprintf("%s%s:%s@%s", https, booruSecrets.Login, booruSecrets.Key, domain)
	baseURL = fmt.Sprintf("%s%s", baseURL, path)
	return baseURL
}

func convertTagsToParams(tags []string) string {
	return "tags=" + strings.Join(tags, "+") + "+"
}

func parseBody(resp *http.Response) []interface{} {
	body, err := io.ReadAll(resp.Body)
	errors.StandardErrorHandler("parseBody: ", err)

	if string(body)[:15] == "<!doctype html>" {
		logger.Error("booru.parseBody: ", "API request failed")
		return make([]interface{}, 0)
	}

	var parsed []interface{}
	json.Unmarshal(body, &parsed)

	return parsed
}

func pluckImages(item map[string]interface{}) []string {
	images := make([]string, 2)

	images[0] = item[largeImage].(string)
	images[1] = item[source].(string)

	if strings.Contains(source, pixivCDN) {
		urlParts := strings.Split(source, "/")
		imageFile := urlParts[len(urlParts)-1]
		imageID := strings.Split(imageFile, "_")[0]
		images[1] = basePixivURL + imageID
	}

	return images
}
