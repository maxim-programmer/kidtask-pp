package parent
import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestParentValidate(t *testing.T) {
	tests := []struct {
		name      string
		parent    Parent
		wantError bool
	}{
		{
			name: "empty email returns error",
			parent: Parent{
				Email: "",
				Name:  "John",
			},
			wantError: true,
		},
		{
			name: "invalid email (no @) returns error",
			parent: Parent{
				Email: "invalid-email",
				Name:  "John",
			},
			wantError: true,
		},
		{
			name: "empty name returns error",
			parent: Parent{
				Email: "john@example.com",
				Name:  "",
			},
			wantError: true,
		},
		{
			name: "valid parent passes validation",
			parent: Parent{
				Email: "john@example.com",
				Name:  "John",
			},
			wantError: false,
		},
		{
			name: "email with numbers passes validation",
			parent: Parent{
				Email: "john123@example.com",
				Name:  "John",
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.parent.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
