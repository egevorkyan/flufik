package nosql

import (
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
)

type TieDot struct {
	tdb       *db.DB
	logger    *logging.Logger
	debugging string
}

func NewTieDot(collectionName string, dataIndexKey string, logger *logging.Logger, debugging string) (*TieDot, error) {
	// (Create if not exist) open a database
	myDB, err := db.OpenDB(core.FlufikNoSqlDbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open or create db: %v", err)
	}
	if !myDB.ColExists(collectionName) {
		err = myDB.Create(collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection %s: %v", collectionName, err)
		}
		col := myDB.Use(collectionName)
		err = col.Index([]string{dataIndexKey})
		if err != nil {
			return nil, fmt.Errorf("failed to index %v in collection %v: %v", dataIndexKey, collectionName, err)
		}
	}
	return &TieDot{
		tdb:       myDB,
		logger:    logger,
		debugging: debugging,
	}, nil
}

func (t *TieDot) Insert(data map[string]interface{}, collectionName string) error {
	if t.debugging == "1" {
		t.logger.Info("insert to db")
	}
	col := t.tdb.Use(collectionName)

	_, err := col.Insert(data)
	if err != nil {
		return fmt.Errorf("failed to insert data to collection %v: %v", collectionName, err)
	}
	return nil
}

func (t *TieDot) Get(queryRequest interface{}, collectionName string) (docId int, readBack map[string]interface{}, err error) {
	if t.debugging == "1" {
		t.logger.Info("get data from db")
	}
	col := t.tdb.Use(collectionName)
	queryResult := make(map[int]struct{})

	err = db.EvalQuery(queryRequest, col, &queryResult)
	if err != nil {
		return docId, nil, fmt.Errorf("failed to evaluate query: %v", err)
	}
	for id := range queryResult {
		readBack, err = col.Read(id)
		if err != nil {
			return docId, nil, fmt.Errorf("failed to read data from collection %v with id %v: %v", collectionName, id, err)
		}
		docId = id
	}
	return docId, readBack, nil
}

func (t *TieDot) Update(docId int, updateData map[string]interface{}, collectionName string) error {
	if t.debugging == "1" {
		t.logger.Info("update fields")
	}
	col := t.tdb.Use(collectionName)
	err := col.Update(docId, updateData)
	if err != nil {
		return fmt.Errorf("failed to update data in collection %v with id %v: %v", collectionName, docId, err)
	}
	return nil
}

func (t *TieDot) Delete(docId int, collectionName string) error {
	if t.debugging == "1" {
		t.logger.Info("delete data")
	}
	col := t.tdb.Use(collectionName)
	err := col.Delete(docId)
	if err != nil {
		return fmt.Errorf("failed to delete from collection %v data with id %v: %v", collectionName, docId, err)
	}
	return nil
}

func (t *TieDot) QueryGen(queryValue string, queryRule string, queryWhere string) (interface{}, error) {
	if t.debugging == "1" {
		t.logger.Info("generate query result")
	}
	var query interface{}
	str := fmt.Sprintf("[{\"%s\": \"%s\", \"in\": [\"%s\"]}]", queryRule, queryValue, queryWhere)
	err := json.Unmarshal([]byte(str), &query)
	if err != nil {
		return nil, fmt.Errorf("failed on json query unmarshaling: %v", err)
	}
	return query, nil
}
