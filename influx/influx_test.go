package influx

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("INFLUX_HOSTNAME", "testHost"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_TOKEN", "testToken"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_BUCKET", "testBucket"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INFLUX_ORG", "testOrg"); err != nil {
		t.Fatal(err)
	}

	i := NewSystem()
	in, err := i.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testToken", in.Token)
	assert.Equal(t, "testBucket", in.Bucket)
	assert.Equal(t, "testOrg", in.Org)
	assert.Equal(t, "testHost", in.Host)
}
