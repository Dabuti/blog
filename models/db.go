package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewDB opens a db connection
func NewDB(dataSourceName string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}
	return db, nil
}
