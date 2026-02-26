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
	GetTopSongs(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.SongRankingDTO, error)
	GetTopAlbums(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.AlbumRankingDTO, error)
	GetHabitsByTimeOfDay(ctx context.Context, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error)
	GetHabitsByDayOfWeek(ctx context.Context, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error)
	GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error)
	GetHistoryEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error)
	GetRankedSongs(ctx context.Context, f domain.SpotifyFilters, artistTrack domain.ArtistTrackFilters, limit int) ([]domain.SongRankingDTO, error)
	GetRankedArtist(ctx context.Context, f domain.SpotifyFilters, artist domain.ArtistTrackFilters, limit int) ([]domain.ArtistRankingDTO, error)
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
		clauses = append(clauses, fmt.Sprintf("(artist_name ILIKE $%d OR album_name ILIKE $%d OR track_name ILIKE $%d)", placeholder, placeholder, placeholder))
		args = append(args, "%"+f.Search+"%")
		placeholder++
	}
	if f.Artist != "" {
		clauses = append(clauses, fmt.Sprintf("artist_name ILIKE $%d", placeholder))
		args = append(args, "%"+f.Artist+"%")
		placeholder++
	}
	if f.Track != "" {
		clauses = append(clauses, fmt.Sprintf("track_name ILIKE $%d", placeholder))
		args = append(args, "%"+f.Track+"%")
		placeholder++
	}
	if f.StartHour != nil && f.EndHour != nil {
		clauses = append(clauses, fmt.Sprintf("EXTRACT(HOUR FROM ts) BETWEEN $%d AND $%d", placeholder, placeholder+1))
		args = append(args, *f.StartHour, *f.EndHour)
		placeholder += 2
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

// Filtros de artista y track,  afectan la VISUALIZACIÓN (Qué artista o canción quiero ver)
func buildWhereArtistTrackClause(f domain.ArtistTrackFilters, startPlaceholder int) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	p := startPlaceholder

	if f.Artist != "" {
		clauses = append(clauses, fmt.Sprintf("artist_name ILIKE $%d", p))
		args = append(args, "%"+f.Artist+"%")
		p++
	}
	if f.Track != "" {
		clauses = append(clauses, fmt.Sprintf("track_name ILIKE $%d", p))
		args = append(args, "%"+f.Track+"%")
		p++
	}

	if len(clauses) == 0 {
		return "", nil
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

// GetTotalStats obtiene horas totales y diversidad musical
func (r *spotifyRepo) GetTotalStats(ctx context.Context, f domain.SpotifyFilters) (domain.TotalStatsDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			COALESCE(ROUND(SUM(ms_played) / 3600000.0, 2), 0) as total_hours,
			COALESCE(ROUND(SUM(ms_played) / 60000.0, 2), 0) as total_minutes,
			COALESCE(ROUND(SUM(ms_played) / NULLIF(COUNT(DISTINCT ts::date), 0) / 3600000.0, 2), 0) AS average_daily_hours,
			COUNT(DISTINCT artist_name) as unique_artists,
			COUNT(DISTINCT track_name) as unique_songs
		FROM spotify_history %s`, where)

	var stats domain.TotalStatsDTO
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.TotalHours,
		&stats.TotalMinutes,
		&stats.AverageDailyHours,
		&stats.UniqueArtists,
		&stats.UniqueSongs,
	)
	return stats, err
}

// GetTopArtists obtiene el ranking de artistas
func (r *spotifyRepo) GetTopArtists(ctx context.Context, limit int, f domain.SpotifyFilters) ([]domain.ArtistRankingDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT 
			RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
			artist_name, 
			COALESCE(ROUND(SUM(ms_played) / 60000.0, 2), 0) as minutes_played,
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
// Util para wrappeds segun anio, mes, y estaciones del anio (capa service) LIMIT 100
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

// Momentos del dia por bloque horario, cantidad de escuchas
func (r *spotifyRepo) GetHabitsByTimeOfDay(ctx context.Context, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
        SELECT 
            CASE 
                WHEN EXTRACT(HOUR FROM ts) BETWEEN 6 AND 11 THEN 'Mañana'
                WHEN EXTRACT(HOUR FROM ts) BETWEEN 12 AND 17 THEN 'Tarde'
                WHEN EXTRACT(HOUR FROM ts) BETWEEN 18 AND 23 THEN 'Noche'
                ELSE 'Madrugada'
            END AS label,
            COUNT(*) AS count
        FROM spotify_history %s
        GROUP BY label 
		ORDER BY count DESC`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.HabitTimeDTO
	for rows.Next() {
		var d domain.HabitTimeDTO
		if err := rows.Scan(&d.Label, &d.Count); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

// Escuchas segun dia de la semana (ingles)
func (r *spotifyRepo) GetHabitsByDayOfWeek(ctx context.Context, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
        SELECT 
			EXTRACT(DOW FROM ts) AS num_day, 
			COUNT(*) AS count
        FROM spotify_history %s
        GROUP BY EXTRACT(DOW FROM ts)
        ORDER BY EXTRACT(DOW FROM ts)`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.HabitTimeDTO
	for rows.Next() {
		var d domain.HabitTimeDTO
		if err := rows.Scan(&d.NumDay, &d.Count); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

// Comparativa anual (Tu año en música)
func (r *spotifyRepo) GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
        SELECT 
            EXTRACT(YEAR FROM ts)::int AS year,
            COALESCE(ROUND(SUM(ms_played) / 3600000.0, 2), 0) AS total_hours,
            COALESCE(ROUND(SUM(ms_played) / 60000.0, 2), 0) AS total_minutes,
            COUNT(*) AS total_songs
        FROM spotify_history 
        %s
        GROUP BY year ORDER BY year`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.YearlyStatsDTO
	for rows.Next() {
		var d domain.YearlyStatsDTO
		if err := rows.Scan(&d.Year, &d.TotalHours, &d.TotalMinutes, &d.TotalSongs); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

// Evolucion historica mensual (Grafico lineas)
func (r *spotifyRepo) GetHistoryEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
		SELECT
			TO_CHAR(ts, 'YYYY') AS year,
			TO_CHAR(ts, 'MM') AS month,
			TO_CHAR(ts, 'YYYY-MM') AS year_month,
			COALESCE(SUM(ms_played) / 3600000.0, 0) AS hours_monthly,
			COALESCE(SUM(ms_played) / 60000.0, 0) AS minutes_monthly
		FROM spotify_history
		%s 
		GROUP BY year, month, year_month
		ORDER BY year, month;`, where)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resul []domain.HistoryEvolutionDTO
	for rows.Next() {
		var r domain.HistoryEvolutionDTO
		if err := rows.Scan(&r.Year, &r.Month, &r.YearMonth, &r.HoursMonthly, &r.MinutesMonthly); err != nil {
			return nil, err
		}
		resul = append(resul, r)
	}
	return resul, nil
}

func (r *spotifyRepo) GetRankedSongs(ctx context.Context, f domain.SpotifyFilters, artistTrack domain.ArtistTrackFilters, limit int) ([]domain.SongRankingDTO, error) {
	// 1. Filtros base (van dentro del ranking para acotar el tiempo/duración)
	baseWhere, baseArgs := buildWhereClause(f)

	// 2. Filtros de selección (van fuera para filtrar el resultado final)
	// Estos no cambian el cálculo del ranking, solo qué filas se muestran
	finalWhere, finalArgs := buildWhereArtistTrackClause(artistTrack, len(baseArgs)+1)

	allArgs := append(baseArgs, finalArgs...)

	query := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT
                RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
                track_name,
                artist_name,
                COUNT(*) AS times_played
            FROM spotify_history
            %s
            GROUP BY track_name, artist_name
        )
        SELECT * FROM ranking_completo
        %s
        ORDER BY ranking ASC
		LIMIT $%d`, baseWhere, finalWhere, len(allArgs)+1)

	allArgs = append(allArgs, limit)
	rows, err := r.db.Query(ctx, query, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resul []domain.SongRankingDTO
	for rows.Next() {
		var r domain.SongRankingDTO
		if err := rows.Scan(&r.Ranking, &r.TrackName, &r.ArtistName, &r.TimesPlayed); err != nil {
			return nil, err
		}
		resul = append(resul, r)
	}
	return resul, nil
}

func (r *spotifyRepo) GetRankedArtist(ctx context.Context, f domain.SpotifyFilters, artist domain.ArtistTrackFilters, limit int) ([]domain.ArtistRankingDTO, error) {
	baseWhere, baseArgs := buildWhereClause(f)
	finalWhere, finalArgs := buildWhereArtistTrackClause(artist, len(baseArgs)+1)
	allArgs := append(baseArgs, finalArgs...)

	query := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT
                RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
                artist_name,
				COALESCE(ROUND(SUM(ms_played) / 60000.0, 2), 0) as minutes_played,
                COUNT(*) AS times_played
            FROM spotify_history
            %s
            GROUP BY artist_name
        )
        SELECT * FROM ranking_completo
        %s
        ORDER BY ranking ASC
		LIMIT $%d`, baseWhere, finalWhere, len(allArgs)+1)

	allArgs = append(allArgs, limit)
	rows, err := r.db.Query(ctx, query, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resul []domain.ArtistRankingDTO
	for rows.Next() {
		var r domain.ArtistRankingDTO
		if err := rows.Scan(&r.Ranking, &r.ArtistName, &r.MinutesPlayed, &r.TimesPlayed); err != nil {
			return nil, err
		}
		resul = append(resul, r)
	}
	return resul, nil
}
