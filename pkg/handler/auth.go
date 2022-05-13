package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"

	"github.com/MelchiSalins/go-auth/models"
	"github.com/MelchiSalins/go-auth/pkg/app"
)

var (
	jwtKey = []byte(app.JwtSecret)
)

const (
	expiresIn = 60
)

// CustomClaims holds all standards claims + coinfish specific claims
type CustomClaims struct {
	User  string `json:"user"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// LoginResponse is the JSON payload that is sent back after login attempt
type LoginResponse struct {
	Status      bool    `json:"Status"`
	TokenType   string  `json:"TokenType"`
	AccessToken *string `json:"AccessToken"`
	ExpiresIn   int     `json:"ExpiresIn"`
}

// Authenticator Struct for Oauth2 Authentication
type Authenticator struct {
	Provider     *oidc.Provider
	ClientConfig oauth2.Config
	Ctx          context.Context
}

// NewAuthenticator Return an instance of Authenticator type
func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, app.OAuthIssuer)
	if err != nil {
		log.Printf("Failed to get provider: %v", err)
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  app.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider:     provider,
		ClientConfig: config,
		Ctx:          ctx,
	}, nil
}

// HandleCallback handles the OAuth provider redirect
func (a *Authenticator) HandleCallback(w http.ResponseWriter, r *http.Request) {
	//TODO: The state should be randomised or made a secret

	// Checks the state to avoid CSRF attacks
	if r.URL.Query().Get("state") != "state" {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// Gets the Authorization code from the URL and Exchanges it for a token,
	//token contains AccessToken, Token Type, Refresh Token, Expiry & Raw Meta Data.
	token, err := a.ClientConfig.Exchange(a.Ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawTokenID, ok := token.Extra("id_token").(string)
	log.Println(rawTokenID)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: app.ClientID,
	}

	//Below step verifies if the rawTokenID (JWT) is signed and valid.
	idToken, err := a.Provider.Verifier(oidcConfig).Verify(a.Ctx, rawTokenID)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idtc := new(models.User)
	if err := idToken.Claims(&idtc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	us, err := models.NewUserService()
	// TODO: Check if below causes issues with gorm internal connection pool
	defer us.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if this is a returning user or new registration
	_, err = us.ByEmail(idtc.Email)

	// check for Error: record not found
	if err != nil {
		// fmt.Println("New User") // log this event
		if err = us.ExistOrCreate(idtc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	j, err := issueJWT(idtc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lr := &LoginResponse{
		AccessToken: j,
		ExpiresIn:   expiresIn,
		Status:      true,
		TokenType:   "Bearer",
	}
	lrJSON, err := json.Marshal(lr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", lrJSON) // Returns the created record Email ID.
}

// issueJWT issues custom Json Web Token after successful OIDC login (Google)
// This token will have claims specific to the coinfish.
func issueJWT(tk *models.User) (*string, error) {
	expiresAt := time.Now().Add(expiresIn * time.Minute)

	standardClaims := jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		NotBefore: time.Now().Unix(),
		IssuedAt:  time.Now().Unix(),
		Audience:  "mflow",
		Issuer:    "mflow",
		Subject:   "mflow",
	}
	claims := &CustomClaims{
		User:           tk.Name,
		Email:          tk.Email,
		StandardClaims: standardClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &tokenString, nil
}
