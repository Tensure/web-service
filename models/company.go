package models

import (
	"errors"
	"gorm.io/gorm"
)

type Company struct {
	CompanyName string `json:"companyName"`
	User        string `json:"user"`
	UserEmail   string `json:"userEmail"`
	Position    string `json:"position"`
}

type HasIssues struct {
	Company string `json:"company"`
	Issues  bool   `json:"issues"`
}

type Companies []Company

func (c *Company) CreateCompany(db *gorm.DB) (err error) {
	result := db.Create(&c)
	if result.RowsAffected == 0 {
		return errors.New("error creating company")
	}
	return nil
}

func (c *Companies) GetAllCompanies(db *gorm.DB) (err error) {
	// Get all records
	if err := db.Find(&c).Error; err != nil {
		return err
	}
	return nil
}

func (c *Company) GetCompany(db *gorm.DB) (err error) {

	// Get first matched record
	// SELECT * FROM users WHERE company_name = 'c.companyName' ORDER BY company_name LIMIT 1;
	// Get all records
	if err := db.Where("company_name = ?", c.CompanyName).First(&c).Error; err != nil {
		return err
	}
	return nil
}

// Update updates company
func (c *Company) Update(db *gorm.DB) error {
	if err := db.Save(&c).Error; err != nil {
		return err
	}
	return nil
}

//
// Delete will soft delete company, must implement IDs
//
//func (f *Feedback) Delete(db *gorm.DB, id uuid.UUID) error {
//	if err := db.Where("id = ?", id).Delete(&f).Error; err != nil {
//		return err
//	}
//	return nil
//}
