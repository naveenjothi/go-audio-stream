package models

type Artist struct {
	BaseModel
	Name      string     `json:"name" form:"name"`
	Image     string     `json:"image" form:"image"`
	Followers int64      `json:"followers"`
	Songs     []Song     `gorm:"many2many:artist_song;"`
	Playlists []Playlist `gorm:"foreignKey:CreatorArtistID"`
	Verified  bool       `json:"verified"`
	Bio       string     `json:"bio" form:"bio"`
}
