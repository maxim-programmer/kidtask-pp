package child
import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)
func TestAgeGroupFromBirthday(t *testing.T) {
	tests := []struct {
		name     string
		birthday *time.Time
		expected string
	}{
		{
			name:     "nil birthday returns junior",
			birthday: nil,
			expected: "junior",
		},
		{
			name: "birthday < 11 years returns junior",
			birthday: func() *time.Time {
				t := time.Now().AddDate(-10, 0, 0)
				return &t
			}(),
			expected: "junior",
		},
		{
			name: "birthday >= 11 years returns senior",
			birthday: func() *time.Time {
				t := time.Now().AddDate(-11, 0, 0)
				return &t
			}(),
			expected: "senior",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ageGroupFromBirthday(tt.birthday)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestChildValidate(t *testing.T) {
	tests := []struct {
		name      string
		child     Child
		wantError bool
	}{
		{
			name: "empty username returns error",
			child: Child{
				Username: "",
				Name:     "John",
			},
			wantError: true,
		},
		{
			name: "empty name returns error",
			child: Child{
				Username: "john123",
				Name:     "",
			},
			wantError: true,
		},
		{
			name: "valid child passes validation",
			child: Child{
				Username: "john123",
				Name:     "John",
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.child.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
