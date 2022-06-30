package hash

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestSHA1Hasher_Hash(t *testing.T) {
	t.Run("Print Hashed Password", func(t *testing.T) {
		h := NewSHA1Hasher("pswd_secret_phrase")

		// TestUser0Password
		for i := 0; i < 3; i++ {
			username := "TestUser" + strconv.Itoa(i+1)
			hashedPassword, err := h.Hash("test_password")
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(username+": "+hashedPassword, uuid.New())
		}
	})
}
