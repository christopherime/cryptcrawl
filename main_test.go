package main

import (
	"testing"
)

func TestTeaHandler(t *testing.T) {
	// This is a simple test to ensure teaHandler doesn't crash
	// We can't easily test the actual SSH session
	
	// Just check that initialModel doesn't crash
	m := initialModel()
	if m.width != 97 {
		t.Errorf("Expected width to be 97, got %d", m.width)
	}
	
	if m.height != 30 {
		t.Errorf("Expected height to be 30, got %d", m.height)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Default value when env not set",
			key:          "TEST_ENV_1",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Env value when set",
			key:          "TEST_ENV_2",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if needed
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}
			
			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
