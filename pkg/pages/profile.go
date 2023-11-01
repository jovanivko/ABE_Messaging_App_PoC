package pages

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/client"
	"abeProofOfConcept/pkg/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func NewProfileWindow(mainWindow fyne.Window, email string) fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle("Profile", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	mailboxButton := widget.NewButton(
		"MailBox", func() {
			messages, err := client.LoadMailBox()
			if err != nil {
				ShowErrorAlert("Couldn't load mailbox! Please try again later.", mainWindow)
			}
			mainWindow.SetContent(NewMailBoxScreen(mainWindow, email, messages))
		},
	)

	logOutButton := widget.NewButton(
		"Log Out", func() {
			err := client.LogOut()
			if err != nil {
				logger.Logger.Errorln("Couldn't log out")
				return
			}
			mainWindow.SetContent(NewLogInWindow(mainWindow))
		},
	)
	buttonsContainer := container.NewHBox(
		titleLabel,
		layout.NewSpacer(),
		mailboxButton,
		logOutButton,
	)
	var user *models.User

	var formBox *fyne.Container

	nameLabel := widget.NewLabel("Firstname: ")
	nameValue := widget.NewLabel("")
	nameBox := container.NewHBox(nameLabel, nameValue)

	surnameLabel := widget.NewLabel("Lastname: ")
	surnameValue := widget.NewLabel("")
	surnameBox := container.NewHBox(surnameLabel, surnameValue)

	positionLabel := widget.NewLabel("Position: ")
	positionValue := widget.NewLabel("")
	positionBox := container.NewHBox(positionLabel, positionValue)

	departmentLabel := widget.NewLabel("Department: ")
	departmentValue := widget.NewLabel("")
	departmentBox := container.NewHBox(departmentLabel, departmentValue)

	phoneNumberLabel := widget.NewLabel("Phone number: ")
	phoneNumberValue := widget.NewLabel("")
	phoneBox := container.NewHBox(phoneNumberLabel, phoneNumberValue)

	addressLabel := widget.NewLabel("Address: ")
	addressValue := widget.NewLabel("")
	addressBox := container.NewHBox(addressLabel, addressValue)

	salaryLabel := widget.NewLabel("Salary: ")
	salaryValue := widget.NewLabel("")
	salaryBox := container.NewHBox(salaryLabel, salaryValue)

	profileBox := container.NewVBox(nameBox, surnameBox, positionBox, departmentBox, phoneBox, addressBox, salaryBox)
	profileBox.Hide()

	emails, err := client.Emails()
	if err != nil {
		emails = []string{"", email}
	} else {
		emails = append([]string{""}, emails...)
	}
	pickerLabel := widget.NewLabel("Choose user by email:")
	userPicker := widget.NewSelect(
		emails, func(userEmail string) {
			if userEmail == "" {
				user = nil
			} else {
				user, err = client.LoadProfile(userEmail)
				if err != nil {
					ShowErrorAlert("Couldn't load profile! Please try again later.", mainWindow)
					user = nil
				}
			}
			if user == nil {
				profileBox.Hide()
			} else {
				nameValue.SetText(user.FirstName)
				surnameValue.SetText(user.LastName)
				positionValue.SetText(user.Position)
				departmentValue.SetText(user.Department)
				phoneNumberValue.SetText(user.PhoneNumber)
				if user.Salary != 0 {
					salaryValue.SetText(strconv.Itoa(user.Salary))
				}
				if user.Address != "" {
					addressValue.SetText(user.Address)
				}
				profileBox.Show()
				if user.Salary == 0 {
					salaryBox.Hide()
				} else if salaryBox.Hidden {
					salaryBox.Show()
				}
				if user.Address == "" {
					addressBox.Hide()
				} else if addressBox.Hidden {
					addressBox.Show()
				}
				profileBox.Refresh()
			}
		},
	)
	pickerBox := container.NewHBox(pickerLabel, userPicker)
	formBox = container.NewVBox(pickerBox, profileBox)

	content := container.NewBorder(buttonsContainer, nil, nil, nil, formBox)
	return content
}
