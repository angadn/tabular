package tabular

import (
	"fmt"
	"testing"
)

func TestInsertion(t *testing.T) {
	tab := New("users", "id", "email", "password", "secret")
	fmt.Printf("Insertion:\t%s\n\n", tab.Insertion("%s"))
}

func TestSelection(t *testing.T) {
	tab := New("campaigns", "id", "user_id", "enabled", "name")
	userTab := New("users", "id", "email", "password", "secret")
	fmt.Printf(
		"Selection:\t%s\n\n",
		tab.Selection("SELECT %s FROM `users` WHERE...", userTab),
	)
}
