package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpotifyRepository interface {
	GetTotalStats(ctx context.Context) (domain.TotalStatsDTO, error)
	GetTopArtists(ctx context.Context, f domain.SpotifyFilters) ([]domain.ArtistRankingDTO, int, error)
	GetTopSongs(ctx context.Context, f domain.SpotifyFilters) ([]domain.SongRankingDTO, int, error)
	GetTopAlbums(ctx context.Context, f domain.SpotifyFilters) ([]domain.AlbumRankingDTO, int, error)
	GetHabitsByTimeOfDay(ctx context.Context) ([]domain.HabitTimeDTO, error)
	GetHabitsByDayOfWeek(ctx context.Context) ([]domain.HabitTimeDTO, error)
	GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error)
	GetHistoryEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error)
}

type spotifyRepo struct {
	db *pgxpool.Pool
}

func NewSpotifyRepository(db *pgxpool.Pool) SpotifyRepository {
	return &spotifyRepo{db: db}
}

func buildWhereClause(f domain.SpotifyFilters) (string, []interface{}) {
	// Nota: Ahora filtramos por t.spotify_uri y h.ms_played
	clauses := []string{"h.spotify_uri LIKE 'spotify:track:%'", "h.ms_played > 10000"}
	args := []interface{}{}
	placeholder := 1

	if f.StartDate != nil {
		clauses = append(clauses, fmt.Sprintf("h.played_at >= $%d", placeholder))
		args = append(args, *f.StartDate)
		placeholder++
	}
	if f.EndDate != nil {
		clauses = append(clauses, fmt.Sprintf("h.played_at <= $%d", placeholder))
		args = append(args, *f.EndDate)
		placeholder++
	}
	if f.Album != "" {
		clauses = append(clauses, fmt.Sprintf("t.album_name ILIKE $%d", placeholder))
		args = append(args, "%"+f.Album+"%")
		placeholder++
	}
	if f.Artist != "" {
		clauses = append(clauses, fmt.Sprintf("t.artist_name ILIKE $%d", placeholder))
		args = append(args, "%"+f.Artist+"%")
		placeholder++
	}
	if f.Track != "" {
		clauses = append(clauses, fmt.Sprintf("t.track_name ILIKE $%d", placeholder))
		args = append(args, "%"+f.Track+"%")
		placeholder++
	}
	if f.StartHour != nil && f.EndHour != nil {
		clauses = append(clauses, fmt.Sprintf("EXTRACT(HOUR FROM h.played_at) BETWEEN $%d AND $%d", placeholder, placeholder+1))
		args = append(args, *f.StartHour, *f.EndHour)
		placeholder += 2
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

// GetTotalStats obtiene horas totales y diversidad musical
func (r *spotifyRepo) GetTotalStats(ctx context.Context) (domain.TotalStatsDTO, error) {
	query := `
        SELECT 
            COALESCE(ROUND(SUM(h.ms_played) / 3600000.0, 2), 0) as total_hours,
            COALESCE(ROUND(SUM(h.ms_played) / 60000.0, 2), 0) as total_minutes,
            COALESCE(ROUND(SUM(h.ms_played) / NULLIF(COUNT(DISTINCT h.played_at::date), 0) / 3600000.0, 2), 0) AS average_daily_hours,
            COUNT(DISTINCT t.artist_name) as unique_artists,
            COUNT(DISTINCT h.spotify_uri) as unique_songs
        FROM history h
        JOIN tracks t ON h.spotify_uri = t.spotify_uri
	`

	var stats domain.TotalStatsDTO
	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalHours,
		&stats.TotalMinutes,
		&stats.AverageDailyHours,
		&stats.UniqueArtists,
		&stats.UniqueSongs,
	)
	return stats, err
}

func (r *spotifyRepo) countRows(ctx context.Context, tableQuery string, args []interface{}) (int, error) {
	var total int
	err := r.db.QueryRow(ctx, tableQuery, args...).Scan(&total)
	return total, err
}

// GetTopArtists obtiene el ranking de artistas
func (r *spotifyRepo) GetTopArtists(ctx context.Context, f domain.SpotifyFilters) ([]domain.ArtistRankingDTO, int, error) {
	baseWhere, baseArgs := buildBaseFilters(f)
	searchWhere, searchArgs := buildSearchFilters(f, len(baseArgs)+1)
	allArgs := append(baseArgs, searchArgs...)

	// Obtener el total de registros únicos para la paginación
	countQuery := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT t.artist_name, t.album_name, t.track_name
            FROM history h
            JOIN tracks t ON h.spotify_uri = t.spotify_uri
            %s
            GROUP BY t.artist_name, t.album_name, t.track_name
        )
        SELECT COUNT(DISTINCT artist_name) FROM ranking_completo %s`, baseWhere, searchWhere)

	total, _ := r.countRows(ctx, countQuery, allArgs)

	// Query principal con RANK, LIMIT y OFFSET
	query := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT 
                RANK() OVER (ORDER BY SUM(h.ms_played) DESC) AS ranking,
                t.artist_name,
				STRING_AGG(DISTINCT t.album_name, ' ') as album_name, 
    			STRING_AGG(DISTINCT t.track_name, ' ') as track_name,
                COALESCE(ROUND(SUM(h.ms_played) / 60000.0, 2), 0) as minutes_played,
                COUNT(*) as times_played,
                MAX(a.image_url) as artist_image
            FROM history h
            JOIN tracks t ON h.spotify_uri = t.spotify_uri
            LEFT JOIN artists a ON t.artist_name = a.artist_name
            %s
            GROUP BY t.artist_name
        )
        SELECT ranking, artist_name, minutes_played, times_played, artist_image 
        FROM ranking_completo
        %s
        ORDER BY ranking ASC
        LIMIT $%d OFFSET $%d`,
		baseWhere, searchWhere, len(allArgs)+1, len(allArgs)+2)

	finalArgs := append(allArgs, f.Limit, f.Offset())
	rows, err := r.db.Query(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rankings []domain.ArtistRankingDTO
	for rows.Next() {
		var dto domain.ArtistRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.ArtistName, &dto.MinutesPlayed, &dto.TimesPlayed, &dto.ArtistImage); err != nil {
			return nil, 0, err
		}
		rankings = append(rankings, dto)
	}
	// Si no hay resultados, retornamos un slice vacío (no nil)
	if rankings == nil {
		rankings = []domain.ArtistRankingDTO{}
	}

	return rankings, total, nil
}

// GetTopAlbums obtiene el ranking de álbumes
// Util para wrappeds segun anio, mes, y estaciones del anio (capa service) LIMIT 100
func (r *spotifyRepo) GetTopSongs(ctx context.Context, f domain.SpotifyFilters) ([]domain.SongRankingDTO, int, error) {
	// Separar filtros (Tiempo vs Búsqueda)
	baseWhere, baseArgs := buildBaseFilters(f)
	searchWhere, searchArgs := buildSearchFilters(f, len(baseArgs)+1)
	allArgs := append(baseArgs, searchArgs...)

	// Contamos cuántas canciones (track + artist) coinciden con el filtro de búsqueda
	countQuery := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT t.track_name, t.artist_name, t.album_name
            FROM history h
            JOIN tracks t ON h.spotify_uri = t.spotify_uri
            %s
            GROUP BY t.track_name, t.artist_name, t.album_name
        )
        SELECT COUNT(*) FROM ranking_completo %s`, baseWhere, searchWhere)

	total, err := r.countRows(ctx, countQuery, allArgs)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT 
                RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
                t.track_name, 
                t.artist_name,
                t.album_name, 
                COUNT(*) AS times_played,
                MAX(t.album_image_url) as song_image 
            FROM history h
            JOIN tracks t ON h.spotify_uri = t.spotify_uri
            %s
            GROUP BY t.track_name, t.artist_name, t.album_name
        )
        SELECT ranking, track_name, artist_name, times_played, song_image
        FROM ranking_completo
        %s
        ORDER BY ranking ASC
        LIMIT $%d OFFSET $%d`,
		baseWhere, searchWhere, len(allArgs)+1, len(allArgs)+2)

	finalArgs := append(allArgs, f.Limit, f.Offset())
	rows, err := r.db.Query(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rankings []domain.SongRankingDTO
	for rows.Next() {
		var dto domain.SongRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.TrackName, &dto.ArtistName, &dto.TimesPlayed, &dto.SongImage); err != nil {
			return nil, 0, err
		}
		rankings = append(rankings, dto)
	}
	// Si no hay resultados, retornamos un slice vacío (no nil)
	if rankings == nil {
		rankings = []domain.SongRankingDTO{}
	}

	return rankings, total, nil
}

// GetTopAlbums obtiene el ranking de álbumes
func (r *spotifyRepo) GetTopAlbums(ctx context.Context, f domain.SpotifyFilters) ([]domain.AlbumRankingDTO, int, error) {
	// Obtener cláusulas separadas
	baseWhere, baseArgs := buildBaseFilters(f)
	searchWhere, searchArgs := buildSearchFilters(f, len(baseArgs)+1)
	allArgs := append(baseArgs, searchArgs...)

	countQuery := fmt.Sprintf(`
        WITH ranking_completo AS (
            SELECT t.album_name, t.artist_name, t.track_name
            FROM history h
            JOIN tracks t ON h.spotify_uri = t.spotify_uri
            %s
            GROUP BY t.album_name, t.artist_name, t.track_name
        )
        SELECT COUNT(DISTINCT (album_name, artist_name)) FROM ranking_completo %s`, baseWhere, searchWhere)

	total, err := r.countRows(ctx, countQuery, allArgs)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		WITH ranking_completo AS (
			SELECT 
				RANK() OVER (ORDER BY COUNT(*) DESC) AS ranking,
				t.album_name, 
				t.artist_name,
				STRING_AGG(t.track_name, ' ') as track_name, 
				COUNT(*) AS times_played,
				MAX(t.album_image_url) as album_image
			FROM history h
			JOIN tracks t ON h.spotify_uri = t.spotify_uri
			%s
			GROUP BY t.album_name, t.artist_name
		)
		SELECT ranking, album_name, artist_name, times_played, album_image
		FROM ranking_completo
		%s
		ORDER BY ranking ASC
		LIMIT $%d OFFSET $%d`, 
		baseWhere, searchWhere, len(allArgs)+1, len(allArgs)+2)

	finalArgs := append(allArgs, f.Limit, f.Offset())
	rows, err := r.db.Query(ctx, query, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rankings []domain.AlbumRankingDTO
	for rows.Next() {
		var dto domain.AlbumRankingDTO
		if err := rows.Scan(&dto.Ranking, &dto.AlbumName, &dto.ArtistName, &dto.TimesPlayed, &dto.AlbumImage); err != nil {
			return nil, 0, err
		}
		rankings = append(rankings, dto)
	}

	if rankings == nil {
		rankings = []domain.AlbumRankingDTO{}
	}

	return rankings, total, nil
}

// Momentos del dia por bloque horario, cantidad de escuchas
func (r *spotifyRepo) GetHabitsByTimeOfDay(ctx context.Context) ([]domain.HabitTimeDTO, error) {
	query := `
        SELECT 
            CASE 
                WHEN EXTRACT(HOUR FROM h.played_at) BETWEEN 6 AND 11 THEN 'Mañana'
                WHEN EXTRACT(HOUR FROM h.played_at) BETWEEN 12 AND 17 THEN 'Tarde'
                WHEN EXTRACT(HOUR FROM h.played_at) BETWEEN 18 AND 23 THEN 'Noche'
                ELSE 'Madrugada'
            END AS label,
            COUNT(*) AS count
        FROM history h
        JOIN tracks t ON h.spotify_uri = t.spotify_uri
        GROUP BY label 
        ORDER BY count DESC`

	rows, err := r.db.Query(ctx, query)
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

	// Si la consulta no devuelve nada, retornamos un slice vacío en lugar de nil
	if res == nil {
		return []domain.HabitTimeDTO{}, nil
	}

	return res, nil
}

// Escuchas segun dia de la semana (ingles)
func (r *spotifyRepo) GetHabitsByDayOfWeek(ctx context.Context) ([]domain.HabitTimeDTO, error) {
	query := `
        SELECT 
            EXTRACT(DOW FROM h.played_at) AS num_day, 
            COUNT(*) AS count
        FROM history h
        JOIN tracks t ON h.spotify_uri = t.spotify_uri
        GROUP BY num_day
        ORDER BY num_day`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.HabitTimeDTO
	for rows.Next() {
		var d domain.HabitTimeDTO
		var dayVal int
		if err := rows.Scan(&dayVal, &d.Count); err != nil {
			return nil, err
		}
		d.NumDay = &dayVal
		res = append(res, d)
	}

	if res == nil {
		res = []domain.HabitTimeDTO{}
	}
	return res, nil
}

// Comparativa anual (Tu año en música)
func (r *spotifyRepo) GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
        SELECT 
            EXTRACT(YEAR FROM h.played_at)::int AS year,
            COALESCE(ROUND(SUM(h.ms_played) / 3600000.0, 2), 0) AS total_hours,
            COALESCE(ROUND(SUM(h.ms_played) / 60000.0, 2), 0) AS total_minutes,
            COUNT(*) AS total_songs
        FROM history h
        JOIN tracks t ON h.spotify_uri = t.spotify_uri
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

	if res == nil {
		res = []domain.YearlyStatsDTO{}
	}
	return res, nil
}

// Evolucion historica mensual (Grafico lineas)
func (r *spotifyRepo) GetHistoryEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error) {
	where, args := buildWhereClause(f)
	query := fmt.Sprintf(`
        SELECT
            TO_CHAR(h.played_at, 'YYYY') AS year,
            TO_CHAR(h.played_at, 'MM') AS month,
            TO_CHAR(h.played_at, 'YYYY-MM') AS year_month,
            COALESCE(ROUND(SUM(h.ms_played) / 3600000.0, 2), 0) AS hours_monthly,
            COALESCE(ROUND(SUM(h.ms_played) / 60000.0, 2), 0) AS minutes_monthly
        FROM history h
        JOIN tracks t ON h.spotify_uri = t.spotify_uri
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

	if resul == nil {
		resul = []domain.HistoryEvolutionDTO{}
	}
	return resul, nil
}


func buildWhereArtistTrackClause(f domain.ArtistTrackFilters, startPlaceholder int) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	p := startPlaceholder

	// WHERE actúa sobre las columnas resultantes de la CTE 'ranking_completo'
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


func buildSearchFilters(f domain.SpotifyFilters, startPlaceholder int) (string, []interface{}) {
    if f.Search == "" {
        return "", nil
    }

    clause := fmt.Sprintf("(artist_name ILIKE $%d OR track_name ILIKE $%d OR album_name ILIKE $%d)", 
        startPlaceholder, startPlaceholder, startPlaceholder)
    
    return "WHERE " + clause, []interface{}{"%" + f.Search + "%"}
}

func buildBaseFilters(f domain.SpotifyFilters) (string, []interface{}) {
	// Filtros que SI afectan el cálculo del ranking (Fechas y Horas)
	clauses := []string{"h.spotify_uri LIKE 'spotify:track:%'", "h.ms_played > 10000"}
	args := []interface{}{}
	p := 1

	if f.StartDate != nil {
		clauses = append(clauses, fmt.Sprintf("h.played_at >= $%d", p))
		args = append(args, *f.StartDate)
		p++
	}
	if f.EndDate != nil {
		clauses = append(clauses, fmt.Sprintf("h.played_at <= $%d", p))
		args = append(args, *f.EndDate)
		p++
	}
	if f.StartHour != nil && f.EndHour != nil {
		clauses = append(clauses, fmt.Sprintf("EXTRACT(HOUR FROM h.played_at) BETWEEN $%d AND $%d", p, p+1))
		args = append(args, *f.StartHour, *f.EndHour)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}
