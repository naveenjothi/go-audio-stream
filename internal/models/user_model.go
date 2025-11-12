package models

type User struct {
	BaseModel
	Email     string `json:"email" form:"email" gorm:"uniqueIndex"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Mobile    string `json:"mobile" form:"mobile" gorm:"uniqueIndex"`
	Username  string `json:"user_name" form:"user_name" gorm:"uniqueIndex"`

	Follows   []Artist   `json:"follows,omitempty" gorm:"many2many:user_follows_artist;"`
	Playlists []Playlist `json:"playlists,omitempty" gorm:"foreignKey:CreatorUserID"`
}
