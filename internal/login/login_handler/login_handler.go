package login_handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Server struct {
	router       *gin.Engine
	oauth2Config *oauth2.Config
}

func NewServer(oauth2Config *oauth2.Config) (*Server, error) {
	oauth2Config, err := NewOauth2Config()
	if err != nil {
		return nil, fmt.Errorf("Failed to get oauth2 config: %v", err)
	}
	s := &Server{
		router:       gin.Default(),
		oauth2Config: oauth2Config,
	}
	return s, nil
}

func NewOauth2Config() (*oauth2.Config, error) {
	providerURL := fmt.Sprintf("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	provider, err := oidc.NewProvider(context.Background(), providerURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get provider: %v", err)
	}
	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "photo"},
	}
	return oauth2Config, nil
}

func (s *Server) LoginHandler(c *gin.Context) {
	state, err := generateRandomString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Save State value in session storage
	session := sessions.Default(c)
	session.Set("state", state)

	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, s.oauth2Config.AuthCodeURL(state))
}

// logout handler
func (s *Server) LogoutHandler(c *gin.Context) {
	// clear session storage
	c.SetCookie("user", "", -1, "/", "localhost", false, true)
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.SetCookie("auth-session", "", -1, "/", "localhost", false, true)

	//call auth2 logout
	logoutURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse logout url"})
		return
	}

	//check request http or https
	scheme := "http"
	if c.Request.TLS == nil {
		scheme = "https"
	}

	//redirect to homepage
	redurectionURL, err := url.Parse(scheme + "://" + c.Request.Host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse url"})
		return

	}

	//add url params
	parameters := url.Values{}
	parameters.Add("returnTo", redurectionURL.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))

	logoutURL.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}

func (s *Server) CallbackHandler(c *gin.Context) {

	// check if state match
	session := sessions.Default(c)
	state := session.Get("state")
	if state != c.Request.URL.Query().Get("state") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := s.oauth2Config.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	if !token.Valid() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
		return
	}

	// get user information to display in profile
	client := s.oauth2Config.Client(c, token)
	resp, err := client.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})

	}

	//parse response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
	}

	// save user information in session storage
	c.SetCookie("user", string(b), int(time.Now().Add(1*time.Hour).Unix()), "/", "localhost", true, true)
	c.SetCookie("token", token.AccessToken, int(time.Now().Add(1*time.Hour).Unix()), "/", "localhost", true, true)

	c.Redirect(http.StatusTemporaryRedirect, "/index")
}

func generateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Failed to read random: %v", err)
	}
	s := base64.StdEncoding.EncodeToString(b)
	return s, nil
}

// implement auth middleware to make sure user is authenticated before taking to protected endpoints
func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		//check if token is empty
		if c.MustGet("token") == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
