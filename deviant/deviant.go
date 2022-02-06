package deviant

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jordanjohnston/ayamego/imageresults"
	"github.com/jordanjohnston/ayamego/util/envflags"
	errors "github.com/jordanjohnston/ayamego/util/errors"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

const https string = "https://"

const domain string = "www.deviantart.com/"
const authPath string = "oauth2/token"
const apiVersion string = "api/v1/oauth2/"
const browseTags string = "browse/newest?with_session=false&mature_content=true&"
const writePerm fs.FileMode = 0666

var deviantSecrets struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type oAuth2Token struct {
	ExpiresIn   float64 `json:"expires_in"`
	Status      string  `json:"status"`
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	expireDate  time.Time
}

type urlBuilder func(args ...string) string

type imageValues struct {
	title      string
	tags       string
	contentSrc string
}

var authToken oAuth2Token

func init() {
	fPath := envflags.DeviantPath
	readConfig(fPath)
}

func parseArgs() *string {
	fPath := flag.String("deviant", "", "path to booru .json file")
	flag.Parse()

	if *fPath == "" {
		errors.FatalErrorHandler("parseArgs: ", fmt.Errorf("%v", "no -deviant specified"))
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

	err = json.Unmarshal(data[:count], &deviantSecrets)
	errors.FatalErrorHandler("readConfig: ", err)
}

// Auth authenticates with oauth2
func Auth() {
	result, err := makeGetRequest(buildAuthURL)
	errors.StandardErrorHandler("deviant.auth", err)

	if _, ok := result["access_token"]; ok {
		authToken.ExpiresIn = result["expires_in"].(float64)
		authToken.Status = result["status"].(string)
		authToken.AccessToken = result["access_token"].(string)
		authToken.TokenType = result["token_type"].(string)
		authToken.expireDate = time.Now().Add(time.Duration(authToken.ExpiresIn))
		logger.Info("got access token: ", authToken)
	} else {
		logger.Info("didn't get access token, response: ", result)
	}
}

func makeGetRequest(makeURL urlBuilder, args ...string) (map[string]interface{}, error) {
	requestURL, err := url.Parse(makeURL(args...))
	errors.StandardErrorHandler(requestURL.String(), err)

	resp, err := http.Get(requestURL.String())
	errors.StandardErrorHandler(makeURL(args...), err)
	defer resp.Body.Close()

	return parseBodySingle(resp), err
}

func buildAuthURL(args ...string) string {
	baseURL := https + domain + authPath
	authURL := fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s",
		baseURL, deviantSecrets.ClientID, deviantSecrets.ClientSecret)
	return authURL
}

func parseBodySingle(resp *http.Response) map[string]interface{} {
	body, err := io.ReadAll(resp.Body)
	errors.StandardErrorHandler("parseBody: ", err)

	var parsed interface{}
	json.Unmarshal(body, &parsed)

	parsedBody := parsed.(map[string]interface{})

	if v, ok := parsedBody["error"]; ok {
		errors.StandardErrorHandler("got error", fmt.Errorf("%w", v))
	}

	return parsedBody
}

// Search for values
// todo: handle literature..
func Search(searchTerms string) (bool, imageresults.SearchResults) {
	if time.Since(authToken.expireDate) > time.Hour {
		Auth()
	}
	tags := strings.Split(searchTerms, ", ")

	result, err := makeGetRequest(makeDeviationURL, tags...)
	errors.StandardErrorHandler("deviant.Search", err)

	writeToFile(result)

	if value, ok := result["results"]; ok {
		return getRandomResult(value.([]interface{}))
	}

	return false, imageresults.SearchResults{}
}

func writeToFile(result map[string]interface{}) {
	s, err := json.Marshal(result)
	if err != nil {
		errors.StandardErrorHandler("deviant.writeToFile", err)
	}
	// dump to file for now - todo: fix this at some point, so literature is handled either appropriately, or we just ignore it
	os.WriteFile("./result.json", s, writePerm)
}

func makeDeviationURL(searchTerms ...string) string {
	baseURL := fmt.Sprintf("%s%s%s", https, domain, apiVersion)
	baseURL = fmt.Sprintf("%s%s", baseURL, browseTags)

	for _, v := range searchTerms {
		baseURL += fmt.Sprintf("q=%s&", v)
	}
	baseURL = baseURL[:len(baseURL)-1]
	baseURL = fmt.Sprintf("%s&access_token=%s", baseURL, authToken.AccessToken)

	return baseURL
}

func getRandomResult(r []interface{}) (bool, imageresults.SearchResults) {
	results := make([]map[string]interface{}, len(r))
	for i, v := range r {
		results[i] = make(map[string]interface{})
		results[i] = v.(map[string]interface{})
	}
	rand.Seed(time.Now().Unix())
	if len(results) == 0 {
		return false, imageresults.SearchResults{}
	}
	i := rand.Intn(len(results))

	v := results[i]

	if v != nil {
		imageValues := getValuesFromMap(v)

		sr := imageresults.SearchResults{
			Title: imageValues.title,
			Images: imageresults.ImageResults{
				ImageURL:  imageValues.contentSrc,
				Thumbnail: imageValues.contentSrc,
			},
			Tags: imageValues.tags,
		}
		return true, sr
	}

	return false, imageresults.SearchResults{}
}

func getValuesFromMap(m map[string]interface{}) imageValues {
	titleI, ok := m["title"]
	if !ok {
		logger.Error("no title in response ", m)
		return imageValues{}
	}
	title := titleI.(string)

	categoryPathI, ok := m["category_path"]
	if !ok {
		logger.Error("no category_path in response ", m)
		return imageValues{}
	}
	categoryPath := categoryPathI.(string)

	contentI, ok := m["content"]
	if !ok {
		logger.Error("no category_path in response ", m)
		return imageValues{}
	}
	content := contentI.(map[string]interface{})

	tagsSplit := strings.Split(categoryPath, "/")
	tags := strings.Join(tagsSplit, ", ")

	return imageValues{
		title:      title,
		contentSrc: content["src"].(string),
		tags:       tags,
	}
}

func makeURL(path string) string {
	baseURL := fmt.Sprintf("%s%s:%s@%s", https, deviantSecrets.ClientID, deviantSecrets.ClientSecret, domain)
	baseURL = fmt.Sprintf("%s%s", baseURL, path)
	return baseURL
}
