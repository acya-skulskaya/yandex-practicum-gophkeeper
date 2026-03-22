package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCredits_Valid(t *testing.T) {
	var tests = []struct {
		name      string
		userCreds UserCredits
		hasError  bool
		want      Problems
	}{
		{
			name: "valid",
			userCreds: UserCredits{
				Login:    "test",
				Password: "test",
			},
			hasError: false,
			want:     Problems{List: map[string]string{}},
		},
		{
			name: "not valid",
			userCreds: UserCredits{
				Login:    "",
				Password: "",
			},
			hasError: true,
			want: Problems{List: map[string]string{
				RequestFieldUsername: ErrorMsgEmptyLogin,
				RequestFieldPassword: ErrorMsgEmptyPassword,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems, err := IsValid(tt.userCreds)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, problems)
		})
	}
}
