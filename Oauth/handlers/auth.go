// handlers/auth.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sidgupt12/learning/Oauth/config"
	"golang.org/x/oauth2"
)

type Handler struct {
	OAuthConfig *oauth2.Config
}

func NewHandler() *Handler {
	return &Handler{
		OAuthConfig: config.GetGoogleOAuthConfig(),
	}
}

func (h *Handler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate authorization URL
	url := h.OAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Exchange authorization code for token
	token, err := h.OAuthConfig.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch user info
	client := h.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse user information
	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Simple response (you'd typically create a session or JWT here)
	fmt.Fprintf(w, "Welcome %s! Your email is %s", userInfo.Name, userInfo.Email)
}
