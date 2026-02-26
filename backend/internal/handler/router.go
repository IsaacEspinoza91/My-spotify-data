package handler

import (
	"net/http"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/service"
)

func NewRouter(spotifySvc service.SpotifyService) http.Handler {
	mux := http.NewServeMux()
	h := NewSpotifyHandler(spotifySvc)

	// 1. Estadísticas Generales
	mux.HandleFunc("GET /api/v1/spotify/stats", h.GetStats)

	// 2. Rankings (Top List)
	mux.HandleFunc("GET /api/v1/spotify/top/artists", h.GetTop)
	mux.HandleFunc("GET /api/v1/spotify/top/songs", h.GetTop)
	mux.HandleFunc("GET /api/v1/spotify/top/albums", h.GetTop)

	// 3. Hábitos (type=time o type=dow)
	mux.HandleFunc("GET /api/v1/spotify/habits", h.GetHabits)

	// 4. Evolución Mensual
	mux.HandleFunc("GET /api/v1/spotify/evolution", h.GetEvolution)

	// 5. Stats Anuales
	mux.HandleFunc("GET /api/v1/spotify/yearly", h.GetYearly)

	// 6. Búsqueda de Ranking Específico
	mux.HandleFunc("GET /api/v1/spotify/search-rank", h.SearchRanking)

	// 7. Wrappeds
	mux.HandleFunc("GET /api/v1/spotify/wrapped", h.GetWrapped)

	var handler http.Handler = mux
	handler = JSONResponse(handler)
	handler = Logger(handler)
	handler = CORS(handler)
	handler = Recovery(handler)

	return handler
}
