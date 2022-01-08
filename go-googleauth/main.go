package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/callback",
		ClientID:     "782280680980-mg73o4fhqllch96s65mbqkcp47mvckkk.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/api/v1/callback", handleCallback)
	http.ListenAndServe(":8080", nil)

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">Google Sign</a></body></html>`
	fmt.Fprintf(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfGl.AuthCodeURL(oauthStateStringGl)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateStringGl {
		fmt.Println("error, state not valid")
		http.Redirect(w, r, "/api/v1/google/login", http.StatusTemporaryRedirect)
		return
	}
	token, err := oauthConfGl.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Sprintf("could not get token %s\n", err.Error())
		http.Redirect(w, r, "/api/v1/google/login", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Sprintf("could not get request %s\n", err.Error())
		http.Redirect(w, r, "/api/v1/google/login", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintf("could not parse response %s\n", err.Error())
		http.Redirect(w, r, "/api/v1/google/login", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "response : %s", content)
}
