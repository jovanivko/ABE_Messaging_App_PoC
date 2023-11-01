package routes

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/internal/security"
	"abeProofOfConcept/pkg/models"
	"abeProofOfConcept/pkg/store"
	"abeProofOfConcept/pkg/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func RegisterHandler(c *gin.Context) {
	var reg RegisterReq
	if err := c.BindJSON(&reg); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": "Failed to extract register data from body",
				"error":   err.Error(),
			},
		)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Errorln("Failed to hash password: ", err.Error())
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to hash password",
				"error":   err.Error(),
			},
		)
		return
	}
	//encrypt salary and address information
	salary, address, err := security.EncryptSensitiveInfo(reg.Email, reg.Department, reg.Salary, reg.Address)
	if err != nil {
		logger.Logger.Errorln("Failed to encrypt sensitive information: ", err.Error())
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to encrypt sensitive information",
				"error":   err.Error(),
			},
		)
		return
	}

	user, err := store.CreateUser(
		c, reg.Email, hashedPassword, reg.FirstName, reg.LastName, reg.Position, reg.Department,
		reg.PhoneNumber, salary, address,
	)
	if err != nil {
		logger.Logger.Errorln("Failed to create user: ", err.Error())
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to create user",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}

func LoginHandler(c *gin.Context) {
	var login LoginReq
	if err := c.BindJSON(&login); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": "Failed to extract login info from body",
				"error":   err.Error(),
			},
		)
		return
	}

	user, err := store.GetUser(c, login.Email)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": "There is no user with provided email",
				"error":   err.Error(),
			},
		)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password))
	if err != nil {
		c.JSON(
			http.StatusUnauthorized, gin.H{
				"message": "Password is incorrect",
				"error":   err.Error(),
			},
		)
		return
	}

	session := sessions.Default(c)
	session.Set("email", user.Email)
	session.Set("position", user.Position)
	session.Set("department", user.Department)
	err = session.Save()

	logger.Logger.Println("Session: ", session)
	logger.Logger.Println(session.Get("email"))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to save user session",
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"famePublicKey": security.FamePublicKey})
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func MessageHandler(c *gin.Context) {
	var msg MessageReq
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": "Failed extract message from request body",
				"error":   err.Error(),
			},
		)
		return
	}

	header, err := store.CreateMessage(c, msg.From, msg.To, msg.Title)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to create a new message",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, header.ID)
}

func FragmentHandler(c *gin.Context) {
	var frg FragmentReq
	if err := c.BindJSON(&frg); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": "Failed extract fragment from request body",
				"error":   err.Error(),
			},
		)
		return
	}
	msgFrg, err := store.CreateMessageFragment(c, frg.MsgID, frg.Content, frg.FragmentID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to create a new message fragment",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, msgFrg)
}

func MailboxHandler(c *gin.Context) {
	var messages []models.Message
	session := sessions.Default(c)

	//logger.Logger.Println("Session: ", session)
	//logger.Logger.Println(session.Get("email"))

	email := session.Get("email").(string)
	department := session.Get("department").(string)
	position := session.Get("position").(string)
	logger.Logger.Println(
		"Decrypting for attributes email:%v position:%v department:%v position:%v", email, department, position,
	)
	msgs, err := store.GetMessages(c, email)
	logger.Logger.Printf("Found %v messages\n", len(msgs))
	if err != nil {

		logger.Logger.Errorln("Failed to get messages", err.Error())
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Failed to get messages",
				"error":   err.Error(),
			},
		)
		return
	}
	for _, msg := range msgs {
		titleEncoded := &security.FameCipher{}
		err := utils.DecodeToObject(msg.Title, titleEncoded)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed to decode title for message: %v", msg.ID), err.Error())
			c.JSON(
				http.StatusInternalServerError, gin.H{
					"message": "Failed to decode title",
					"error":   err.Error(),
				},
			)
			return
		}
		title, err := security.DecryptFameCipher(titleEncoded, email, position, department)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed to decrypt title for message: %v", msg.ID), err.Error())
			c.JSON(
				http.StatusUnauthorized, gin.H{
					"message": "Failed to decrypt title",
					"error":   err.Error(),
				},
			)
			return
		}
		frgs, err := store.GetMessageFragments(c, msg.ID)
		logger.Logger.Printf("Found %v fragments for messageID %v\n", len(frgs), msg.ID)
		if err != nil {
			logger.Logger.Errorln(fmt.Sprintf("Failed to get message fragments for message: %v", msg.ID), err.Error())
			c.JSON(
				http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Failed to get message fragments for message: %v", msg.ID),
					"error":   err.Error(),
				},
			)
			return
		}
		var content string
		for _, frg := range frgs {
			frgContentEncrypted := &security.FameCipher{}
			err := utils.DecodeToObject(frg.Content, frgContentEncrypted)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError, gin.H{
						"message": fmt.Sprintf("Failed to decode message fragment: %v", frg.ID),
						"error":   err.Error(),
					},
				)
				return
			}
			frgContent, err := security.DecryptFameCipher(frgContentEncrypted, email, position, department)
			if err == nil {
				content += frgContent + "\r\n"
			} else {
				logger.Logger.Errorln(
					"message",
					fmt.Sprintf(
						"Failed to decrypt fragment %v for message %v", frg.FragmentID, frg.MsgID,
					),
					"error",
					err,
				)
			}
		}
		messages = append(
			messages, models.Message{
				From: msg.Sender, To: msg.Recipients, Title: title, Content: content, CreatedAt: msg.CreatedAt,
			},
		)
	}

	c.JSON(http.StatusOK, messages)
}

func ProfileHandler(c *gin.Context) {
	email := c.Param("email")
	user, err := store.GetUser(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "There is no user with provided email"})
		return
	}

	session := sessions.Default(c)
	email = session.Get("email").(string)
	position := session.Get("position").(string)
	department := session.Get("department").(string)
	salaryCipher := &security.GPSWCipher{}
	err = utils.DecodeToObject(user.Salary, salaryCipher)
	salaryString, err := security.DecryptGpswCipher(salaryCipher, email, position, department)
	var salary int
	if err != nil {
		salary = 0
	} else {
		salary, _ = strconv.Atoi(salaryString)
	}
	addressCipher := &security.GPSWCipher{}
	err = utils.DecodeToObject(user.Address, addressCipher)
	address, err := security.DecryptGpswCipher(addressCipher, email, position, department)
	if err != nil {
		address = ""
	}
	c.JSON(
		http.StatusOK, models.User{
			Email:       user.Email,
			Position:    user.Position,
			Department:  user.Department,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Address:     address,
			Salary:      salary,
		},
	)
}

func AllEmailsHandler(c *gin.Context) {
	var emails []string
	emails, err := store.GetAllUserEmails(c)
	if err != nil {
		logger.Logger.Errorln(
			"message", "Error getting all emails",
			"error", err.Error(),
		)
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"message": "Error getting all emails",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, emails)
}

// UserSession returns user session
func UserSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	return session
}
