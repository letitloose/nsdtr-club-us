package models

import (
	"database/sql"
	"testing"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func TestLegacyMemberList(t *testing.T) {
	db, err := sql.Open("mysql", "lougar:thewarrior@/nsdtrc_members?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}

	legacyModel := LegacyModel{DB: db}

	legacyMembers, err := legacyModel.List()
	if err != nil {
		t.Fatal(err)
	}

	expected := 1791
	if len(legacyMembers) != expected {
		t.Fatalf("wrong number of results returned. expecting %d, got %d", expected, len(legacyMembers))
	}
}
