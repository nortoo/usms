package store

import "gorm.io/gorm"

type DBName string

var dbs = make(map[DBName]*gorm.DB)

const (
	Default DBName = "default"
)
