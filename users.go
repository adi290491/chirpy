package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adi290491/chirpy/internal/auth"
	"github.com/adi290491/chirpy/internal/database"
	"github.com/google/uuid"
)

var (
	expirationTime time.Duration = 1
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (c *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {

	var params LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hash error", err)
		return
	}

	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hash,
		Email:          params.Email,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User creation error", err)
		return
	}

	resp := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func (c *apiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {

	var userRequest LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := c.db.GetUserByEmail(r.Context(), userRequest.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	log.Printf("USER: %+v", user)
	match, err := auth.CheckPasswordHash(userRequest.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	jwt, err := auth.MakeJWT(user.ID, c.JWT_SECRET, time.Hour*expirationTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JWT token creation failed", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Refresh token creation failed", err)
		return
	}

	c.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (c *apiConfig) UpdateEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Access Token is missing or malformed", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, c.JWT_SECRET)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	var userRequest LoginRequest
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	hashPassword, err := auth.HashPassword(userRequest.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hash error", err)
		return
	}

	user, err := c.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          userRequest.Email,
		HashedPassword: hashPassword,
		ID:             userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User update error", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed,
	})
}
