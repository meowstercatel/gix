package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var data = []string{"plik1", "plik2", "plik3"} // tu lista itemkow z kompa

func main() {
	myApp := app.New()
	makeTray(myApp)

	myWindow := myApp.NewWindow("Form Widget")
	FolderPathLabel := widget.NewLabel("nie wybrales folderu")

	Commit := widget.NewEntry()
	ComitArea := widget.NewMultiLineEntry()
	RepoAea := widget.NewMultiLineEntry()
	passwordArea := widget.NewPasswordEntry()
	ChooseType := widget.NewLabel("Pick smth")
	CheckAllPull := widget.NewButton("Check all", func() {})
	CheckAllPush := widget.NewButton("Check all", func() {})
	PullButn := widget.NewButton("Pull", func() {})
	SaveButton := widget.NewButton("Save", func() { log.Println("Wysywa wartosc do serwera czy cos jakos wykminicie") })
	var FOlderPath string
	SelectFolderButton := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				if err.Error() == "Cancelled" {
					FolderPathLabel.SetText("Zostal wykonany cancell podczas wyboru folderu")
				} else {
					log.Println("Error Podczas wyboru folderu:", err)
					FolderPathLabel.SetText("Error:" + err.Error())
				}
				return
			}
			if uri == nil {
				FolderPathLabel.SetText("Brak wybranego folderu")
				return
			}
			FolderPathLabel.SetText(uri.Path())
			fmt.Println("Wybrany folder to:", uri.Path())
			FOlderPath = uri.Path()
			log.Println(FOlderPath)

		}, myWindow)
	})
	ListOfitems := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: "Name of Repo", Widget: RepoAea,
			}, {
				Text: "Password(if included)", Widget: passwordArea,
			}},
		OnSubmit: func() {
			log.Println("Repo:", RepoAea.Text)
			log.Println("Pasword:", passwordArea.Text)
			myWindow.Close()
		},
	}
	form.Hide()

	YourRepo := widget.NewSelect([]string{"I tu ma ", "Wyswitlac epo"}, func(value string) {

	})
	YourRepo.Hide()

	Choose := widget.NewSelect([]string{"Your Repo", "Public Repo"}, func(value string) {
		if value == "Public Repo" {
			YourRepo.Hide()
			log.Println("Tu sie pojawia form")
			form.Show()
		} else {
			form.Hide()
			log.Println("Tu ma sie wypisywac przypisane repo")
			YourRepo.Show()
		}
	})

	PushCon := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: "Verion", Widget: Commit,
			},
			{
				Text: "Value of Commit", Widget: ComitArea,
			}},
		OnSubmit: func() {
			log.Println("Version: ", Commit.Text)
			log.Println("Commit: ", ComitArea.Text)
		},
	}

	Pull := container.NewHSplit(ListOfitems, container.NewVBox(SelectFolderButton, FolderPathLabel, CheckAllPull, PullButn))
	Push := container.NewHSplit(ListOfitems, container.NewVBox(PushCon, CheckAllPush))
	Repositor := container.NewVBox(ChooseType, Choose, form, YourRepo, SaveButton)

	tabs := container.NewAppTabs(
		container.NewTabItem("Repository", Repositor),
		container.NewTabItem("Pull", Pull),
		container.NewTabItem("Push", Push),
		container.NewTabItemWithIcon("settings", theme.SettingsIcon(), widget.NewLabel("Kiedy tu beda ustawienia konta itp ze z jakiego brancha bierzesz albo mozna to dac jako kolejnego taba")),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	Repositor.Position()
	myWindow.CenterOnScreen()
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(600, 350))
	myWindow.ShowAndRun()
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}
