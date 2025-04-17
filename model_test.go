package main

import (
	"testing"
)

func TestInitialModel(t *testing.T) {
	m := initialModel()

	// Check that the model is initialized with the expected values
	if m.width != 97 {
		t.Errorf("Expected width to be 97, got %d", m.width)
	}

	if m.height != 30 {
		t.Errorf("Expected height to be 30, got %d", m.height)
	}

	if m.gold != 0 {
		t.Errorf("Expected gold to be 0, got %d", m.gold)
	}

	if m.level != 1 {
		t.Errorf("Expected level to be 1, got %d", m.level)
	}

	if m.gameOver {
		t.Errorf("Expected gameOver to be false, got true")
	}

	if m.gameWon {
		t.Errorf("Expected gameWon to be false, got true")
	}

	if len(m.messages) == 0 {
		t.Errorf("Expected messages to be initialized with a welcome message")
	}

	// Check that the dungeon is initialized
	if len(m.dungeon) == 0 {
		t.Errorf("Expected dungeon to be initialized")
	}

	// Check that the player is initialized
	if m.player.Health <= 0 {
		t.Errorf("Expected player health to be positive, got %d", m.player.Health)
	}

	if m.player.Damage <= 0 {
		t.Errorf("Expected player damage to be positive, got %d", m.player.Damage)
	}
}

func TestIsVisible(t *testing.T) {
	m := initialModel()
	playerX := m.player.Pos.X
	playerY := m.player.Pos.Y

	// Test visibility within range
	if !m.isVisible(playerX, playerY) {
		t.Errorf("Expected player position to be visible")
	}

	if !m.isVisible(playerX+1, playerY) {
		t.Errorf("Expected position adjacent to player to be visible")
	}

	if !m.isVisible(playerX, playerY+1) {
		t.Errorf("Expected position adjacent to player to be visible")
	}

	// Test visibility at the edge of range
	if !m.isVisible(playerX+5, playerY) {
		t.Errorf("Expected position at edge of visibility range to be visible")
	}

	// Test visibility outside range
	if m.isVisible(playerX+6, playerY) {
		t.Errorf("Expected position outside visibility range to be invisible")
	}

	if m.isVisible(playerX, playerY+6) {
		t.Errorf("Expected position outside visibility range to be invisible")
	}
}

func TestAddMessage(t *testing.T) {
	m := initialModel()
	initialMessageCount := len(m.messages)
	
	m.addMessage("Test message")
	
	if len(m.messages) != initialMessageCount+1 {
		t.Errorf("Expected message count to increase by 1, got %d, want %d", 
			len(m.messages), initialMessageCount+1)
	}
	
	if m.messages[len(m.messages)-1] != "Test message" {
		t.Errorf("Expected last message to be 'Test message', got '%s'", 
			m.messages[len(m.messages)-1])
	}
}

func TestAbsFunction(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{5, 5},
		{-5, 5},
		{10, 10},
		{-10, 10},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

func TestDungeonToString(t *testing.T) {
	m := initialModel()
	
	// Force reveal map to ensure consistent output
	m.revealMap = true
	
	result := m.dungeonToString()
	
	// Basic checks on the result
	if result == "" {
		t.Errorf("dungeonToString() returned empty string")
	}
	
	// Check that the result has the expected number of lines
	lines := 0
	for i := 0; i < len(result); i++ {
		if result[i] == '\n' {
			lines++
		}
	}
	
	if lines != m.height {
		t.Errorf("Expected %d lines in dungeonToString() result, got %d", m.height, lines)
	}
}
