package models

type Vhds struct {
	Id        int64  `gorm:"id" json:"id"`
	ValueData string `gorm:"value_data" json:"valueData"`
	Version   string `gorm:"version" json:"version"`
}

func (Vhds) TableName() string {
	return "vhds"
}
