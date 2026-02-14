package model

type Country struct {
	Code string `gorm:"primaryKey;type:char(2)"` // ISO 3166-1 (ex: FR, US)
	Name string `gorm:"type:varchar(100);not null"`
}
