package wish
import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestWishValidate(t *testing.T) {
	tests := []struct {
		name      string
		wish      Wish
		wantError bool
	}{
		{
			name: "empty title returns error",
			wish: Wish{
				Title: "",
			},
			wantError: true,
		},
		{
			name: "valid wish passes validation",
			wish: Wish{
				Title: "New bicycle",
			},
			wantError: false,
		},
		{
			name: "wish with description passes validation",
			wish: Wish{
				Title:       "New bicycle",
				Description: stringPtr("Red bicycle with basket"),
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.wish.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func stringPtr(s string) *string {
	return &s
}
