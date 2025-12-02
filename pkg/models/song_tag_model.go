package models

type SongTag struct {
	BaseModel
	Name string `json:"name"`

	Songs []Song `gorm:"many2many:song_tag_map;" json:"songs"`
}
