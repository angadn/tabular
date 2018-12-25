package tabular

import (
	"fmt"
	"testing"
)

func TestInsertion(t *testing.T) {
	tab := New("users", "id", "email", "password", "secret", "created_at", "updated_at")
	fmt.Printf("Insertion:\t%s\n\n", tab.Insertion(
		"%s",
		"!created_at",
		"NOW()",
		"updated_at",
		"NOW()",
	))
}

func TestBatchInsertion(t *testing.T) {
	tab := New("users", "id", "email", "password", "secret", "created_at", "updated_at")
	fmt.Printf("Batch Insertion:\t%s\n\n", tab.BatchInsertion(
		"%s",
		3,
		"!created_at",
		"NOW()",
		"updated_at",
		"NOW()",
	))
}

func TestSelection(t *testing.T) {
	tab := New("campaigns", "id", "user_id", "enabled", "name")
	userTab := New("users", "id", "email", "password", "secret")
	fmt.Printf(
		"Selection:\t%s\n\n",
		tab.Selection("SELECT %s FROM `users` WHERE ...", userTab),
	)
}

func TestDualJoinSelection(t *testing.T) {
	tab := New("campaigns", "id", "user_id", "secondary_user_id", "enabled", "name")
	userTab := New("users", "id", "email", "password", "secret")
	fmt.Printf(
		"Dual Join:\t%s\n\n",
		tab.Selection("SELECT %s FROM `users` WHERE ...", userTab, userTab.WithAlias("secondary_user")),
	)
}

func TestNullSelection(t *testing.T) {
	tab := New("campaigns", "id", "user_id", "enabled", "name").WithNullSelection()
	userTab := New("users", "id", "email", "password", "secret")
	fmt.Printf(
		"Null Selection:\t%s\n\n",
		tab.Selection("SELECT %s FROM `users` WHERE ...", userTab),
	)
}

func TestTablePrefix(t *testing.T) {
	tab := New("campaigns", "id", "user_id", "enabled", "name")
	userTab := New("users", "id", "email", "password", "secret")
	fmt.Printf(
		"Prefix Selection:\t%s\n\n",
		tab.PrefixedSelection("SELECT %s FROM `users` WHERE ...", userTab),
	)
}
