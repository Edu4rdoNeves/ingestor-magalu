package entity

import (
	"time"

	"gorm.io/gorm"
)

type PulseData struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Tenant     string         `gorm:"column:tenant;index:idx_unique_pulse,unique" json:"tenant"`
	ProductSku string         `gorm:"column:product_sku;index:idx_unique_pulse,unique" json:"product_sku"`
	UseUnity   string         `gorm:"column:use_unity;index:idx_unique_pulse,unique" json:"use_unity"`
	UsedAmount float64        `gorm:"column:used_amount" json:"used_amount"`
	CreatedAt  time.Time      `gorm:"auto_created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (PulseData) TableName() string {
	return "pulses"
}
