package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Repository holds the data for each repository
type Repository struct {
	// LastUpdateScripts is the last date when the aplicaction checked for updates on scripts
	LastUpdateScripts string
	// Repo is actually the username on github that acts as a repo (i.e. "DarkAnHell")
	Repo string
}

// RepoDataGithub holds the information we want to retrieve from the update request to github
type RepoDataGithub struct {
	Files map[string]map[string]interface{} `json:"files"`       // Has the raw url to get the content
	ID    string                            `json:"id"`          // Unique ID of the script
	Desc  string                            `json:"description"` // Contains json with the data
}

// GetUpdates will return the new scripts to be added, and save the date if the update was successful
func (r *Repository) GetUpdates() (newScripts Scripts, err error) {

	var trailing string

	// Only update since the last update, if there was one
	if r.LastUpdateScripts != "" {
		var t string
		t, err = AddOneSecond(r.LastUpdateScripts)
		if err != nil {
			return newScripts, errors.New(translate.ErrAddSecond + ": " + err.Error())
		}

		trailing = "?since=" + t

	} else {
		trailing = ""
	}

	response, err := http.Get("https://api.github.com/users/" + r.Repo + "/gists" + trailing)

	if err != nil {
		return newScripts, errors.New(translate.ErrGithubResponse + ": " + err.Error())
	}

	if err == nil {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return newScripts, errors.New(translate.ErrBodyRead + ": " + err.Error())
		}

		// Parse JSON
		var bunchOfData []RepoDataGithub

		if err := json.Unmarshal(contents, &bunchOfData); err != nil {
			return newScripts, errors.New(translate.ErrGithubJSONParse + ": " + err.Error())
		}

		// GET ALL THE SCRIPTS!!!
		for _, s := range bunchOfData {
			script, err := generateScriptFromData(s)

			// We ignore JSON parsing errors, since they will be - most likely - gists other than scripts
			// Nontheless, we log them, just in case
			if err != nil && strings.Contains(err.Error(), translate.NoticeScriptJSONParse) {
				l.Println(err)
			} else if err != nil {
				return newScripts, errors.New(translate.ErrScriptGeneration + ": " + err.Error())
			}

			newScripts = append(newScripts, script)
		}
	}

	// Save time of last update
	r.LastUpdateScripts = time.Now().UTC().Format("2006-01-02T15:04:05-0700")

	return
}

func generateScriptFromData(data RepoDataGithub) (s Script, err error) {
	s.ID = data.ID

	// Assign data from description
	if err := json.Unmarshal([]byte(data.Desc), &s); err != nil {
		return s, errors.New(translate.NoticeScriptJSONParse + ": " + err.Error())
	}

	// Get remote URL from file info
	for _, v := range data.Files {
		for k, v2 := range v {

			if k != "raw_url" {
				continue
			}

			s.RemotePath = v2.(string)
			break
		}
	}

	s.Remote = true

	return

}

// AddOneSecond will return the same update date but with one second more.
// This makes github not answer with already updated scripts when requesting them
func AddOneSecond(date string) (tS string, err error) {
	t, err := time.Parse(githubRFC, date)
	if err != nil {
		return
	}

	tS = t.Add(time.Second).String()

	return
}
