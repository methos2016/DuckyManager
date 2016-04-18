package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Repository holds the data for each repository
type Repository struct {
	// LastUpdateScripts is the last date when the aplicaction checked for updates on scripts
	LastUpdateScripts string
	// Repo is actually the username on github that acts as a repo (i.e. "DarkAnHell")
	Repo string
}

// GetUpdates will return the new scripts to be added
func (r *Repository) GetUpdates() (newScripts Scripts, err error) {
	t, err := AddOneSecond(r.LastUpdateScripts)
	if err != nil {
		return //TODO return translated error
	}

	response, err := http.Get("http://api.github.com/users/" + r.Repo + "/gists?since=" + t)

	if err == nil {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return // TODO return translated error
		}

		// TODO get new scripts from json recieved
	}

	return
}

// AddOneSecond will return the same update date but with one second more.
// This makes github not answer with already updated scripts when requesting them
func AddOneSecond(date string) (tS string, err error) {
	t, err := time.Parse(time.RFC3339Nano, date) // Close RFC as the defined by github's API
	if err != nil {
		return
	}

	tS = t.Add(time.Second).String()

	return
}
