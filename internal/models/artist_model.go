package models

type ArtistModel struct {
	BaseModel
	Name      string          `json:"name" form:"name"`
	Image     string          `json:"image" form:"image"`
	Followers int64           `json:"followers"`
	Songs     []SongModel     `gorm:"many2many:artist_song;"`
	Playlists []PlaylistModel `gorm:"foreignKey:CreatorArtistID"`
}
