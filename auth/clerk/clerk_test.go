package clerk

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("CLERK_SECRET_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY", "testPublicKey"); err != nil {
		t.Fatal(err)
	}

	c := NewSystem()
	ck, err := c.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", ck.Key)
	assert.Equal(t, "testPublicKey", ck.PublicKey)
}
