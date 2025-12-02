package models

type Song struct {
	BaseModel
	Name        string `json:"name"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	VideoURL    string `json:"video_url"`
	Duration    int32  `json:"duration"`
	Explicit    bool   `json:"explicit"`
	TrackNumber *int16 `json:"track_number"`
	Language    string `json:"language"`

	Artists       []Artist         `gorm:"many2many:artist_song;" json:"artists"`
	PlaylistSongs []PlaylistSong   `gorm:"foreignKey:SongID" json:"playlist_songs"`
	Features      *SongFeatures    `gorm:"foreignKey:SongID" json:"features"`
	Tags          []SongTag        `gorm:"many2many:song_tag_map;" json:"tags"`
	Instruments   []SongInstrument `gorm:"many2many:song_instrument_map;" json:"instruments"`
}
