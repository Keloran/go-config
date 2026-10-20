package resend

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("RESEND_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}

	r := NewSystem()
	rd, err := r.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", rd.Key)
}
