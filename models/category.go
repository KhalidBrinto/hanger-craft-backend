package models

import (
	"errors"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name         null.String `gorm:"size:100;not null"`
	CategoryType null.String `gorm:"size:100;not null;check:category_type IN ('parent', 'child')"`
	ParentID     *uint
	Products     []Product `gorm:"foreignKey:CategoryID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {

	if c.CategoryType.String == "child" && c.ParentID == nil {
		return errors.New("must provide parent category id when creating child category")
	}

	return nil

}
