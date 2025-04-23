package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rhajizada/gazette/internal/oauth"
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
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
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

	var claims oauth.Claims
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "could not parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	appToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims.GetJwtMap(time.Hour),
	).SignedString(h.Secret)
	if err != nil {
		msg := fmt.Sprintf("failed to sign app token: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "app_jwt",
		Value:    appToken,
		Path:     "/api/feeds/",
		HttpOnly: true,
		MaxAge:   3600,
	})
	http.Redirect(w, r, "/api/feeds/", http.StatusSeeOther)
}
