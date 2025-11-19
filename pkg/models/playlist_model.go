package models

type Playlist struct {
	BaseModel
	Name            string `json:"name" form:"name"`
	Image           string `json:"image" form:"image"`
	Private         bool   `json:"private" form:"private"`
	Description     string `json:"description" form:"description"`
	IsCollaborative bool   `json:"is_collaborative" form:"is_collaborative"`

	CreatorUserID *string `json:"creator_user_id,omitempty" gorm:"index"`
	CreatorUser   *User   `gorm:"foreignKey:CreatorUserID"`

	CreatorArtistID *string `json:"creator_artist_id,omitempty" gorm:"index"`
	CreatorArtist   *Artist `gorm:"foreignKey:CreatorArtistID"`

	PlaylistSongs []PlaylistSong `gorm:"foreignKey:PlaylistID"`
}
