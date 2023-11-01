package main

import (
	"abeProofOfConcept/pkg/client"
	"abeProofOfConcept/pkg/pages"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var mainWindow fyne.Window

func main() {
	client.NewHttpClient()
	a := app.New()
	mainWindow = a.NewWindow("Company messaging application")

	// Podešavanje dimenzija prozora na srednje veličine
	mainWindow.Resize(fyne.NewSize(700, 600))
	pages.ShowStartScreen(mainWindow)

	// Postavljanje teme da bi izgledalo lepo
	a.Settings().SetTheme(theme.DarkTheme())

	// Pokretanje aplikacije
	mainWindow.ShowAndRun()
}
