package whatsapp

import (
	"context"
	"database/sql"
	"testing"
)

func TestAdmin(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	admin := NewAdmin(db)

	if err := admin.Migrate(ctx); err != nil {
		t.Fatal(err)
	}

	jid := "628123456789@s.whatsapp.net"

	isAdmin, err := admin.IsAdmin(ctx, jid)
	if err != nil {
		t.Fatal(err)
	}

	if isAdmin {
		t.Fatal("expected user to not be admin")
	}

	if err := admin.AddAdmin(ctx, jid, "admin"); err != nil {
		t.Fatal(err)
	}

	isAdmin, err = admin.IsAdmin(ctx, jid)
	if err != nil {
		t.Fatal(err)
	}

	if !isAdmin {
		t.Fatal("expected user to be admin")
	}

	if err := admin.RemoveAdmin(ctx, jid); err != nil {
		t.Fatal(err)
	}

	isAdmin, err = admin.IsAdmin(ctx, jid)
	if err != nil {
		t.Fatal(err)
	}

	if isAdmin {
		t.Fatal("expected user to no longer be admin")
	}
}
