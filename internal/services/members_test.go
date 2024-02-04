package services

import (
	"database/sql"
	"testing"

	"github.com/letitloose/nsdtr-club-us/internal/models"
)

func TestMigrateLegacyMembers(t *testing.T) {
	db := models.NewTestDB(t)

	legacydb, err := sql.Open("mysql", "lougar:thewarrior@/nsdtrc_members?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}

	members := &models.MemberModel{DB: db}
	legacyModel := &models.LegacyModel{DB: legacydb}
	memberService := MemberService{MemberModel: members, Legacy: legacyModel}

	err = memberService.MigrateLegacyMembers()
	if err != nil {
		t.Fatal(err)
	}
}
