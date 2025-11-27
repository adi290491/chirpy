package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/adi290491/chirpy/internal/auth"
	"github.com/adi290491/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	JWT_SECRET     string
	PLATFORM       string
	API_KEY        string
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) GetFileServerHits(w http.ResponseWriter, r *http.Request) {
	hits := c.fileServerHits.Load()

	template := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, hits)
	log.Printf("Hits Fetched: %d", hits)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, template)

}

func (c *apiConfig) ResetFileServerHits(w http.ResponseWriter, r *http.Request) {

	if c.PLATFORM != "dev" {
		respondWithError(w, http.StatusForbidden, "", errors.New("operation not permitted"))
		return
	}

	err := c.db.DeleteUser(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while deleting users", err)
		return
	}
	// c.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (c *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	type RefreshTokenResponse struct {
		Token string `json:"token"`
	}

	user, err := c.db.GetUserFromRefreshToken(r.Context(), token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token is invalid or has expired", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, c.JWT_SECRET, expirationTime*time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while refreshing token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshTokenResponse{
		Token: newToken,
	})
}

func (c *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	err = c.db.UpdateRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while revoking refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (c *apiConfig) initDB() {
	InitDB(c)
}
