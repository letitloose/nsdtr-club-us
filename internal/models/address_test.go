package models

import (
	"database/sql"
	"testing"
)

func TestInsertAddress(t *testing.T) {
	db := NewTestDB(t)

	am := AddressModel{DB: db}

	address := &Address{City: sql.NullString{String: "NY", Valid: true}, StateProvince: sql.NullString{String: "NY", Valid: true}, CountryCode: sql.NullString{String: "USA", Valid: true}}
	_, err := am.Insert(address)
	if err != nil {
		t.Fatal(err)
	}
}
