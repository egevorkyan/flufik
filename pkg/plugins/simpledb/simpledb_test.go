package simpledb

import (
	"fmt"
	"testing"
)

func TestSimpleDb_CreateTable(t *testing.T) {
	db := NewSimpleDB("test.db")
	if err := db.CreateTable("kvdb"); err != nil {
		fmt.Println(err)
	}
	if err := db.CreateTable("app"); err != nil {
		fmt.Println(err)
	}

}

func TestSimpleDb_Insert(t *testing.T) {
	appName := "art"
	appVersion := "4"
	appArch := "x86_64"
	appOsVersion := "linux"
	appLocation := "/opt/art/arm64"
	db := NewSimpleDB("test.db")
	if err := db.Insert("app", appName, appVersion, appArch, appOsVersion, appLocation); err != nil {
		fmt.Println(err)
	}
}

func TestSimpleDb_Get(t *testing.T) {
	db := NewSimpleDB("test.db")
	key := map[string]interface{}{"appName": "art", "appArch": "x86_64", "appOsVersion": "darwin"}
	a, err := db.GetLatestApp(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}

func TestSimpleDb_GetAppByVersion(t *testing.T) {
	db := NewSimpleDB("test.db")
	key := map[string]interface{}{"appName": "art", "appArch": "x86_64", "appOsVersion": "darwin", "appVersion": "4"}
	a, err := db.GetAppByVersion(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}

func TestSimpleDb_GetKey(t *testing.T) {
	db := NewSimpleDB("test.db")
	d, err := db.GetKey("test10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
}

func TestSimpleDb_Delete(t *testing.T) {
	db := NewSimpleDB("test.db")
	if err := db.Delete("test10"); err != nil {
		t.Fail()
	}
}

func TestSimpleDb_DeleteApp(t *testing.T) {
	db := NewSimpleDB("test.db")
	if err := db.DeleteApp("art", "4", "arm64", ""); err != nil {
		t.Fail()
	}
}
