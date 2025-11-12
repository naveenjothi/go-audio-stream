package models

type UserModel struct {
	BaseModel
	Email     string `json:"email" form:"email"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Mobile    string `json:"mobile" form:"mobile"`
	Username  string `json:"user_name" form:"user_name"`

	Follows   []ArtistModel   `gorm:"many2many:user_follows_artist;"`
	Playlists []PlaylistModel `gorm:"foreignKey:CreatorUserID"`
}
