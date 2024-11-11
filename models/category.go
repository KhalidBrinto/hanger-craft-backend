package models

import (
	"backend/utils"
	"errors"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name         null.String `gorm:"size:100;not null"`
	CategoryType null.String `gorm:"size:100;not null;check:category_type IN ('parent', 'child', 'grandchild')"`
	ParentID     *uint
	Image        *CategoryImage `gorm:"foreignKey:CategoryID"`
	Products     []Product      `gorm:"foreignKey:CategoryID"`
}

type CategoryImage struct {
	ID         uint `gorm:"primaryKey"`
	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID" json:"-"`
	Image      string   `gorm:"-" json:"Image"`
	ImageBytes []byte   `gorm:"column:image;type:bytea" json:"-"` // Binary data for the image
}

func (c *CategoryImage) AfterFind(tx *gorm.DB) (err error) {
	c.Image = utils.EncodeImageToBase64(c.ImageBytes)

	return nil

}

func (c *CategoryImage) BeforeCreate(tx *gorm.DB) (err error) {
	bt, err := utils.DecodeBase64Image(c.Image)
	if err != nil {
		return err
	}

	c.ImageBytes = bt

	return nil

}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {

	if (c.CategoryType.String == "child" || c.CategoryType.String == "grandchild") && c.ParentID == nil {
		return errors.New("must provide parent category id when creating sub category")
	}

	return nil

}
