package models

type Playlist struct {
	BaseModel
	Name            string `json:"name"`
	Image           string `json:"image"`
	Private         bool   `json:"private"`
	Description     string `json:"description"`
	IsCollaborative bool   `json:"is_collaborative"`

	CreatorUserID   *string `gorm:"index" json:"creator_user_id"`
	CreatorUser     *User   `gorm:"foreignKey:CreatorUserID"`
	CreatorArtistID *string `gorm:"index" json:"creator_artist_id"`
	CreatorArtist   *Artist `gorm:"foreignKey:CreatorArtistID"`

	PlaylistSongs []PlaylistSong `gorm:"foreignKey:PlaylistID" json:"playlist_songs"`
}
