package models

type PlaylistModel struct {
	BaseModel
	Name    string `json:"name" form:"name"`
	Image   string `json:"image" form:"image"`
	Private bool   `json:"private" form:"private"`

	CreatorUserID *string    `json:"creator_user_id,omitempty" gorm:"index"`
	CreatorUser   *UserModel `gorm:"foreignKey:CreatorUserID"`

	CreatorArtistID *string      `json:"creator_artist_id,omitempty" gorm:"index"`
	CreatorArtist   *ArtistModel `gorm:"foreignKey:CreatorArtistID"`

	Songs []SongModel `gorm:"many2many:playlist_song;"`
}
