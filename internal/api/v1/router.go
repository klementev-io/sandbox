package v1

import "net/http"

func SetupRouter() http.Handler {
	h := new(Handler)

	router := http.NewServeMux()

	router.HandleFunc("/api/v1/health", h.Health)

	return router
}
