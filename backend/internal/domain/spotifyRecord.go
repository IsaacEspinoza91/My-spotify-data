package domain

import (
	"strings"
	"time"
)

// Entidad de Base de Datos
type SpotifyRecord struct {
	ID          int       `json:"id"`
	TS          time.Time `json:"ts"`
	Platform    string    `json:"platform"`
	MsPlayed    int       `json:"ms_played"`
	ConnCountry string    `json:"conn_country"`
	TrackName   string    `json:"track_name"`
	ArtistName  string    `json:"artist_name"`
	AlbumName   string    `json:"album_name"`
	SpotifyURI  string    `json:"spotify_uri"`
}

// DTO para Estadísticas Generales
type TotalStatsDTO struct {
	TotalHours        float64 `json:"total_hours"`
	TotalMinutes      float64 `json:"total_minutes"`
	AverageDailyHours float64 `json:"average_daily_hours"`
	UniqueArtists     int     `json:"unique_artists"`
	UniqueSongs       int     `json:"unique_songs"`
}

// DTO para Rankings
type ArtistRankingDTO struct {
	Ranking       int     `json:"ranking"`
	ArtistName    string  `json:"artist_name"`
	MinutesPlayed float64 `json:"minutes_played"`
	TimesPlayed   int     `json:"times_played"`
}

type SongRankingDTO struct {
	Ranking     int    `json:"ranking"`
	TrackName   string `json:"track_name"`
	ArtistName  string `json:"artist_name"`
	TimesPlayed int    `json:"times_played"`
}

type AlbumRankingDTO struct {
	Ranking     int    `json:"ranking"`
	AlbumName   string `json:"album_name"`
	ArtistName  string `json:"artist_name"`
	TimesPlayed int    `json:"times_played"`
}

type HabitTimeDTO struct {
	Label  string `json:"label,omitempty"` // Mañana, Tarde, Lunes, Martes, 2023, etc.
	NumDay *int   `json:"num_day,omitempty"`
	Count  int    `json:"count"`
}

type YearlyStatsDTO struct {
	Year         int     `json:"year"`
	TotalHours   float64 `json:"total_hours"`
	TotalMinutes float64 `json:"total_minutes"`
	TotalSongs   int     `json:"total_songs"`
}

type HistoryEvolutionDTO struct {
	Year           string  `json:"year"`       // YYYY
	Month          string  `json:"month"`      // MM
	YearMonth      string  `json:"year_month"` // YYYY-MM
	HoursMonthly   float64 `json:"hours_monthly"`
	MinutesMonthly float64 `json:"minutes_monthly"`
}

// Filtros de búsqueda
type SpotifyFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	Search    string // Para artista o álbum
	Artist    string // Filtro específico
	Track     string // Filtro específico
	StartHour *int   // 0-23
	EndHour   *int   // 0-23
	Page      int
	Limit     int
}
type ArtistTrackFilters struct {
	Artist string
	Track  string
}

// Limpieza y validación de filtros
func (f *SpotifyFilters) CleanAndValidate() {
	// 1. Trim de strings para evitar espacios accidentales
	f.Search = strings.TrimSpace(f.Search)
	f.Artist = strings.TrimSpace(f.Artist)
	f.Track = strings.TrimSpace(f.Track)

	// 2. Validación de rango de horas
	if f.StartHour != nil {
		if *f.StartHour < 0 {
			*f.StartHour = 0
		}
		if *f.StartHour > 23 {
			*f.StartHour = 23
		}
	}
	if f.EndHour != nil {
		if *f.EndHour < 0 {
			*f.EndHour = 0
		}
		if *f.EndHour > 23 {
			*f.EndHour = 23
		}
	}

	// 3. Validación de lógica temporal (opcional)
	if f.StartDate != nil && f.EndDate != nil {
		if f.StartDate.After(*f.EndDate) {
			// Si la fecha inicio es mayor a la fin, podrías resetearlas o swapearlas
			f.StartDate, f.EndDate = f.EndDate, f.StartDate
		}
	}

	// 4. Validaciones paginacion
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}
}

// Offset calcula el salto para SQL
func (f *SpotifyFilters) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *ArtistTrackFilters) Clean() {
	f.Artist = strings.TrimSpace(f.Artist)
	f.Track = strings.TrimSpace(f.Track)
}

type Season string

const (
	Summer Season = "summer"
	Autumn Season = "autumn"
	Winter Season = "winter"
	Spring Season = "spring"
)
