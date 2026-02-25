package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SpotifyRepository define las operaciones de base de datos
type SpotifyRepository interface {
	GetTotalStats(ctx context.Context, f domain.SpotifyFilters) (domain.TotalStatsDTO, error)
	GetTopArtists(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.ArtistRankingDTO, error)
}

type spotifyRepo struct {
	db *pgxpool.Pool
}

func NewSpotifyRepository(db *pgxpool.Pool) SpotifyRepository {
	return &spotifyRepo{db: db}
}

// Función auxiliar para construir WHERE dinámico
// Solo parametro search es obligatorio, pero puede ser "" para no filtrar por busqueda
func buildWhereClause(f domain.SpotifyFilters) (string, []interface{}) {
	clauses := []string{"spotify_uri LIKE 'spotify:track:%'", "ms_played > 10000"}
	args := []interface{}{}
	placeholder := 1

	if f.StartDate != nil {
		clauses = append(clauses, fmt.Sprintf("ts >= $%d", placeholder))
		args = append(args, *f.StartDate)
		placeholder++
	}
	if f.EndDate != nil {
		clauses = append(clauses, fmt.Sprintf("ts <= $%d", placeholder))
		args = append(args, *f.EndDate)
		placeholder++
	}
	if f.Search != "" {
		clauses = append(clauses, fmt.Sprintf("(artist_name ILIKE $%d OR album_name ILIKE $%d)", placeholder, placeholder))
		args = append(args, "%"+f.Search+"%")
		placeholder++
	}
	if f.StartHour != nil && f.EndHour != nil {
		clauses = append(clauses, fmt.Sprintf("EXTRACT(HOUR FROM ts) BETWEEN $%d AND $%d", placeholder, placeholder+1))
		args = append(args, *f.StartHour, *f.EndHour)
		placeholder += 2
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

// GetTotalStats obtiene horas totales y diversidad musical
func (r *spotifyRepo) GetTotalStats(ctx context.Context, f domain.SpotifyFilters) (domain.TotalStatsDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			ROUND(SUM(ms_played) / 3600000.0, 2) as total_hours,
			ROUND(SUM(ms_played) / 60000.0, 2) as total_minutes,
			ROUND(SUM(ms_played) / (COUNT(DISTINCT ts::date) * 3600000.0), 2) AS average_daily_hours,
			COUNT(DISTINCT artist_name) as unique_artists,
			COUNT(DISTINCT track_name) as unique_songs
		FROM spotify_history %s`, where)

	var stats domain.TotalStatsDTO
	err := r.db.QueryRow(ctx, query, args...).Scan(&stats.TotalHours, &stats.TotalMinutes, &stats.AverageDailyHours, &stats.UniqueArtists, &stats.UniqueSongs)
	return stats, err
}

// GetTopArtists obtiene el ranking de artistas
func (r *spotifyRepo) GetTopArtists(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.ArtistRankingDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
			artist_name, 
			ROUND(SUM(ms_played) / 60000.0, 2) as minutes_played,
			COUNT(*) as times_played
		FROM spotify_history 
		%s
		GROUP BY artist_name
		ORDER BY minutes_played DESC
		LIMIT $%d`, where, len(args)+1)

	args = append(args, limit)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []domain.ArtistRankingDTO
	for rows.Next() {
		var dto domain.ArtistRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.ArtistName, &dto.MinutesPlayed, &dto.TimesPlayed); err != nil {
			return nil, err
		}
		rankings = append(rankings, dto)
	}
	return rankings, nil
}

// GetTopAlbums obtiene el ranking de álbumes
// Util para wrappeds segun anio, mes, y estaciones del anio (capa service)
func (r *spotifyRepo) GetTopSongs(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.SongRankingDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking, 
			track_name, 
			artist_name, 
			COUNT(*) AS times_played
		FROM spotify_history
		%s
		GROUP BY track_name, artist_name
		ORDER BY times_played DESC
		LIMIT $%d`, where, len(args)+1)
	args = append(args, limit)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []domain.SongRankingDTO
	for rows.Next() {
		var dto domain.SongRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.TrackName, &dto.ArtistName, &dto.TimesPlayed); err != nil {
			return nil, err
		}
		rankings = append(rankings, dto)
	}
	return rankings, nil
}

// GetTopAlbums obtiene el ranking de álbumes
func (r *spotifyRepo) GetTopAlbums(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.AlbumRankingDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking, 
			album_name, 
			artist_name, 
			COUNT(*) AS times_played
		FROM spotify_history
		%s
		GROUP BY album_name, artist_name
		ORDER BY times_played DESC
		LIMIT $%d`, where, len(args)+1)
	args = append(args, limit)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []domain.AlbumRankingDTO
	for rows.Next() {
		var dto domain.AlbumRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.AlbumName, &dto.ArtistName, &dto.TimesPlayed); err != nil {
			return nil, err
		}
		rankings = append(rankings, dto)
	}
	return rankings, nil
}
