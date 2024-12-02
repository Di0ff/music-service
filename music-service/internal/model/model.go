package model

type Song struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Group       string `gorm:"size:20;not null" json:"group"`
	SongName    string `gorm:"size:20;not null" json:"song"`
	ReleaseDate string `gorm:"size:10;not null" json:"release_date"`
	Text        string `gorm:"type:text" json:"text"`
	Link        string `gorm:"size:40" json:"link"`
}
