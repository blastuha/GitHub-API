package api

import "net/http"

var GithubClient = &http.Client{
	Transport: &http.Transport{DisableKeepAlives: false},
}
