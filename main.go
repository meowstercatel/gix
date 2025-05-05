package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type repo struct {
	Name    string   `json:"name"`
	Commits []commit `json:"commits"`
}

type commit struct {
	Message     string `json:"message"`
	ZipLocation string `json:"zip_location"`
	Author      string `json:"author"`
}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("init repo", func() {
			hello.SetText("oke")
		}),
	))

	w.ShowAndRun()
}
