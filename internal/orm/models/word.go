package models

type Word struct {
	BaseModel
	Value  string `gorm:"NOT NULL" json:"value"`
	Tags []*Tag `gorm:"many2many:tag_word"`
}

