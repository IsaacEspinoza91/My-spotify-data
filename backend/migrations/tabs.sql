CREATE TABLE spotify_history (
    id SERIAL PRIMARY KEY,
    ts TIMESTAMP,
    platform TEXT,
    ms_played INTEGER,
    conn_country VARCHAR(5),
    track_name TEXT,
    artist_name TEXT,
    album_name TEXT,
    spotify_uri TEXT
);

-- 1. Índice para el rango de fechas y análisis temporal (Crucial para el Wrapped)
CREATE INDEX idx_spotify_ts ON spotify_history (ts);

-- 2. Índice para filtrar podcasts y agrupar por canción de forma única
CREATE INDEX idx_spotify_uri ON spotify_history (spotify_uri) 
WHERE (spotify_uri LIKE 'spotify:track:%');

-- 3. Índice para acelerar los rankings de artistas y álbumes
CREATE INDEX idx_spotify_artist_album ON spotify_history (artist_name, album_name);