package flags

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		os.Clearenv()

		f, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "", f.AgentID)
		assert.Equal(t, "", f.EnvironmentID)
		assert.Equal(t, "", f.ProjectID)
	})

	t.Run("custom values", func(t *testing.T) {
		os.Clearenv()

		if err := os.Setenv("FLAGS_ENVIRONMENT_ID", "envId"); err != nil {
			assert.NoError(t, err)
		}
		if err := os.Setenv("FLAGS_PROJECT_ID", "projId"); err != nil {
			assert.NoError(t, err)
		}
		if err := os.Setenv("FLAGS_AGENT_ID", "agentId"); err != nil {
			assert.NoError(t, err)
		}

		f, err := Build()
		assert.NoError(t, err)
		assert.Equal(t, "envId", f.EnvironmentID)
		assert.Equal(t, "projId", f.ProjectID)
		assert.Equal(t, "agentId", f.AgentID)
	})
}
