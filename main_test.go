package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func newTestDB(t *testing.T) *Driver {
	t.Helper()

	db, err := New(t.TempDir(), nil)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	return db
}

func sampleUser(name string) User {
	return User{
		Name:    name,
		Age:     json.Number("21"),
		Contact: "9999999999",
		Company: "ExampleCo",
		Address: Address{
			City:    "Bangalore",
			State:   "Karnataka",
			Country: "India",
			Pincode: json.Number("560001"),
		},
	}
}

func TestWriteAndReadRoundTrip(t *testing.T) {
	db := newTestDB(t)

	want := sampleUser("Alice")

	if err := db.Write("users", want.Name, want); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	var got User
	if err := db.Read("users", want.Name, &got); err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if got.Name != want.Name || got.Company != want.Company || got.Age != want.Age {
		t.Fatalf("round-trip mismatch: got=%+v want=%+v", got, want)
	}
}

func TestReadAllReturnsAllRecords(t *testing.T) {
	db := newTestDB(t)

	users := []User{sampleUser("Alice"), sampleUser("Bob"), sampleUser("Carol")}
	for _, u := range users {
		if err := db.Write("users", u.Name, u); err != nil {
			t.Fatalf("write failed for %s: %v", u.Name, err)
		}
	}

	records, err := db.ReadAll("users")
	if err != nil {
		t.Fatalf("read all failed: %v", err)
	}

	if len(records) != len(users) {
		t.Fatalf("unexpected record count: got=%d want=%d", len(records), len(users))
	}

	combined := strings.Join(records, "\n")
	for _, u := range users {
		if !strings.Contains(combined, "\"Name\": \""+u.Name+"\"") {
			t.Fatalf("record for user %s not found in ReadAll output", u.Name)
		}
	}
}

func TestDeleteRemovesRecord(t *testing.T) {
	db := newTestDB(t)

	u := sampleUser("ToDelete")
	if err := db.Write("users", u.Name, u); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if err := db.Delete("users", u.Name); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	var out User
	if err := db.Read("users", u.Name, &out); err == nil {
		t.Fatal("expected read to fail after delete, got nil error")
	}
}

func TestValidationErrors(t *testing.T) {
	db := newTestDB(t)

	if err := db.Write("", "alice", sampleUser("alice")); err == nil {
		t.Fatal("expected write to fail for empty collection")
	}

	if err := db.Write("users", "", sampleUser("alice")); err == nil {
		t.Fatal("expected write to fail for empty resource")
	}

	var u User
	if err := db.Read("", "alice", &u); err == nil {
		t.Fatal("expected read to fail for empty collection")
	}

	if err := db.Delete("users", ""); err == nil {
		t.Fatal("expected delete to fail for empty resource")
	}

	if _, err := db.ReadAll(""); err == nil {
		t.Fatal("expected read all to fail for empty collection")
	}
}
