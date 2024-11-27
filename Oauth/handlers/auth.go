// handlers/auth.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	url := h.OAuthConfig.AuthCodeURL("state",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
	//redirect to generated URL
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	//redirect to logout==success
	http.Redirect(w, r, "/?logout=success", http.StatusTemporaryRedirect)
}

func (h *Handler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Log the entire request for debugging
	fmt.Println("Callback URL Query:", r.URL.Query())

	// Get the authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code found", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for token
	token, err := h.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Token Exchange Error:", err)
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Print token details for debugging
	fmt.Println("Access Token:", token.AccessToken)
	//fmt.Println("Token Type:", token.TokenSource)
	fmt.Printf("Token Type: %T\n", token)

	// Fetch user info
	client := h.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Println("User Info Fetch Error:", err)
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

	// Read and log the raw response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	fmt.Println("Raw User Info Response:", string(bodyBytes))

	// Decode the user info
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		fmt.Println("JSON Decode Error:", err)
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Direct HTML response with user info
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Login Successful</title>
        </head>
        <body>
            <h1>Login Successful</h1>
            <p>Name: %s</p>
            <p>Email: %s</p>
           
			<img src="%s" alt="Profile Picture" width="100" height="100">
            <a href="/">Return to Home</a>
        </body>
        </html>
    `, userInfo.Name, userInfo.Email, userInfo.Picture)
}
