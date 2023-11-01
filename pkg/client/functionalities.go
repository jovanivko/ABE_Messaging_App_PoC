package client

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/models"
	"abeProofOfConcept/pkg/routes"
	"abeProofOfConcept/pkg/utils"
	"fmt"
	"github.com/fentec-project/gofe/abe"
	"strings"
)

func LoadMailBox() ([]models.Message, error) {
	var messages []models.Message
	res, err := GetMailbox()
	if err != nil {
		logger.Logger.Errorln("Failed sending the get mailbox request", "error", err)
		return nil, err
	}
	err = ReadResponse(res, &messages)
	if err != nil {
		logger.Logger.Errorln("Failed reading the messages from the mailbox response body", "error", err)
		return nil, err
	}
	return messages, nil
}

type FragmentSpec struct {
	Content string
	Policy  string
}

func SendMessage(from string, to []string, title string, contents []string, policies []string) error {
	emailPolicy := utils.CreatePolicyFromEmails(from, to)
	titleCipher, err := EncryptFame(title, emailPolicy)
	if err != nil {
		logger.Logger.Errorln("Failed encrypting the message title", "error", err)
		return err
	}
	res, err := PostMessage(&routes.MessageReq{To: to, From: from, Title: titleCipher})
	if err != nil {
		logger.Logger.Errorln("Failed posting the message", "error", err)
		return err
	}

	for i, pol := range policies {
		if pol == "" {
			policies[i] = emailPolicy
		} else {
			policies[i] = "email:" + from + " OR (" + pol + ")"
		}
		policies[i] = strings.Replace(policies[i], "position", "role", -1)
	}

	var msgId int
	err = ReadResponse(res, &msgId)
	if err != nil {
		logger.Logger.Errorln("Failed reading the messageID", "error", err)
		return err
	}
	for i, con := range contents {
		fragmentCipher, err := EncryptFame(con, policies[i])
		if err != nil {
			logger.Logger.Errorln(
				fmt.Sprintf("Failed encrypting the fragment content for fragment: %v", i), "error", err,
			)
			return err
		}
		res, err = PostFragment(
			&routes.FragmentReq{
				MsgID:      msgId,
				Content:    fragmentCipher,
				FragmentID: i + 1,
			},
		)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed sending the post message fragment request: %v", i), "error", err)
			return err
		}
		err = ReadResponse(res, nil)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed posting the message fragment: %v", i), "error", err)
			return err
		}
	}

	return nil
}

func SendMessageFromFile(from string, to []string, title string, filePath string) error {
	emailPolicy := utils.CreatePolicyFromEmails(from, to)
	titleCipher, err := EncryptFame(title, emailPolicy)
	if err != nil {
		logger.Logger.Errorln("Failed encrypting the message title", "error", err)
		return err
	}
	res, err := PostMessage(&routes.MessageReq{To: to, From: from, Title: titleCipher})
	if err != nil {
		logger.Logger.Errorln("Failed posting the message", "error", err)
		return err
	}
	var msgId int
	err = ReadResponse(res, &msgId)
	if err != nil {
		logger.Logger.Errorln("Failed reading the messageID", "error", err)
		return err
	}
	policies, contents, err := utils.ParseMessageFile(filePath)
	for i, v := range policies {
		if v == "" {
			policies[i] = emailPolicy
		} else {
			policies[i] = "email:" + from + " OR (" + v + ")"
		}
		policies[i] = strings.Replace(policies[i], "position", "role", -1)
	}
	if err != nil {
		logger.Logger.Errorln("Failed parsing the message file", "error", err)
		return err
	}
	for i, policy := range policies {
		fragmentCipher, err := EncryptFame(contents[i], policy)
		if err != nil {
			logger.Logger.Errorln(
				fmt.Sprintf("Failed encrypting the fragment content for fragment: %v", i), "error", err,
			)
			return err
		}
		fragment := &routes.FragmentReq{
			MsgID:      msgId,
			Content:    fragmentCipher,
			FragmentID: i + 1,
		}
		//logger.Logger.Println("Posting fragment ", fragment)
		res, err = PostFragment(
			fragment,
		)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed sending the post message fragment request: %v", i), "error", err)
			return err
		}
		err = ReadResponse(res, nil)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed posting the message fragment: %v", i), "error", err)
			return err
		}
	}
	return nil
}

func LogOut() error {
	res, err := PostLogout()
	if err != nil {
		logger.Logger.Errorln("Failed sending the logout request", "error", err)
		return err
	}
	sessionCookie = ""
	return ReadResponse(res, nil)
}

// "famePublicKey": security.FamePublicKey, "gpswPublicKey": security.GpswPublicKey
type KeysResponse struct {
	FamePublicKey *abe.FAMEPubKey `json:"famePublicKey"`
}

func LogIn(email string, password string) (*KeysResponse, error) {
	res, err := PostLogin(&routes.LoginReq{Email: email, Password: password})
	if err != nil {
		logger.Logger.Errorln("Failed sending the logout request", "error", err)
		return nil, err
	}
	keys := &KeysResponse{}
	err = ReadResponse(res, keys)
	if err != nil {
		logger.Logger.Errorln("Failed reading the login response", "error", err)
		return nil, err
	}
	return keys, nil
}

func LoadProfile(email string) (*models.User, error) {
	res, err := GetProfile(email)
	if err != nil {
		logger.Logger.Errorln("Failed sending the get profile request", "error", err)
	}
	user := &models.User{}
	err = ReadResponse(res, user)
	if err != nil {
		logger.Logger.Errorln("Failed reading the profile from the response", "error", err)
		return nil, err
	}
	return user, nil
}

func Emails() ([]string, error) {
	res, err := GetAllEmails()
	if err != nil {
		logger.Logger.Errorln("Failed sending the get all emails request", "error", err)
		return nil, err
	}
	var emails []string
	err = ReadResponse(res, &emails)
	if err != nil {
		logger.Logger.Errorln("Failed extracting the emails response", "error", err)
		return nil, err
	}
	return emails, nil

}
