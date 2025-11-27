package main

import "net/http"

func (apiCfg *apiConfig) registerRoutes(mux *http.ServeMux, handler http.Handler) {
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.GetFileServerHits)

	mux.HandleFunc("POST /admin/reset", apiCfg.ResetFileServerHits)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)

	mux.HandleFunc("POST /api/login", apiCfg.LoginUser)

	mux.HandleFunc("PUT /api/users", apiCfg.UpdateEmailAndPassword)

	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshToken)

	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeToken)

	mux.HandleFunc("POST /api/chirps", apiCfg.CreateChirp)

	mux.HandleFunc("GET /api/chirps", apiCfg.GetAllChirps)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpById)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.RunWebhook)

}
