package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/adi290491/chirpy/internal/auth"
	"github.com/adi290491/chirpy/internal/database"
	"github.com/google/uuid"
)

var (
	profanes = map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func validateChirp(body string) (string, error) {

	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return cleanupChirps(body), nil
}

func cleanupChirps(body string) string {
	tmp := strings.Split(body, " ")
	for i, word := range tmp {
		if profanes[strings.ToLower(word)] {
			tmp[i] = "****"
		}
	}
	return strings.Join(tmp, " ")
}

func (c *apiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {

	authorization, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	uuid, err := auth.ValidateJWT(authorization, c.JWT_SECRET)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	type requestParams struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	var req requestParams
	err = decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	log.Printf("Request: %+v", req)

	cleanedBody, err := validateChirp(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	chirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uuid,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create a chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (c *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {

	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	var res []database.Chirp
	var err error
	if authorId == "" {
		res, err = c.db.GetAllChirps(r.Context())
	} else {
		userID, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error while parsing author id", err)
			return
		}
		res, err = c.db.GetAllChirpsByUserId(r.Context(), userID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while fetching all chirps", err)
	}

	sort.Slice(res, func(i, j int) bool {
		if sortOrder == "desc" {
			return res[j].CreatedAt.Before(res[i].CreatedAt)
		} else {
			return res[i].CreatedAt.Before(res[j].CreatedAt)
		}
	})

	var chirps []Chirp
	for _, row := range res {
		chirps = append(chirps, Chirp{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			Body:      row.Body,
			UserId:    row.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (c *apiConfig) GetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse chirp ID", err)
		return
	}

	chirp, err := c.db.GetChirpById(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "error while fetching chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (c *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	authorization, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	userID, err := auth.ValidateJWT(authorization, c.JWT_SECRET)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusForbidden, "could not parse chirp ID", err)
		return
	}

	chirp, err := c.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error while fetching chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you cannot delete another user's chirp", nil)
		return
	}

	err = c.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while deleting chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
