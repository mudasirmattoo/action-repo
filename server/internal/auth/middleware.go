package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTKeys struct {
	Keys []JWTKey `json:"keys"`
}

type JWTKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Cache for JWT keys
var jwksCache = struct {
	keys   *JWTKeys
	expiry time.Time
}{}

// SupabaseAuthMiddleware validates JWT tokens issued by Supabase
func SupabaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the Authorization header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			return
		}

		// Set user ID in the context for later use
		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
			return
		}

		// Set claims in context
		c.Set("userId", userID)
		c.Set("claims", claims)

		// Continue to the next middleware/handler
		c.Next()
	}
}

// validateToken validates a JWT token against Supabase's JWK
func validateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the key function
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from the token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("key ID not found in token header")
		}

		// Get the public key
		publicKey, err := getPublicKey(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})
}

// getPublicKey fetches and returns the RSA public key for the given key ID
func getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Get Supabase URL from environment
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		return nil, fmt.Errorf("SUPABASE_URL environment variable not set")
	}

	// Check if we need to refresh the keys
	if jwksCache.keys == nil || time.Now().After(jwksCache.expiry) {
		// Fetch JWKS from Supabase
		jwksURL := fmt.Sprintf("%s/auth/v1/jwks", strings.TrimSuffix(supabaseURL, "/"))
		resp, err := http.Get(jwksURL)
		if err != nil {
			return nil, fmt.Errorf("error fetching JWKS: %v", err)
		}
		defer resp.Body.Close()

		// Decode JWKS
		var keys JWTKeys
		if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
			return nil, fmt.Errorf("error decoding JWKS: %v", err)
		}

		// Update cache
		jwksCache.keys = &keys
		jwksCache.expiry = time.Now().Add(24 * time.Hour) // Cache for 24 hours
	}

	// Find the key with matching ID
	for _, key := range jwksCache.keys.Keys {
		if key.Kid == kid {
			// Parse the RSA public key
			n, err := base64URLDecode(key.N)
			if err != nil {
				return nil, fmt.Errorf("invalid modulus: %v", err)
			}

			e, err := base64URLDecode(key.E)
			if err != nil {
				return nil, fmt.Errorf("invalid exponent: %v", err)
			}

			// Convert n to big.Int
			nInt := new(big.Int).SetBytes(n)

			// Convert e to int
			var eInt int
			for i := 0; i < len(e); i++ {
				eInt = eInt<<8 | int(e[i])
			}

			// Create and return RSA public key
			return &rsa.PublicKey{
				N: nInt,
				E: eInt,
			}, nil
		}
	}

	return nil, fmt.Errorf("key ID not found: %s", kid)
}

// base64URLDecode decodes a base64url encoded string
func base64URLDecode(str string) ([]byte, error) {
	// Add padding if necessary
	padding := 4 - len(str)%4
	if padding < 4 {
		str += strings.Repeat("=", padding)
	}

	// Replace base64url encoding with standard base64 encoding
	str = strings.ReplaceAll(str, "-", "+")
	str = strings.ReplaceAll(str, "_", "/")

	// Decode
	return base64.StdEncoding.DecodeString(str)
}
