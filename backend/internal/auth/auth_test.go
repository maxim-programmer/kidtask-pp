package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HashPassword / CheckPassword тесты
func TestAuthPassword(t *testing.T) {
	password := "mySecretPassword123"

	t.Run("hash is not equal to original password", func(t *testing.T) {
		hash, err := HashPassword(password)
		require.NoError(t, err)
		assert.NotEqual(t, password, hash)
	})

	t.Run("correct password passes verification", func(t *testing.T) {
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPassword(hash, password)
		assert.True(t, result)
	})

	t.Run("incorrect password fails verification", func(t *testing.T) {
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPassword(hash, "wrongPassword")
		assert.False(t, result)
	})

	t.Run("empty password fails verification", func(t *testing.T) {
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPassword(hash, "")
		assert.False(t, result)
	})
}

// GenerateToken / ParseToken тесты
func TestAuthToken(t *testing.T) {
	secret := "test-secret-key"
	userID := "12345"
	role := "child"
	parentID := "parent-123"

	t.Run("token can be parsed successfully", func(t *testing.T) {
		token, err := GenerateToken(userID, role, parentID, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ParseToken(token, secret)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, parentID, claims.ParentID)
	})

	t.Run("wrong secret fails to parse token", func(t *testing.T) {
		token, err := GenerateToken(userID, role, parentID, secret)
		require.NoError(t, err)

		claims, err := ParseToken(token, "wrong-secret-key")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("expired token fails to parse", func(t *testing.T) {
		// Создаем токен с просроченным временем через прямую манипуляцию
		expiredClaims := Claims{
			UserID:   userID,
			Role:     role,
			ParentID: parentID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString([]byte(secret))
		require.NoError(t, err)

		claims, err := ParseToken(token, secret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("malformed token fails to parse", func(t *testing.T) {
		claims, err := ParseToken("not-a-valid-token", secret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("token with wrong signing method fails", func(t *testing.T) {
		// Создаем токен с неправильным методом подписи
		wrongClaims := Claims{
			UserID:   userID,
			Role:     role,
			ParentID: parentID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			},
		}

		// Используем неправильный метод подписи (должен быть HS256)
		token, err := jwt.NewWithClaims(jwt.SigningMethodNone, wrongClaims).SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		claims, err := ParseToken(token, secret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
