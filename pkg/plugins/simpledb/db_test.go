package simpledb

import (
	"fmt"
	"testing"
)

func TestCreateInternalDb(t *testing.T) {
	db, err := CreateInternalDb("./demo.db")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFluffDb_InsertUsers(t *testing.T) {
	db, err := OpenInternalDB("./demo.db")
	if err != nil {
		t.Fatal(err)
	}
	err = db.InsertUsers("admin", "demo1233", "admin")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFluffDb_GetUserByName(t *testing.T) {
	db, err := OpenInternalDB("./demo.db")
	if err != nil {
		t.Fatal(err)
	}
	user, err := db.GetUserByName("admin")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
}
