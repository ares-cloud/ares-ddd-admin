package database

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = fmt.Errorf("record not found")
)

func IfErrorNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// IsUniqueIndexError ， 判断是否为索引错误
// 参数：
//
//	err ： desc
//
// 返回值：
//
//	bool ：desc
func IsUniqueIndexError(err error) bool {
	errType := reflect.TypeOf(err).String()
	if errType == "*mysql.MySQLError" && err.(*mysql.MySQLError).Number == 1062 {
		return true
	}
	return false
}
func IsForeignKeyError(err error) bool {
	errType := reflect.TypeOf(err).String()
	if errType == "*mysql.MySQLError" && err.(*mysql.MySQLError).Number == 1452 {
		return true
	}
	return false
}
func IsDuplicateEntryError(err error) bool {
	errType := reflect.TypeOf(err).String()
	if errType == "*mysql.MySQLError" && err.(*mysql.MySQLError).Number == 1062 {
		return true
	}
	return false
}
