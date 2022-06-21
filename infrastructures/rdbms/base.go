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

type BaseRepository struct {
	Db *gorm.DB
}

func (b *BaseRepository) Transact(fc func() error) (error) {
   var err error = nil
   b.Db.Transaction(func(tx *gorm.DB) error {
      err = fc()
      return err
   })
   return err
}