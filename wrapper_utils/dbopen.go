package wrapper_utils

import (
	"github.com/jinzhu/gorm"
	"log"
)

var (
	DB *gorm.DB
)

func DBopen(driver string) *gorm.DB {
	connectionstring := Sqlenv("CONN_STR", "")
	_db, dberr := gorm.Open(driver, connectionstring)
	if dberr != nil {
		log.Println(dberr)
		return nil
	}
	_db.SingularTable(true)
	return _db
}
