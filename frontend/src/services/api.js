import axios from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1/spotify',
    headers: {
        'Content-Type': 'application/json',
    },
});

export default {
    // 1. Top Artists
    getTopArtists(params) {
        return api.get('/top/artists', { params });
    },

    // 2. Top Songs
    getTopSongs(params) {
        return api.get('/top/songs', { params });
    },

    // 3. Top Albums
    getTopAlbums(params) {
        return api.get('/top/albums', { params });
    },

    // 4. General Stats
    getStats(params) {
        return api.get('/stats', { params });
    },

    // 5. Habits
    getHabits(type) {
        return api.get('/habits', { params: { type } });
    },

    // 6. Evolution
    getEvolution(params) {
        return api.get('/evolution', { params });
    },

    // 7 & 8. Search Rank
    getSearchRank(params) {
        return api.get('/search-rank', { params });
    },

    // 9. Yearly Snapshot
    getYearly(params) {
        return api.get('/yearly', { params });
    },

    // 10, 11, 12. Wrapped
    getWrapped(params) {
        return api.get('/wrapped', { params });
    }
};
