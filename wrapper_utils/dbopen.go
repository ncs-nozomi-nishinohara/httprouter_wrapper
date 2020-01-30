package wrapper_utils

import (
	"github.com/jinzhu/gorm"
	"github.com/prometheus/common/log"
)

var (
	DB *gorm.DB
)

func DBopen(driver string) *gorm.DB {
	connectionstring := Migration_env("CONN_STR", "")
	_db, dberr := gorm.Open(driver, connectionstring)
	if dberr != nil {
		log.Errorln(dberr)
		return nil
	}
	_db.SingularTable(true)
	return _db
}
