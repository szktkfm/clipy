package clipy

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type Clipboard struct {
	gorm.Model
	Bytes string
}

func NewRepository(dbpath string) *Repository {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Clipboard{})

	return &Repository{db: db}
}

// Database操作
func (r *Repository) write(b []byte) {
	r.db.Create(&Clipboard{Bytes: string(b)})
}

func (r *Repository) read(offset, limit int) []Clipboard {
	cb := []Clipboard{}
	result := r.db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&cb)
	if result.Error != nil {
		log.Fatal(result)
	}
	return cb
}

func (r *Repository) count() int64 {
	var count int64
	r.db.Model(&Clipboard{}).Count(&count)
	return count
}
