package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// Script holds all the data from a script
type Script struct {
	// Path holds the path to the script
	Path string
	// Name holds the name of the script
	Name string
	// Desc holds the description for the script
	Desc string
	// Tags holds the tags for the script, such as OS, Desktop Enviroment, etc.
	Tags string
	// User saves the creator's name/nick
	User string
	//Hash holds a md5 of the last known state of the script
	Hash string
}

// NewScript returns a new object of type Script
func NewScript() *Script { return &Script{} }

// SearchLocal will search on the path for valid scripts and load them onto the already loaded scripts
func SearchLocal(path string, scripts *[]Script) (count uint, err error) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		isNew := true
		for _, script := range *scripts {
			if script.Path == path+"/"+f.Name() {
				isNew = false
				break
			}
		}

		if isNew {

			h, err := HashFile(path + "/" + f.Name())
			if err != nil {
				return count, err
			}

			count++
			*scripts = append(*scripts, Script{
				Path: path + "/" + f.Name(),
				Hash: h})
		}
	}

	return
}

// CheckLocal will check the local database of scripts.
// Will create a new one if it doesn't exists.
// Returns all the loaded data, plus info about the changes made
func CheckLocal(path, scriptsPath string) (
	scripts []Script,
	totalValid, deleted, modified, newOnes uint,
	err error,
) {

	// Read / create database
	f, err := ioutil.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		// database does not exist, create
		f2, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return scripts, totalValid, deleted, modified, newOnes, err
		}
		f2.WriteString("[{}]")
		f2.Close()

		f, err = ioutil.ReadFile(path)
		if err != nil {
			return scripts, totalValid, deleted, modified, newOnes, err
		}

	}

	if err = json.Unmarshal(f, &scripts); err != nil {
		return
	}

	// Check database integrity
	for i, script := range scripts {
		fE, hE, hash := script.CheckIntegrity()
		if fE {
			deleted++
			scripts = append(scripts[:i], scripts[i+1:]...)
		} else if !hE {
			scripts[i].Hash = hash
			modified++
			totalValid++
		} else {
			totalValid++
		}
	}

	// Look for new files
	newOnes, err = SearchLocal(scriptsPath, &scripts)
	if err != nil {
		return
	}

	totalValid += newOnes

	b, err := json.Marshal(scripts)
	if err != nil {
		return
	}

	f2, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}

	_, err = f2.Write(b)
	if err != nil {
		return
	}

	return
}

// HashFile will hash a file and return the hexsum
func HashFile(path string) (h string, err error) {
	var result []byte
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	h = hex.EncodeToString(hash.Sum(result))

	return
}

// GetName will return the name for the script, or it's path if it does not
// have a name assigned
func (s *Script) GetName() string {
	if s.Name != "" {
		return s.Name
	}
	return s.Path
}

// CheckIntegrity will load the file from the path and check it against known data.
func (s *Script) CheckIntegrity() (fileErr, hashEq bool, h string) {
	h, err := HashFile(s.Path)
	if err == nil {
		fileErr = false
		if h != s.Hash {
			hashEq = false
		} else {
			hashEq = true
		}
	} else {
		fileErr = true
		hashEq = false
	}

	return
}

// Equals will check if the scripts are the same object
func (s *Script) Equals(s2 Script) bool { return s.Hash == s2.Hash }

// Save will save the data to the json database
func Save(path string, scripts []Script) {
	b, err := json.Marshal(scripts)
	if err != nil {
		l.Println(err.Error())
		return
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		l.Println(err.Error())
		return
	}
	defer f.Close()
	f.Write(b)
}

/*** SORT FUNCTIONS ***/

// Scripts is a dummy type for qsort
type Scripts []Script

func (s Scripts) Len() int {
	return len(s)
}
func (s Scripts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Scripts) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
}

/*** END SORT FUNCTIONS***/

// SortScripts sorts the slice based on name and path
func SortScripts(scripts Scripts) []Script {
	sort.Sort(scripts)

	return scripts
}

// ListByName will return all scripts which contains the name on it
// An empty string is interpreted as "any".
// This function IS NOT case sensitive.
func ListByName(scripts []Script, name string) (matches []Script) {
	for _, script := range scripts {
		if strings.Contains(strings.ToLower(script.Name), strings.ToLower(name)) {
			matches = append(matches, script)
		}
	}
	return
}

// ListByUser will return all scripts that were created by a user containing the string.
// An empty string is interpreted as "any".
// This function IS NOT case sensitive.
func ListByUser(scripts []Script, user string) (matches []Script) {
	for _, script := range scripts {
		if strings.Contains(strings.ToLower(script.User), strings.ToLower(user)) {
			matches = append(matches, script)
		}
	}
	return
}

// ListByTags will return all scripts that contain a tag which contains the passed keyword
// An empty string is interpreted as "any".
// This function IS NOT case sensitive.
func ListByTags(scripts []Script, inTag string) (matches []Script) {
	for _, script := range scripts {
		if strings.Contains(strings.ToLower(script.Tags), strings.ToLower(inTag)) {
			matches = append(matches, script)
		}

	}
	return
}

// ListByDesc will return all scripts that has a description containing a string
// An empty string is interpreted as "any".
// This function IS NOT case sensitive.
func ListByDesc(scripts []Script, desc string) (matches []Script) {
	for _, script := range scripts {
		if strings.Contains(strings.ToLower(script.Desc), strings.ToLower(desc)) {
			matches = append(matches, script)
		}
	}
	return
}

// TrimRepeated will return an alice formed from the passed one,
// but the result will only include unique values.
func TrimRepeated(scripts []Script) (valid []Script) {

	for _, s := range scripts {
		l.Println(scripts, valid)
		eq := false
		for _, v := range valid {
			if s.Equals(v) {
				eq = true
				break
			}
		}
		if !eq {
			valid = append(valid, s)
		}
	}

	return

}
