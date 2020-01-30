package wrapper_utils

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func Migration(driver string, dirname string) {
	files, _ := filepath.Glob(filepath.Join(dirname, "*.sql"))
	DB = DBopen(driver)
	if DB == nil {
		for _ = 0; ; {
			time.Sleep(time.Second * 2)
			DB = DBopen(driver)
			if DB != nil {
				break
			}
		}
	}

	err := DB.DB().Ping()
	if err != nil {
		for _ = 0; ; {
			time.Sleep(2)
			err = DB.DB().Ping()
			if err == nil {
				break
			}
		}
	}
	DB.DB().SetMaxIdleConns(50)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Second * 30)

	err = DB.Transaction(func(tx *gorm.DB) error {
		for _, file := range files {
			byte_sql, _ := ioutil.ReadFile(file)
			sqls := strings.Split(string(byte_sql), ";")
			for _, sql := range sqls {
				if sql != "" {
					err := tx.Exec(sql).Error
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}
