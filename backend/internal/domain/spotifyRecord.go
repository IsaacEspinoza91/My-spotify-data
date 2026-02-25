package domain

import "time"

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

// Filtros de búsqueda
type SpotifyFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	Search    string // Para artista o álbum
	StartHour *int   // 0-23
	EndHour   *int   // 0-23
}
