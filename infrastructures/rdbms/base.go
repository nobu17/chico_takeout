package rdbms

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
   ID        string `gorm:"primary_key"`
   CreatedAt time.Time
   UpdatedAt time.Time
   DeletedAt gorm.DeletedAt `gorm:"index"`
  }