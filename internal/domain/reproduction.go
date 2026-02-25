package domain

type SpotifyRecord struct {
	TS          string `json:"ts"`
	Platform    string `json:"platform"`
	MsPlayed    int    `json:"ms_played"`
	ConnCountry string `json:"conn_country"`
	TrackName   string `json:"master_metadata_track_name"`
	ArtistName  string `json:"master_metadata_album_artist_name"`
	AlbumName   string `json:"master_metadata_album_album_name"`
	SpotifyURI  string `json:"spotify_track_uri"`
}

type InputFilter struct {
	ArtistName string
}
