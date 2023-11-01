package pages

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/client"
	"abeProofOfConcept/pkg/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

func NewMailBoxScreen(mainWindow fyne.Window, email string, messages []models.Message) fyne.CanvasObject {

	titleLabel := widget.NewLabelWithStyle("MailBox", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	newMessageButton := widget.NewButton("New Message", nil)

	profileButton := widget.NewButton(
		"Profile", func() {
			mainWindow.SetContent(NewProfileWindow(mainWindow, email))
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

	var scrollableMessagesList *container.Scroll

	reloadButton := widget.NewButtonWithIcon(
		"Refresh", theme.ViewRefreshIcon(), func() {
			msgs, err := client.LoadMailBox()
			if err == nil {
				messages = msgs
				scrollableMessagesList.Refresh()
			} else {
				logger.Logger.Errorln("message", "Failed reloading the mailbox", "error", err)
				ShowErrorAlert("Couldn't reload mailbox! Please try again later.", mainWindow)
			}
		},
	)

	buttonsContainer := container.NewHBox(
		newMessageButton,
		profileButton,
		layout.NewSpacer(), // Make spacing
		reloadButton,
		logOutButton,
	)

	//Messages list
	messagesList := widget.NewList(
		func() int {
			return len(messages)
		},
		func() fyne.CanvasObject {

			return container.NewVBox(
				widget.NewLabelWithStyle("Title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Sender"),
				widget.NewLabel("Date"),
			)
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			// Ažurira element za prikaz poruke
			vbox := obj.(*fyne.Container)
			msgTitleLabel := vbox.Objects[0].(*widget.Label)
			senderLabel := vbox.Objects[1].(*widget.Label)
			dateLabel := vbox.Objects[2].(*widget.Label)

			msgTitleLabel.SetText(messages[i].Title)
			senderLabel.SetText(messages[i].From)
			dateLabel.SetText(messages[i].CreatedAt.Format("02/01/2006\t15:04"))
		},
	)

	scrollableMessagesList = container.NewVScroll(messagesList)
	// Kontejner za sadržaj desne strane
	rightContainer := container.NewVBox()

	// Kontent poruke (kada se klikne na poruku)
	messageTitleLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	messageFromLabel := widget.NewLabel("From: ")
	messageToLabel := widget.NewLabel("To: ")
	messageDateTimeLabel := widget.NewLabel("Date and time: ")
	messageContentLabel := widget.NewLabel("Select a message to view its content.")

	messagesList.OnSelected = func(id widget.ListItemID) {
		messageTitleLabel.SetText("Title: " + messages[id].Title)
		messageFromLabel.SetText("From: " + messages[id].From)
		messageToLabel.SetText("To: " + strings.Join(messages[id].To, ", "))
		messageDateTimeLabel.SetText("Date and time: " + messages[id].CreatedAt.Format("02/01/2006\t15:04"))
		messageContentLabel.SetText(messages[id].Content)

		// Postavljanje svih komponenti za prikaz poruke unutar desne strane
		rightContainer.Objects = []fyne.CanvasObject{
			messageTitleLabel, messageFromLabel, messageToLabel, messageDateTimeLabel, messageContentLabel,
		}
		rightContainer.Refresh()
	}

	////////////////////////////////
	//
	//Send message sub window
	var filePath string

	toLabel := widget.NewLabel("To:")
	toEntry := widget.NewEntry()
	titleNewMessageLabel := widget.NewLabel("Title:")
	titleNewMessageEntry := widget.NewEntry()
	separator := widget.NewSeparator()

	var entries []*widget.Entry
	var recipientsPolicies []*widget.Select
	var positionPolicies []*widget.Select
	var departmentPolicies []*widget.Select
	var radioGroups []*widget.RadioGroup
	createEntry := func() *widget.Entry {
		entry := widget.NewMultiLineEntry()
		entries = append(entries, entry)
		return entry
	}

	toEntry.OnChanged = func(s string) {
		opt := strings.Split(strings.TrimSpace(toEntry.Text), ";")
		for _, pol := range recipientsPolicies {
			pol.Options = opt
			pol.Refresh()
		}
	}

	sendMessageButton := widget.NewButton(
		"Send Message", func() {
			msgTo := strings.Split(toEntry.Text, ";")
			for i := 0; i < len(msgTo); i++ {
				msgTo[i] = strings.TrimSpace(msgTo[i])
			}
			msgTitle := strings.TrimSpace(titleNewMessageEntry.Text)
			if len(msgTo) == 0 || msgTitle == "" {
				//TODO alert the user
				return
			}
			if filePath == "" {
				var contents []string
				var policies []string
				for i, entry := range entries {
					content := strings.TrimSpace(entry.Text)
					if content == "" {
						continue
					}
					contents = append(contents, content)
					var policy string
					if radioGroups[i].Selected == "Specific Recipient" && recipientsPolicies[i].Selected != "" {
						policy = "email:" + recipientsPolicies[i].Selected

					} else {
						var deptPolicy string
						var posPolicy string
						if departmentPolicies[i].Selected != "" {
							deptPolicy = "department:" + departmentPolicies[i].Selected
						}
						if positionPolicies[i].Selected != "" {
							posPolicy = "position:" + positionPolicies[i].Selected
						}
						if deptPolicy != "" && posPolicy != "" {
							policy = deptPolicy + " AND " + posPolicy
						} else {
							policy = deptPolicy + posPolicy
						}
					}
					policies = append(policies, policy)
				}
				err := client.SendMessage(email, msgTo, msgTitle, contents, policies)
				if err != nil {
					ShowErrorAlert("Server error sending message! Please try again later.", mainWindow)
				}
			} else {
				err := client.SendMessageFromFile(email, msgTo, msgTitle, filePath)
				if err != nil {
					ShowErrorAlert("Server error sending message! Please try again later.", mainWindow)
				}
			}
			msgs, err := client.LoadMailBox()
			if err == nil {
				messages = msgs
				scrollableMessagesList.Refresh()
			} else {
				logger.Logger.Errorln("message", "Failed reloading the mailbox", "error", err)
				ShowErrorAlert("Couldn't reload mailbox! Please try again later.", mainWindow)
			}

		},
	)

	sendMessageButtonColor := canvas.NewRectangle(color.NRGBA{R: 255, B: 0, G: 0, A: 255})

	messageContentList := container.NewVBox() // Ovo će čuvati sve MultiLineEntry-e
	scrollableContentList := container.NewVScroll(messageContentList)
	scrollableContentList.SetMinSize(fyne.NewSize(scrollableContentList.Size().Width, 400))
	container1 := container.New(
		layout.NewMaxLayout(),
		sendMessageButtonColor,
		sendMessageButton,
	)

	resetMessageButton := widget.NewButton(
		"Clear Message", func() {
			for _, entry := range entries {
				entry.SetText("")
			}
			messageContentList.Objects = nil
			messageContentList.Refresh()
		},
	)

	newParagraphButton := widget.NewButton(
		"New Paragraph", func() {

			employeeDropdownLocal := widget.NewSelect(
				strings.Split(strings.TrimSpace(toEntry.Text), ";"), nil,
			)
			positionDropdownLocal := widget.NewSelect([]string{"junior", "senior", "manager"}, nil)
			departmentDropdownLocal := widget.NewSelect([]string{"engineering", "marketing", "finance"}, nil)

			recipientTypeLocal := widget.NewRadioGroup(
				[]string{"Specific Recipient", "Attributes"}, func(selected string) {
					switch selected {
					case "Specific Recipient":
						employeeDropdownLocal.Enable()
						positionDropdownLocal.Disable()
						departmentDropdownLocal.Disable()
					case "Attributes":
						employeeDropdownLocal.Disable()
						positionDropdownLocal.Enable()
						departmentDropdownLocal.Enable()
					}
				},
			)
			recipientTypeLocal.SetSelected("Specific Recipient")
			dropdownContainerLocal := container.NewHBox(
				employeeDropdownLocal,
				positionDropdownLocal,
				departmentDropdownLocal,
			)

			// Add all instances to corresponding slices/lists
			entry := createEntry()
			radioGroups = append(radioGroups, recipientTypeLocal)
			recipientsPolicies = append(recipientsPolicies, employeeDropdownLocal)
			positionPolicies = append(positionPolicies, positionDropdownLocal)
			departmentPolicies = append(departmentPolicies, departmentDropdownLocal)

			messageContentList.Add(recipientTypeLocal)
			messageContentList.Add(dropdownContainerLocal)
			messageContentList.Add(entry)
			messageContentList.Refresh()
		},
	)

	//Load file button

	loadFileButton := widget.NewButton(
		"Pick a File", func() {
			fileDialog := dialog.NewFileOpen(
				func(reader fyne.URIReadCloser, err error) {
					if err == nil && reader != nil {
						// Get the selected file's path
						filePath = reader.URI().Path()
						logger.Logger.Println("File selected " + filePath)
						reader.Close()
					}
				}, mainWindow,
			)

			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"})) // Add desired file extensions
			fileDialog.Show()
		},
	)

	buttonContainer := container.NewHBox(newParagraphButton, loadFileButton, resetMessageButton, container1)

	newMessageButton.OnTapped = func() {
		rightContainer.Objects = []fyne.CanvasObject{
			toLabel, toEntry,
			titleNewMessageLabel, titleNewMessageEntry,
			buttonContainer,
			separator,
			scrollableContentList,
		}
		rightContainer.Refresh()
	}

	scrollableMessagesList.SetMinSize(fyne.NewSize(scrollableMessagesList.Size().Width, 600))

	splitContainer := container.NewHSplit(
		scrollableMessagesList, // left pane
		rightContainer,         // right pane
	)

	splitContainer.Offset = 0.25 // left pane

	// Group all components vertically
	content := container.NewVBox(
		titleLabel,
		buttonsContainer,
		splitContainer,
	)

	return content
}
