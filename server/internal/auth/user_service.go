package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// UserService provides methods for user management via Supabase API
type UserService struct {
	SupabaseURL      string
	SupabaseAPIKey   string
	ServiceRoleKey   string
	AdminHeaderName  string
	AdminHeaderValue string
}

// NewUserService creates a new UserService instance
func NewUserService() (*UserService, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAPIKey := os.Getenv("SUPABASE_ANON_KEY")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	adminHeaderName := os.Getenv("SUPABASE_ADMIN_HEADER_NAME")
	adminHeaderValue := os.Getenv("SUPABASE_ADMIN_HEADER_VALUE")

	if supabaseURL == "" || supabaseAPIKey == "" {
		return nil, fmt.Errorf("missing required environment variables for Supabase")
	}

	// Ensure URL ends with a '/'
	if !strings.HasSuffix(supabaseURL, "/") {
		supabaseURL = supabaseURL + "/"
	}

	return &UserService{
		SupabaseURL:      supabaseURL,
		SupabaseAPIKey:   supabaseAPIKey,
		ServiceRoleKey:   serviceRoleKey,
		AdminHeaderName:  adminHeaderName,
		AdminHeaderValue: adminHeaderValue,
	}, nil
}

// UserProfile represents a user's profile data
type UserProfile struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

// GetUser retrieves a user's profile by ID
func (s *UserService) GetUser(userID string) (*UserProfile, error) {
	url := fmt.Sprintf("%sauth/v1/admin/users/%s", s.SupabaseURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Use service role key if available, otherwise use anon key
	apiKey := s.ServiceRoleKey
	if apiKey == "" {
		apiKey = s.SupabaseAPIKey
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	// Add admin header if configured
	if s.AdminHeaderName != "" && s.AdminHeaderValue != "" {
		req.Header.Set(s.AdminHeaderName, s.AdminHeaderValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", resp.Status)
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &profile, nil
}

// UpdateUserProfile updates a user's profile data
func (s *UserService) UpdateUserProfile(userID string, profile *UserProfile) error {
	url := fmt.Sprintf("%srest/v1/profiles", s.SupabaseURL)
	
	// Prepare the request body
	body, err := json.Marshal(map[string]interface{}{
		"id":          userID,
		"first_name":  profile.FirstName,
		"last_name":   profile.LastName,
		"display_name": profile.DisplayName,
		"avatar_url":  profile.AvatarURL,
	})
	if err != nil {
		return fmt.Errorf("error marshaling profile: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("apikey", s.SupabaseAPIKey)
	req.Header.Set("Authorization", "Bearer "+s.SupabaseAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response: %s", resp.Status)
	}

	return nil
}
