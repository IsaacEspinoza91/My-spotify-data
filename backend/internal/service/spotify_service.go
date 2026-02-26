package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/domain"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/repository"
)

type SpotifyService interface {
	GetDashboardStats(ctx context.Context, f domain.SpotifyFilters) (domain.TotalStatsDTO, error)
	GetTopList(ctx context.Context, listType string, limit int, f domain.SpotifyFilters) (interface{}, error)
	GetHabitAnalysis(ctx context.Context, habitType string, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error)
	GetGlobalEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error)
	SearchRankedItem(ctx context.Context, f domain.SpotifyFilters, target domain.ArtistTrackFilters, limit int) (interface{}, error)
	GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error)
	GetYearlyWrapped(ctx context.Context, year int) (interface{}, error)
	GetMonthlyWrapped(ctx context.Context, year, month int) (interface{}, error)
	GetSeasonalWrapped(ctx context.Context, year int, season domain.Season) (interface{}, error)
}

type spotifyService struct {
	repo repository.SpotifyRepository
}

func NewSpotifyService(repo repository.SpotifyRepository) SpotifyService {
	return &spotifyService{repo: repo}
}

// Implementación de SpotifyService

func (s *spotifyService) GetDashboardStats(ctx context.Context, f domain.SpotifyFilters) (domain.TotalStatsDTO, error) {
	f.CleanAndValidate()
	return s.repo.GetTotalStats(ctx, f)
}

func (s *spotifyService) GetTopList(ctx context.Context, listType string, limit int, f domain.SpotifyFilters) (interface{}, error) {
	f.CleanAndValidate()
	if limit <= 0 {
		limit = 10
	}

	switch listType {
	case "artists":
		return s.repo.GetTopArtists(ctx, limit, f)
	case "songs":
		return s.repo.GetTopSongs(ctx, limit, f)
	case "albums":
		return s.repo.GetTopAlbums(ctx, limit, f)
	default:
		return nil, nil
	}
}

func (s *spotifyService) GetHabitAnalysis(ctx context.Context, habitType string, f domain.SpotifyFilters) ([]domain.HabitTimeDTO, error) {
	f.CleanAndValidate()
	if habitType == "dow" { // Day of Week. Domingo = 0
		return s.repo.GetHabitsByDayOfWeek(ctx, f)
	}
	// Tiempo del dia, tarde, noche, etc
	return s.repo.GetHabitsByTimeOfDay(ctx, f)
}

func (s *spotifyService) GetYearlyStats(ctx context.Context, f domain.SpotifyFilters) ([]domain.YearlyStatsDTO, error) {
	f.CleanAndValidate()
	return s.repo.GetYearlyStats(ctx, f)
}

func (s *spotifyService) GetGlobalEvolution(ctx context.Context, f domain.SpotifyFilters) ([]domain.HistoryEvolutionDTO, error) {
	f.CleanAndValidate()
	return s.repo.GetHistoryEvolution(ctx, f)
}

// SearchRankedItem permite buscar dónde quedó un artista o canción específica en el ranking global
func (s *spotifyService) SearchRankedItem(ctx context.Context, f domain.SpotifyFilters, target domain.ArtistTrackFilters, limit int) (interface{}, error) {
	f.CleanAndValidate()
	target.Clean()

	if limit <= 0 {
		limit = 10
	}

	if target.Track != "" {
		return s.repo.GetRankedSongs(ctx, f, target, limit)
	}
	return s.repo.GetRankedArtist(ctx, f, target, limit)
}

// Metodos para obtener wrappeds segun el año, mes o estacion
func (s *spotifyService) GetYearlyWrapped(ctx context.Context, year int) (interface{}, error) {
	loc, _ := time.LoadLocation("America/Santiago")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(1, 0, 0).Add(-time.Second)

	f := domain.SpotifyFilters{StartDate: &start, EndDate: &end}
	return s.repo.GetTopSongs(ctx, 100, f)
}

func (s *spotifyService) GetMonthlyWrapped(ctx context.Context, year, month int) (interface{}, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("el mes %d no es válido (debe ser 1-12)", month)
	}

	loc, _ := time.LoadLocation("America/Santiago")
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	f := domain.SpotifyFilters{StartDate: &start, EndDate: &end}
	return s.repo.GetTopSongs(ctx, 100, f)
}

func (s *spotifyService) GetSeasonalWrapped(ctx context.Context, year int, season domain.Season) (interface{}, error) {
	validSeasons := map[domain.Season]bool{
		domain.Summer: true, domain.Autumn: true,
		domain.Winter: true, domain.Spring: true,
	}
	if !validSeasons[season] {
		return nil, fmt.Errorf("estación '%s' no válida. Use: summer, autumn, winter o spring", season)
	}

	loc, _ := time.LoadLocation("America/Santiago")
	var start, end time.Time

	switch season {
	case domain.Summer: // 21 Dic (año anterior) - 20 Mar
		start = time.Date(year-1, 12, 21, 0, 0, 0, 0, loc)
		end = time.Date(year, 3, 20, 23, 59, 59, 0, loc)
	case domain.Autumn: // 21 Mar - 20 Jun
		start = time.Date(year, 3, 21, 0, 0, 0, 0, loc)
		end = time.Date(year, 6, 20, 23, 59, 59, 0, loc)
	case domain.Winter: // 21 Jun - 20 Sep
		start = time.Date(year, 6, 21, 0, 0, 0, 0, loc)
		end = time.Date(year, 9, 20, 23, 59, 59, 0, loc)
	case domain.Spring: // 21 Sep - 20 Dic
		start = time.Date(year, 9, 21, 0, 0, 0, 0, loc)
		end = time.Date(year, 12, 20, 23, 59, 59, 0, loc)
	}

	f := domain.SpotifyFilters{StartDate: &start, EndDate: &end}
	return s.repo.GetTopSongs(ctx, 100, f)
}
