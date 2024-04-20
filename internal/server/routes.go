package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	s.Store = session.New(session.Config{
		Expiration:   3600,
		CookieSecure: false, //if true , it will only send to https and not http
	})
	s.App.Get("/auth/google", s.authGoogleHandler)
	s.App.Get("/auth/google/callback", s.authGoogleCallbackHandler)
	s.App.Get("/profile", s.profileHandler)
}

func (s *FiberServer) authGoogleHandler(c *fiber.Ctx) error {
	from := c.Query("from", "/")
	url := googleOauthConfig.AuthCodeURL(from)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (s *FiberServer) authGoogleCallbackHandler(c *fiber.Ctx) error {
	log.Info("context", c)

	var pathUrl string = "/"
	state := c.Query("state")
	if state != "" {
		pathUrl = state
	}

	code := c.Query("code")
	if code == "" {
		c.SendStatus(http.StatusUnauthorized)
	}

	//get token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.SendStatus(http.StatusUnauthorized)
	}

	//exchange token with userInfo
	userInfoJSON, err := fetchGoogleUserInfo(token)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}

	//save userInfo to session
	//Because session cant keep large files , so json needs to be turn to bytes
	jsonMarshal, err := json.Marshal(userInfoJSON)
	if err != nil {
		return err
	}
	session, err := s.Store.Get(c)
	if err != nil {
		return err
	}
	session.Set("user", jsonMarshal)

	if err := session.Save(); err != nil {
		return err
	}

	return c.Redirect(fmt.Sprint("http://localhost:8000", pathUrl), http.StatusTemporaryRedirect)
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

func (s *FiberServer) profileHandler(c *fiber.Ctx) error {
	session, err := s.Store.Get(c)
	if err != nil {
		return err
	}

	//get from session key "user"
	userMarshal := session.Get("user")
	if userMarshal == nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	//unMarshal the userInfo
	var userJSON interface{}
	if err := json.Unmarshal(userMarshal.([]byte), &userJSON); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"data":   userJSON,
	})
}
