package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NeekUP/roadmaps/core"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type openAuthenticator struct {
	providers     map[string]*oauth2.Config
	cache         core.DistributedCache
	cachePrefix   string
	baseReturnUrl string
}

func NewOpenAuthenticator(cache core.DistributedCache, baseReturnUrl string) core.OpenAuthenticator {
	return &openAuthenticator{
		providers:     map[string]*oauth2.Config{},
		cache:         cache,
		cachePrefix:   "oauth_",
		baseReturnUrl: baseReturnUrl,
	}
}

func (auth *openAuthenticator) AddProvider(providerName, clientId, clientSecret string, scope []string) {
	if _, ok := auth.providers[providerName]; ok {
		return
	}

	redirectUrl, err := url.Parse(auth.baseReturnUrl)
	if err != nil {
		panic(err)
	}

	parameters := url.Values{}
	parameters.Add("provider", providerName)
	redirectUrl.RawQuery = parameters.Encode()

	if ok, endpoint := auth.getEndpoint(providerName); ok {
		auth.providers[providerName] = &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl.String(),
			Scopes:       scope,
			Endpoint:     endpoint,
		}
	}
}

func (auth *openAuthenticator) HasProvider(providerName string) bool {
	_, exists := auth.providers[providerName]
	return exists
}

func (auth *openAuthenticator) getEndpoint(providerName string) (bool, oauth2.Endpoint) {
	switch providerName {
	case "facebook":
		return true, facebook.Endpoint
	case "github":
		return true, github.Endpoint
	default:
		return false, oauth2.Endpoint{}
	}
}
func (auth *openAuthenticator) LoginLink(providerName string) (string, error) {
	fakeName := "fake_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	link, err := auth.RegisterLink("fake_"+fakeName, providerName)
	return link, err
}

func (auth *openAuthenticator) RegisterLink(username, providerName string) (string, error) {

	var provider *oauth2.Config
	var ok bool
	if provider, ok = auth.providers[providerName]; !ok {
		return "", errors.New(fmt.Sprintf("auth provider %v not found.", providerName))
	}

	Url, err := url.Parse(provider.Endpoint.AuthURL)
	if err != nil {
		return "", err
	}

	state := uuid.New().String()
	auth.saveState(username, state)

	parameters := url.Values{}
	parameters.Add("client_id", provider.ClientID)
	parameters.Add("scope", strings.Join(provider.Scopes, " "))
	parameters.Add("redirect_uri", provider.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", state)
	Url.RawQuery = parameters.Encode()
	return Url.String(), nil
}

func (auth *openAuthenticator) Auth(providerName, state, code string) (name, email, openid string, err error) {
	var provider *oauth2.Config
	var ok bool
	if provider, ok = auth.providers[providerName]; ok {
		return "", "", "", errors.New(fmt.Sprintf("auth provider %v not found.", providerName))
	}

	username := auth.getUsernameByState(state)
	if username == "" {
		return "", "", "", errors.New("Name not found by state: " + state)
	}

	token, err := provider.Exchange(context.Background(), code)
	if err != nil {
		return "", "", "", err
	}

	resp, err := http.Get("https://graph.facebook.com/me?fields=id,email&access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	data := new(facebookResponse)
	err = decoder.Decode(data)
	defer resp.Body.Close()

	if err != nil {
		return "", "", "", err
	}

	if data.Error.Message != "" {
		return "", "", "", errors.New(data.Error.Message)
	}

	return username, data.Email, data.Id, nil
}

func (auth *openAuthenticator) saveState(username, state string) {
	auth.cache.Save(auth.cachePrefix+username, state, 15*time.Minute)
	auth.cache.Save(auth.cachePrefix+state, username, 15*time.Minute)
}

func (auth *openAuthenticator) getUsernameByState(username string) string {
	if state, ok := auth.cache.Get(auth.cachePrefix + username); ok {
		return state.(string)
	}
	return ""
}

func (auth *openAuthenticator) getStateByUsername(username string) string {
	if state, ok := auth.cache.Get(auth.cachePrefix + username); ok {
		return state.(string)
	}
	return ""
}

type facebookResponse struct {
	Id    string
	Email string
	Error facebookError
}

type facebookError struct {
	Message string
	Code    int
}
