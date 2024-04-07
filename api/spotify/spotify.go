package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	apiLog      = log.New(os.Stdout, "API: ", log.Ldate|log.Ltime|log.Lshortfile)
	accessToken *tokenResponse
)

type tokenResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    float64 `json:"expires_in"`
	Scope        string  `json:"scope"`
}

type device struct {
	Id           string `json:"id"`
	IsActive     bool   `json:"is_active"`
	IsRestricted bool   `json:"is_restricted"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}
type image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"Width"`
}

type album struct {
	Images []image
}

type item struct {
	Album album  `json:"album"`
	Name  string `json:"name"`
	Id    string `json:"id"`
}

type PlayerData struct {
	Device device `json:"device"`
	Item   item   `json:"item"`
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		screamAndDie(err, apiLog)
	}
}

func Init() {
	loadDotEnv()
	var (
		clientId     = os.Getenv("SPOTIFY_CLIENT_ID")
		redirectUri  = os.Getenv("SPOTIFY_REDIRECT_URI")
		authorizeUrl = os.Getenv("SPOTIFY_AUTHORIZE_URL")
	)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		apiLog.Println("Got  request in /login")
		state := "aaaaaaaaaaaaaaaa"
		http.SetCookie(w, &http.Cookie{
			Name:  "stateKey",
			Value: state,
		})
		scope := "playlist-read-private playlist-read-collaborative user-read-currently-playing user-read-playback-state user-modify-playback-state"

		spotifyRedirectUri := authorizeUrl + url.Values{
			"response_type": {"code"},
			"client_id":     {clientId}, // Replace clientID with your Spotify client ID
			"scope":         {scope},
			"redirect_uri":  {redirectUri}, // Replace redirectURI with your redirect URI
			"state":         {state},
		}.Encode()

		http.Redirect(w, r, spotifyRedirectUri, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		apiLog.Println("Got  request in /callback")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		storedState, err := r.Cookie("stateKey")
		if err != nil {
			storedState = nil
		}
		// Check if state is missing or mismatched
		if state == "" || state != storedState.Value {
			http.Redirect(w, r, "/#"+url.Values{"error": {"state_mismatch"}}.Encode(), http.StatusFound)
			return
		}
		// Clear the state cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "stateKey",
			Value:  "",
			MaxAge: -1,
		})

		tokenResponse, err := getToken(code)
		accessToken = tokenResponse
		if err != nil {
			screamAndDie(err, apiLog)
		}
	})
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
}

func TriggerAuth() {
	apiLog.Printf("Go to http://raspberrypi:8080/login to login to spotify")
	// rsp, err := http.Get("http://raspberrypi:8080/login")
	// screamAndDie(err, apiLog)
	// defer rsp.Body.Close()

	// body, err := io.ReadAll(rsp.Body)
	// screamAndDie(err, apiLog)
	// fmt.Print(string(body))
	// fmt.Print(accessToken)
}

func DoStuff() {
	fmt.Println("accessToken", accessToken)

	type device struct {
		Id           string `json:"id"`
		IsActive     bool   `json:"is_active"`
		IsRestricted bool   `json:"is_restricted"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	}
	var playerData struct {
		Device device `json:"device"`
	}
	Get("/me/player", &playerData)
	fmt.Println(&playerData.Device)
	Get("/users/bakr.chakk/playlists", &playerData)
}

func Get(endpoint string, jsonStruct interface{}) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1"+endpoint, nil)
	screamAndDie(err, apiLog)
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)
	// Make the API request
	client := &http.Client{}
	resp, err := client.Do(req)
	screamAndDie(err, apiLog)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	screamAndDie(err, apiLog)
	err = json.Unmarshal(body, &jsonStruct)
	screamAndDie(err, apiLog)
	// fmt.Println(endpoint, ": ", jsonStruct.Device)

}

func GetPlaybackState() *PlayerData {
	var pd PlayerData
	Get("/me/player", &pd)
	return &pd

}

func GetTrackName() string {
	var pd PlayerData
	Get("/me/player", &pd)
	return pd.Item.Name

}

func getToken(code string) (token *tokenResponse, err error) {

	var (
		clientId     = os.Getenv("SPOTIFY_CLIENT_ID")
		clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
		authUrl      = os.Getenv("SPOTIFY_AUTH_URL")
		redirectUri  = os.Getenv("SPOTIFY_REDIRECT_URI")
	)
	fmt.Println("clientId: ", clientId)
	fmt.Println("clientSecret:", clientSecret)
	fmt.Println("authUrl:", authUrl)
	authFormData := url.Values{}
	authFormData.Set("grant_type", "authorization_code")
	authFormData.Set("redirect_uri", redirectUri)
	authFormData.Set("code", code)

	req, err := http.NewRequest("POST", authUrl, strings.NewReader(authFormData.Encode()))
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("%s:%s", clientId, clientSecret)
	authHeaderValue := base64.StdEncoding.EncodeToString([]byte(authHeader))
	req.Header.Set("Authorization", "Basic "+authHeaderValue)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = authFormData.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var body tokenResponse
		err := json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
		return &body, nil
	} else {
		return nil, err
	}

}

func screamAndDie(err error, logger *log.Logger) {
	if err != nil {
		logger.Panic(err)
	}
}
