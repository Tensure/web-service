package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Feedback struct {
	Company   string    `json:"company"`
	Happiness int       `json:"happiness"`
	Timestamp time.Time `json:"timestamp"`
}

func (f *Feedback) CreateFeedback(db *gorm.DB) (err error) {
	result := db.Create(&f)
	if result.RowsAffected == 0 {
		return errors.New("feedback not created")
	}
	return nil
}

func (f *Feedback) GetAllFeedback(db *gorm.DB) (err error) {
	// Get all records
	if err := db.Find(&f).Error; err != nil {
		return err
	}
	return nil
}

// Update updates feedback
func (f *Feedback) Update(db *gorm.DB) error {
	if err := db.Save(&f).Error; err != nil {
		return err
	}
	return nil
}

//
// Delete will soft delete feedback, must implement IDs
//
//func (f *Feedback) Delete(db *gorm.DB, id uuid.UUID) error {
//	if err := db.Where("id = ?", id).Delete(&f).Error; err != nil {
//		return err
//	}
//	return nil
//}
