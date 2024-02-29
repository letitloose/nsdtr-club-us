package services

import (
	"database/sql"
	"testing"
	"time"

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

func TestAddMembership(t *testing.T) {
	db := models.NewTestDB(t)

	members := &models.MemberModel{DB: db}
	memberService := MemberService{MemberModel: members}

	memberID, err := memberService.Insert("Lou", "garwood", "", "", "", "", "", 1, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	membershipForm := MembershipForm{MemberID: memberID, Year: 2003, MembershipType: "SI", MembershipAmount: 25, RosterAmount: 5, HealthDonations: 10, RescueDonations: 10}

	err = memberService.AddMembership(&membershipForm)
	if err != nil {
		t.Fatal(err)
	}

	member, err := memberService.GetMemberProfile(memberID)
	if err != nil {
		t.Fatal(err)
	}

	expected := 10.0
	if member.Memberships[0].ResueAmount.Float64 != expected {
		t.Fatalf("rescue amount is wrong.  expected: %f but got: %f", expected, member.Memberships[0].ResueAmount.Float64)
	}
}
