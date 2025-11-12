package models

type SongModel struct {
	BaseModel
	Name     string `json:"name" form:"name"`
	Image    string `json:"image" form:"image"`
	Url      string `json:"url"`
	VideoUrl string `json:"video_url"`

	Artists   []ArtistModel   `gorm:"many2many:artist_song;"`
	Playlists []PlaylistModel `gorm:"many2many:playlist_song;"`
}
