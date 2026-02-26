package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/domain"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/service"
)

type SpotifyHandler struct {
	service service.SpotifyService
}

func NewSpotifyHandler(s service.SpotifyService) *SpotifyHandler {
	return &SpotifyHandler{service: s}
}

// Helper para parsear los filtros comunes de la URL
func parseSpotifyFilters(r *http.Request) domain.SpotifyFilters {
	f := domain.SpotifyFilters{
		Search: r.URL.Query().Get("search"),
		Artist: r.URL.Query().Get("artist"),
		Track:  r.URL.Query().Get("track"),
	}

	// Cargar la zona horaria de Chile
	loc, _ := time.LoadLocation("America/Santiago")

	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		// Intentar parsear como fecha simple YYYY-MM-DD
		if t, err := time.ParseInLocation("2006-01-02", startStr, loc); err == nil {
			f.StartDate = &t
		}
	}
	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", endStr, loc); err == nil {
			// Para el EndDate, sumamos 23h 59m para incluir todo el día
			endOfDay := t.Add(24*time.Hour - time.Second)
			f.EndDate = &endOfDay
		}
	}
	if hStr := r.URL.Query().Get("start_hour"); hStr != "" {
		if h, err := strconv.Atoi(hStr); err == nil {
			f.StartHour = &h
		}
	}
	if hStr := r.URL.Query().Get("end_hour"); hStr != "" {
		if h, err := strconv.Atoi(hStr); err == nil {
			f.EndHour = &h
		}
	}
	return f
}

func (h *SpotifyHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetDashboardStats(r.Context(), parseSpotifyFilters(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (h *SpotifyHandler) GetTop(w http.ResponseWriter, r *http.Request) {
	f := parseSpotifyFilters(r)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	// El tipo (artists, songs, albums) viene de la URL
	listType := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

	res, err := h.service.GetTopList(r.Context(), listType, limit, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *SpotifyHandler) GetHabits(w http.ResponseWriter, r *http.Request) {
	// habit_type puede ser "time" o "dow" (day of week)
	hType := r.URL.Query().Get("type")
	res, err := h.service.GetHabitAnalysis(r.Context(), hType, parseSpotifyFilters(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *SpotifyHandler) GetEvolution(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetGlobalEvolution(r.Context(), parseSpotifyFilters(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *SpotifyHandler) GetYearly(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetYearlyStats(r.Context(), parseSpotifyFilters(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *SpotifyHandler) SearchRanking(w http.ResponseWriter, r *http.Request) {
	f := parseSpotifyFilters(r)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	target := domain.ArtistTrackFilters{
		Artist: r.URL.Query().Get("target_artist"),
		Track:  r.URL.Query().Get("target_track"),
	}
	res, err := h.service.SearchRankedItem(r.Context(), f, target, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *SpotifyHandler) GetWrapped(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	// 1. Validación de Año (Requerido para todos los Wrappeds)
	yearStr := q.Get("year")
	if yearStr == "" {
		http.Error(w, "El parámetro 'year' es obligatorio", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2010 || year > time.Now().Year()+1 {
		http.Error(w, "Año inválido", http.StatusBadRequest)
		return
	}

	var res interface{}
	var svcErr error

	// 2. Lógica de selección de Wrapped
	season := q.Get("season")
	monthStr := q.Get("month")

	if season != "" {
		// Caso Estacional
		res, svcErr = h.service.GetSeasonalWrapped(ctx, year, domain.Season(strings.ToLower(season)))
	} else if monthStr != "" {
		// Caso Mensual
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			http.Error(w, "El mes debe ser un número", http.StatusBadRequest)
			return
		}
		res, svcErr = h.service.GetMonthlyWrapped(ctx, year, month)
	} else {
		// Caso Anual por defecto
		res, svcErr = h.service.GetYearlyWrapped(ctx, year)
	}

	// 3. Manejo de errores del Servicio
	if svcErr != nil {
		http.Error(w, "Error al procesar el Wrapped: "+svcErr.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Respuesta Exitosa (El middleware JSONResponse se encarga del header)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}
