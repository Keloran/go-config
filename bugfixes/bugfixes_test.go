package bugfixes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildGeneric(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("BUGFIXES_AGENT_KEY", "testKey"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("BUGFIXES_AGENT_SECRET", "testSecret"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("BUGFIXES_SERVER", "http://bob.bob"); err != nil {
		t.Fatal(err)
	}

	b := NewSystem()
	bf, err := b.Build()
	assert.NoError(t, err)
	assert.Equal(t, "testKey", bf.AgentKey)
	assert.Equal(t, "testSecret", bf.AgentSecret)
	assert.Equal(t, "http://bob.bob", bf.Server)
}

func TestBuildGenericInvalidHost(t *testing.T) {
	os.Clearenv()
	if err := os.Setenv("BUGFIXES_SERVER", "bob.bob"); err != nil {
		t.Fatal(err)
	}

	b := NewSystem()
	_, err := b.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bugfixes: unable to use server without protocol")
}
