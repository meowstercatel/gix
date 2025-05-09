package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"os"
	"os/exec"
	"resty.dev/v3"
	"strconv"
	"strings"
	"time"
)

var repo Repo

var settings Settings

type Settings struct {
	Server string `json:"server"`
}

type Commit struct {
	Message     string `json:"message"`
	ZipLocation string `json:"zip_location"`
	Author      string `json:"author"`
	Date        int64  `json:"date"`
}

type Repo struct {
	Author   string   `json:"author"`
	Name     string   `json:"name"`
	Commits  []Commit `json:"commits"`
	Location string
}

func (r *Repo) Save() {
	jsonRepo, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	if os.WriteFile(r.Location+"/.gix/repo.json", jsonRepo, 0644) != nil {
		panic(err)
	}
}

func LoadRepo(location string) bool {
	r := Repo{}
	jsonRepo, err := os.ReadFile(location + "/.gix/repo.json")
	if err != nil {
		return false
	}
	if json.Unmarshal(jsonRepo, &r) != nil {
		return false
	}
	return true
}

func (r *Repo) AddCommit(commit Commit) {
	r.Commits = append(r.Commits, commit)
	r.Save()
}

func InitRepo(name string, location string) bool {
	repo = Repo{
		Name:     name,
		Location: location,
	}
	repo.Save()
	return true
}

func CreateCommit(message string) {
	files, _ := os.ReadDir(repo.Location)
	counter := 0
	for _, f := range files {
		name := strings.Split(f.Name(), ".")[0]
		nameInt, _ := strconv.Atoi(name)
		if nameInt > counter {
			counter = nameInt
		}
	}
	lastFullTar := counter - counter%5
	counter++ //this line is important because we don't want to overwrite the last tarball
	tarName := strconv.Itoa(counter) + ".tar.gz"
	tar := exec.Command("tar",
		"-cvf",
		tarName,
		"--listed-incremental=.gix/"+strconv.Itoa(lastFullTar)+".snar",
		"~",
	)
	tar.Dir = repo.Location
	if err := tar.Run(); err != nil {
		panic(err)
	}
	commit := Commit{
		Message:     message,
		ZipLocation: tarName, //incremental tarball file name
		Author:      "",
		Date:        time.Now().Unix(),
	}
	repo.AddCommit(commit)
}

func Push() {
	lastCommit := repo.Commits[len(repo.Commits)-1]

	file, err := os.ReadFile(lastCommit.ZipLocation)
	if err != nil {
		panic(err)
	}

	client := resty.New()
	defer client.Close()

	//implementing this server-side is going to be tricky
	//we should also send the incremental file (.snar) but i'll leave that to someone else
	_, err = client.R().
		SetBody(file).
		SetHeaders(map[string]string{
			"Content-Type": "application/octet-stream",
			"Author":       "",
		}).
		Post(settings.Server)
	if err != nil {
		panic(err)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("gix")

	commitMessage := widget.NewEntry()

	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(container.NewVBox(
		commitMessage,
		widget.NewButton("commit", func() {
			CreateCommit(commitMessage.Text)
		}),
		widget.NewButton("open/init repo", func() {
			folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				if !LoadRepo(uri.Path()) {
					InitRepo(uri.Name(), uri.Path())
				}
			}, w)
			folderDialog.Show()
		}),
	))

	w.ShowAndRun()
}
