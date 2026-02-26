package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/git-ecosystem/git-bundle-server/pkg/auth"
)

/* Built-in auth modes */
// Authorize users with credentials matching a static username/password pair
// that applies to the whole server.
type fixedCredentialAuth struct {
	usernameHash [32]byte
	passwordHash [32]byte
}

type fixedCredentialAuthParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

func NewFixedCredentialAuth(rawParameters json.RawMessage) (auth.AuthMiddleware, error) {
	if len(rawParameters) == 0 {
		return nil, fmt.Errorf("parameters JSON must exist")
	}

	var params fixedCredentialAuthParams
	err := json.Unmarshal(rawParameters, &params)
	if err != nil {
		return nil, err
	}

	// Check for invalid username characters
	if strings.Contains(params.Username, ":") {
		return nil, fmt.Errorf("username contains a colon (\":\")")
	}

	// Additional username validation
	if len(params.Username) == 0 {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(params.Username) > 255 {
		return nil, fmt.Errorf("username too long (max 255 characters)")
	}

	// Check for other potentially problematic characters
	invalidChars := []string{"\x00", "\n", "\r", "\t"}
	for _, char := range invalidChars {
		if strings.Contains(params.Username, char) {
			return nil, fmt.Errorf("username contains invalid character")
		}
	}

	// Make sure password hash is a valid hash
	passwordHashBytes, err := hex.DecodeString(params.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("passwordHash is invalid: %w", err)
	} else if len(passwordHashBytes) != 32 {
		return nil, fmt.Errorf("passwordHash is incorrect length (%d vs. expected 32)", len(passwordHashBytes))
	}

	return &fixedCredentialAuth{
		usernameHash: sha256.Sum256([]byte(params.Username)),
		passwordHash: [32]byte(passwordHashBytes),
	}, nil
}

func (a *fixedCredentialAuth) Authorize(r *http.Request, _ string, _ string) auth.AuthResult {
	username, password, ok := r.BasicAuth()

	// Always perform the same amount of work to prevent timing attacks
	var result auth.AuthResult
	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], a.usernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], a.passwordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			result = auth.Allow()
		} else {
			// Return a 404 status even though the issue is that the user is
			// forbidden so we don't indirectly reveal which repositories are
			// configured in the bundle server.
			result = auth.Deny(404)
		}
	} else {
		result = auth.Deny(401, auth.Header{Key: "WWW-Authenticate", Value: `Basic realm="restricted", charset="UTF-8"`})
	}

	// Add a small constant delay to further prevent timing attacks
	time.Sleep(50 * time.Millisecond)

	return result
}
