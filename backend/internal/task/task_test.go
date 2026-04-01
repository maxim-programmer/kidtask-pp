package task
import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestTaskValidate(t *testing.T) {
	tests := []struct {
		name      string
		task      Task
		wantError bool
	}{
		{
			name: "empty title returns error",
			task: Task{
				Title:  "",
				Reward: 100,
			},
			wantError: true,
		},
		{
			name: "reward <= 0 returns error",
			task: Task{
				Title:  "Clean room",
				Reward: 0,
			},
			wantError: true,
		},
		{
			name: "negative reward returns error",
			task: Task{
				Title:  "Clean room",
				Reward: -10,
			},
			wantError: true,
		},
		{
			name: "valid task passes validation",
			task: Task{
				Title:  "Clean room",
				Reward: 100,
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
