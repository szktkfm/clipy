package clipy

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type Clipboard struct {
	gorm.Model
	ClipText string
}

func NewRepository(dbpath string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Clipboard{})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) write(b []byte) {
	r.db.Create(&Clipboard{ClipText: string(b)})
}

func (r *Repository) read(offset, limit int) ([]Clipboard, error) {
	cb := []Clipboard{}
	result := r.db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&cb)
	if result.Error != nil {
		return nil, result.Error
	}
	return cb, nil
}

func (r *Repository) count() (int64, error) {
	var count int64
	result := r.db.Model(&Clipboard{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
