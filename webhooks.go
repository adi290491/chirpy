package main

import (
	"encoding/json"
	"net/http"

	"github.com/adi290491/chirpy/internal/auth"
	"github.com/adi290491/chirpy/internal/database"
	"github.com/google/uuid"
)

type Webhook struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	UserID uuid.UUID `json:"user_id"`
}

var eventType = map[string]struct{}{
	"user.upgraded": {},
}

func (c *apiConfig) RunWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	if apiKey != c.API_KEY {
		respondWithError(w, http.StatusUnauthorized, "api key does not match", nil)
		return
	}

	var webhook Webhook
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhook)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while processing request", err)
		return
	}

	if _, ok := eventType[webhook.Event]; !ok {
		respondWithError(w, http.StatusNoContent, "unrecognized event", nil)
		return
	}

	user, err := c.db.GetUserByID(r.Context(), webhook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error while fetching user", err)
		return
	}

	err = c.db.UpdateUserSubscription(r.Context(), database.UpdateUserSubscriptionParams{
		IsChirpyRed: true,
		ID:          user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, "error while updating user subscription", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
