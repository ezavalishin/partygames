package models

type Tag struct {
	BaseModel
	Value  string `gorm:"NOT NULL" json:"value"`
	Words []*Word `gorm:"many2many:tag_word"`
}

