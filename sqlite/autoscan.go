package sqlite

import (
	"errors"
	"reflect"

	"github.com/harrybrwn/db"
)

func autoscan[T any](row db.Scanner, dst T) error {
	v := reflect.ValueOf(dst)
	tp := reflect.TypeOf(dst)
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
		v = v.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return errors.New("cannot autoscan non-structs for now")
	}
	fields := make([]any, tp.NumField())
	for i := range tp.NumField() {
		fieldVal := v.Field(i)
		fields[i] = fieldVal.Addr().Interface()
	}
	return row.Scan(fields...)
}

type scanner interface {
	db.Scanner
	Next() bool
}

func autoscanRows[T any](rows scanner) ([]T, error) {
	res := make([]T, 0)
	for rows.Next() {
		var v T
		err := autoscan(rows, &v)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func scanStrings(rows scanner) ([]string, error) {
	res := make([]string, 0)
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}
