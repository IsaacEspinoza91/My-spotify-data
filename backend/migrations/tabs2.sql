CREATE TABLE artists (
    spotify_artist_uri TEXT PRIMARY KEY, -- Se llenará con la API
    artist_name TEXT UNIQUE NOT NULL,    -- Usaremos esto para el primer enlace
    image_url TEXT
);

-- 2. Tabla de Tracks (Relacionada con Artistas)
CREATE TABLE tracks (
    spotify_uri TEXT PRIMARY KEY,
    track_name TEXT NOT NULL,
    album_name TEXT,
    album_image_url TEXT,
    artist_name TEXT REFERENCES artists(artist_name) -- Enlace inicial por nombre
);

-- 3. Tabla de Historial
CREATE TABLE history (
    id SERIAL PRIMARY KEY,
    spotify_uri TEXT REFERENCES tracks(spotify_uri),
    played_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ms_played INTEGER
);



-- Para acelerar los filtros por fecha (StartDate, EndDate) y ordenamiento cronológico
CREATE INDEX idx_history_played_at ON history (played_at);

-- Crucial para los JOINs frecuentes con la tabla tracks
CREATE INDEX idx_history_spotify_uri ON history (spotify_uri);

-- Índice compuesto para optimizar GetHabitsByTimeOfDay (EXTRACT HOUR)
-- Solo si tienes millones de filas y esa consulta es lenta
--  ERROR INMUTABLE  CREATE INDEX idx_history_hour_played ON history (EXTRACT(HOUR FROM played_at));


-- Para acelerar GetTopSongs y GetRankedSongs (agrupación por nombre y artista)
CREATE INDEX idx_tracks_name_artist ON tracks (track_name, artist_name);

-- Para búsquedas rápidas por álbum (GetTopAlbums)
CREATE INDEX idx_tracks_album_artist ON tracks (album_name, artist_name);

-- Optimización para búsquedas de texto parcial (ILIKE '%search%')
-- Nota: Requiere la extensión pg_trgm en PostgreSQL
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE INDEX idx_tracks_track_name_trgm ON tracks USING gin (track_name gin_trgm_ops);
-- CREATE INDEX idx_tracks_artist_name_trgm ON tracks USING gin (artist_name gin_trgm_ops);

-- El PRIMARY KEY ya crea un índice sobre spotify_artist_uri.
-- Como usas artist_name para el JOIN desde tracks, este índice es obligatorio:
CREATE INDEX idx_artists_name ON artists (artist_name);

-- Para búsquedas de texto en el nombre del artista
-- CREATE INDEX idx_artists_name_trgm ON artists USING gin (artist_name gin_trgm_ops);