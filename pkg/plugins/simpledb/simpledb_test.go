package simpledb

import (
	"fmt"
	"testing"
)

func TestSimpleDb_CreateTable(t *testing.T) {
	db := NewSimpleDB()
	if err := db.CreateTable(); err != nil {
		fmt.Println(err)
	}
	if err := db.Insert("test", "test123", "testing", "demo"); err != nil {
		fmt.Println(err)
	}
	if err := db.Insert("test1", "", "", "demo123"); err != nil {
		fmt.Println(err)
	}

	db = NewSimpleDB()
	value, err := db.Get("test1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value.KeyValue, value.PrivateKeyValue, value.PublicKeyValue, value.TokenValue)
	db.CloseDb()
}

func TestSimpleDb_Delete(t *testing.T) {
	db := NewSimpleDB()
	if err := db.Delete("test1"); err != nil {
		fmt.Println(err)
	}
	db.CloseDb()
}
