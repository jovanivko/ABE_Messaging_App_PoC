package pages

import (
	"abeProofOfConcept/pkg/client"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// NewLogInWindow vraća novi prozor za login
func NewLogInWindow(mainWindow fyne.Window) fyne.CanvasObject {

	titleLabel := widget.NewLabelWithStyle("Log In", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	usernameLabel := widget.NewLabel("Email:")
	usernameEntry := widget.NewEntry()

	passwordLabel := widget.NewLabel("Password:")
	passwordEntry := widget.NewPasswordEntry()

	errorLabel := canvas.NewText(" ", color.RGBA{R: 233, A: 255})
	errorLabel.Hide()

	var content *fyne.Container
	showErrorMessage := func(message string) {

		if message != "" {
			errorLabel.Text = message
			errorLabel.Show()
			errorLabel.TextStyle = fyne.TextStyle{Italic: true}
		} else {
			errorLabel.Hide()
		}
		content.Refresh()
	}

	loginButton := widget.NewButton(
		"Log In", func() {
			showErrorMessage("")
			var loginError string
			if usernameEntry.Text == "" && passwordEntry.Text == "" {
				loginError = "Email and password are required for login!"
			} else if usernameEntry.Text == "" {
				loginError = "Email is required for login!"
			} else if passwordEntry.Text == "" {
				loginError = "Password is required for login!"
			}
			if loginError != "" {
				showErrorMessage(loginError)
				return
			}
			keys, err := client.LogIn(usernameEntry.Text, passwordEntry.Text)
			if err == nil {
				client.FamePublicKey = keys.FamePublicKey
				showErrorMessage("")
				mailbox, err := client.LoadMailBox()
				if err != nil {
					ShowErrorAlert("Server error sending message! Please try again later.", mainWindow)
				}
				mainWindow.SetContent(NewMailBoxScreen(mainWindow, usernameEntry.Text, mailbox))

			} else {
				if err.Error() == "Recovered from panic in ReadResponse" {
					ShowErrorAlert(
						"The connection with the server could not be established. Please try again later.", mainWindow,
					)
				}
				showErrorMessage("Wrong email or password!")
			}
		},
	)

	backButton := widget.NewButton(
		"Back", func() {
			ShowStartScreen(mainWindow)
		},
	)

	bottomBar := container.NewHBox(layout.NewSpacer(), backButton)

	content = container.NewBorder(
		nil, bottomBar, nil, nil,
		container.NewVBox(
			titleLabel,
			usernameLabel,
			usernameEntry,
			passwordLabel,
			passwordEntry,
			loginButton,
			errorLabel,
		),
	)

	return content
}

func ShowStartScreen(mainWindow fyne.Window) {

	title := widget.NewLabelWithStyle(
		"ABE secured company messaging application", fyne.TextAlignCenter, fyne.TextStyle{Bold: true},
	)
	subtitle := widget.NewLabelWithStyle(
		"~Attribute-based encryption~", fyne.TextAlignCenter, fyne.TextStyle{Italic: true},
	)

	startButton := widget.NewButton(
		"Start", func() {
			mainWindow.SetContent(NewLogInWindow(mainWindow))

		},
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		title,
		subtitle,
		layout.NewSpacer(),
		startButton,
		layout.NewSpacer(),
	)

	mainWindow.SetContent(content)
}

func ShowErrorAlert(errorMessage string, parent fyne.Window) {
	dialog.ShowError(errors.New(errorMessage), parent)
}
