package main

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/internal/security"
	routes "abeProofOfConcept/pkg/routes"
	"abeProofOfConcept/pkg/store"
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"time"
)

func main() {
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			logger.Logger.Fatalln("Failed to close database")
		}
	}(store.ConnectDB())
	err := store.CreateDatabase(context.Background())
	if err != nil {
		logger.Logger.Fatalln("Failed to create database")
	}
	r := gin.Default()
	if err = security.InitFame(); err != nil {
		logger.Logger.Fatalln("Failed to initialize FAME security: ", err.Error())
	}
	if err = security.InitGpsw(); err != nil {
		logger.Logger.Fatalln("Failed to initialize GPSW security: ", err.Error())
	}

	cookieStore := cookie.NewStore([]byte("secret"))
	cookieStore.Options(
		sessions.Options{
			MaxAge: int(15 * time.Minute), // Set session timeout to 15 minutes
			Path:   "/",
		},
	)
	r.Use(sessions.Sessions("mysession", cookieStore))
	r.Use(logger.RequestLogger())

	r.POST("/register", routes.RegisterHandler)
	r.POST("/login", routes.LoginHandler)
	r.POST("/logout", routes.LogoutHandler)
	r.POST("/message", routes.MessageHandler)
	r.POST("/fragment", routes.FragmentHandler)
	r.GET("/mailbox", routes.MailboxHandler)
	r.GET("/profile/:email", routes.ProfileHandler)
	r.GET("/emails", routes.AllEmailsHandler)

	r.Run(":8080")
}
