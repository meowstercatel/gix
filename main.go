package main

import (
	"encoding/json"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

type Repo struct {
	Author string `json:"author"`

	Name    string   `json:"name"`
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Message     string `json:"message"`
	ZipLocation string `json:"zip_location"`
	Author      string `json:"author"`
}

func initRepo(name string, location string) Repo {
	repo := Repo{
		Name: name,
	}
	jsonRepo, err := json.Marshal(repo)
	if os.WriteFile(location+"/repo.json", jsonRepo, 0644) != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return Repo{}
}

func main() {
	a := app.New()
	w := a.NewWindow("gix")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("init repo", func() {
			hello.SetText("oke")
		}),
	))

	w.ShowAndRun()
}
