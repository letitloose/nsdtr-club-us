package models

import (
	"database/sql"
	"os"
	"testing"
)

func TestTeardown(t *testing.T) {
	// t.Skip()
	db, err := sql.Open("mysql", "lougar:thewarrior@/nsdtrc_test?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("../../sql/teardown.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	db.Close()
}
