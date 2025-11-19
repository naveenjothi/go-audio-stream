package models

type Song struct {
	BaseModel
	Name        string `json:"name" form:"name"`
	Image       string `json:"image" form:"image"`
	Url         string `json:"url"`
	VideoUrl    string `json:"video_url"`
	Duration    int32  `json:"duration"`
	Explicit    bool   `json:"explicit"` //#parental control
	TrackNumber *int16 `json:"track_number,omitempty"`
	Language    string `json:"language"`

	Artists       []Artist       `gorm:"many2many:artist_song;"`
	PlaylistSongs []PlaylistSong `gorm:"foreignKey:SongID"`
}
