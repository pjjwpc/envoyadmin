package models

type Rls struct {
	Id        int64  `gorm:"id" json:"id"`
	ValueData string `gorm:"value_data" json:"valueData"`
	Version   string `gorm:"version" json:"version"`
}

func (Rls) TableName() string {
	return "rls"
}
