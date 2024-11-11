package models

import (
	"backend/utils"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	Name      null.String `gorm:"size:100;not null"`
	Status    *string     `gorm:"not null;check:status IN ('published', 'unpublished')"`
	Permalink *string     `gorm:"not null;default:''"`
	Logo      *BrandImage `gorm:"foreignKey:BrandID"`
	Products  []Product   `gorm:"foreignKey:BrandID"`
}

type BrandImage struct {
	ID         uint `gorm:"primaryKey"`
	BrandID    uint
	Brand      Brand  `gorm:"foreignKey:BrandID" json:"-"`
	Image      string `gorm:"-" json:"Image"`
	ImageBytes []byte `gorm:"column:image;type:bytea" json:"-"` // Binary data for the image
}

func (c *BrandImage) AfterFind(tx *gorm.DB) (err error) {
	c.Image = utils.EncodeImageToBase64(c.ImageBytes)

	return nil

}

func (c *BrandImage) BeforeCreate(tx *gorm.DB) (err error) {
	bt, err := utils.DecodeBase64Image(c.Image)
	if err != nil {
		return err
	}

	c.ImageBytes = bt

	return nil

}

// func (c *Brand) BeforeCreate(tx *gorm.DB) (err error) {

// 	if c.CategoryType.String == "child" && c.ParentID == nil {
// 		return errors.New("must provide parent category id when creating child category")
// 	}

// 	return nil

// }
