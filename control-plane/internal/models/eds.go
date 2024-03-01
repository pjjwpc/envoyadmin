package models

type Eds struct {
	Id        int64  `gorm:"id" json:"id"`
	ValueData string `gorm:"value_data" json:"valueData"`
	Version   string `gorm:"version" json:"version"`
}

func (Eds) TableName() string {
	return "eds"
}
