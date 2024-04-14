package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	Scopes:       []string{"profile", "email"},
	Endpoint:     google.Endpoint,
}

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/auth/google", s.authGoogleHandler)
	s.App.Get("/auth/google/callback", s.authGoogleCallbackHandler)
}

func (s *FiberServer) authGoogleHandler(c *fiber.Ctx) error {
	from := c.Query("from", "/")
	url := googleOauthConfig.AuthCodeURL(from)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (s *FiberServer) authGoogleCallbackHandler(c *fiber.Ctx) error {
	log.Info("context", c)
	//state := c.Query("state")

	//var urlPath string = "/"
	//if state != "" {
	//	urlPath = state
	//}

	code := c.Query("code")
	if code == "" {
		c.SendStatus(http.StatusUnauthorized)
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.SendStatus(http.StatusUnauthorized)
	}

	userInfoJSON, err := fetchGoogleUserInfo(token)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"status": "Success",
		"data":   userInfoJSON,
	})
}

func fetchGoogleUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := googleOauthConfig.Client(context.Background(), token)

	//Fetch api to google , with the user token
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	//Turns response to JSON
	userInfoJSON := make(map[string]interface{})
	if err := json.NewDecoder(response.Body).Decode(&userInfoJSON); err != nil {
		return nil, err
	}

	fmt.Println(userInfoJSON)
	return userInfoJSON, nil
}
