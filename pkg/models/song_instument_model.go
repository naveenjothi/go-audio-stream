package models

type SongInstrument struct {
	BaseModel
	Name string `json:"name"`

	Songs []Song `gorm:"many2many:song_instrument_map;" json:"songs"`
}
