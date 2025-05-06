package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
)

var repo Repo

type Commit struct {
	Message     string `json:"message"`
	ZipLocation string `json:"zip_location"`
	Author      string `json:"author"`
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
	if os.WriteFile(r.Location+"/repo.json", jsonRepo, 0644) != nil {
		panic(err)
	}
}

func LoadRepo(location string) Repo {
	r := Repo{}
	jsonRepo, err := os.ReadFile(location + "/repo.json")
	if err != nil {
		panic(err)
	}
	if json.Unmarshal(jsonRepo, &r) != nil {
		panic(err)
	}
	return r
}

func (r *Repo) AddCommit(commit Commit) {
	r.Commits = append(r.Commits, commit)
	r.Save()
}

func InitRepo(name string, location string) Repo {
	repo = Repo{
		Name:     name,
		Location: location,
	}
	repo.Save()
	return repo
}

func main() {
	a := app.New()
	w := a.NewWindow("gix")

	commitMessage := widget.NewEntry()

	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(container.NewVBox(
		commitMessage,
		widget.NewButton("commit", func() {
			commit := Commit{
				Message:     commitMessage.Text,
				ZipLocation: "", //zip diff location
				Author:      "",
			}
			repo.AddCommit(commit)
			repo.Save()
		}),
		widget.NewButton("init repo", func() {
			folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
				fmt.Println(uri)
				InitRepo(uri.Name(), uri.Path())
			}, w)
			folderDialog.Show()
		}),
	))

	w.ShowAndRun()
}
