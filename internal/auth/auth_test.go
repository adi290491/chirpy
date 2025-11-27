package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tests := []struct {
		name          string
		makeTokenFunc func() (string, error)
		secretToUse   string
		expectErr     bool
		expectUserID  uuid.UUID
	}{
		{
			name: "valid token",
			makeTokenFunc: func() (string, error) {
				return MakeJWT(userID, secret, time.Hour)
			},
			secretToUse:  secret,
			expectErr:    false,
			expectUserID: userID,
		},
		{
			name: "invalid token",
			makeTokenFunc: func() (string, error) {
				return "invalid-token", nil
			},
			secretToUse: secret,
			expectErr:   true,
		},
		{
			name: "expired token",
			makeTokenFunc: func() (string, error) {
				return MakeJWT(userID, secret, -1*time.Hour)
			},
			secretToUse: secret,
			expectErr:   true,
		},

		{
			name: "wrong secret",
			makeTokenFunc: func() (string, error) {
				return MakeJWT(userID, secret, time.Hour)
			},
			secretToUse: "WRONG-SECRET",
			expectErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := tc.makeTokenFunc()
			if err != nil {
				t.Fatalf("MakeJWT failed unexpectedly: %v", err)
			}

			parsedID, err := ValidateJWT(token, tc.secretToUse)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if parsedID != tc.expectUserID {
				t.Fatalf("expected userID %s, got %s", tc.expectUserID.String(), parsedID.String())
			}
		})
	}
}
