package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rhajizada/gazette/internal/oauth"
	"github.com/rhajizada/gazette/internal/repository"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := randomState()
	if err != nil {
		msg := fmt.Sprintf("failed generating state: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300,
	})
	url := h.OAuth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oidc_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	oauthToken, err := h.OAuth.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		msg := fmt.Sprintf("token exchange failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	rawID := oauthToken.Extra("id_token")
	if rawID == nil {
		http.Error(w, "no 'id_token' in token response", http.StatusInternalServerError)
		return
	}
	idToken, err := h.Verifier.Verify(r.Context(), rawID.(string))
	if err != nil {
		msg := fmt.Sprintf("invalid id_token: %v", err)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	var claims oauth.ProviderClaims
	if err := idToken.Claims(&claims); err != nil {
		msg := fmt.Sprintf("could not parse claims: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	user, err := h.Repo.GetUserBySub(r.Context(), claims.Sub)
	userExists := true
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userExists = false
		} else {
			msg := fmt.Sprintf("failed fetching user: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	if !userExists {
		user, err = h.Repo.CreateUser(r.Context(), repository.CreateUserParams{
			Sub:   claims.Sub,
			Name:  claims.Name,
			Email: claims.Email,
		})
		if err != nil {
			msg := fmt.Sprintf("failed creating user: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	appToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims.GetAppClaims(user.ID, time.Hour),
	).SignedString(h.Secret)
	if err != nil {
		msg := fmt.Sprintf("failed to sign app token: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "app_jwt",
		Value:    appToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	baseURL := "/"
	params := url.Values{}
	params.Add("token", appToken)

	url, err := url.Parse(baseURL)
	if err != nil {
		msg := fmt.Sprintf("failed to sign app token: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	url.RawQuery = params.Encode()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}
