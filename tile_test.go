package main

import (
	"testing"
)

func TestGetTileBySymbol(t *testing.T) {
	tests := []struct {
		name     string
		symbol   rune
		expected TileType
	}{
		{"Empty", ' ', Empty},
		{"Wall", '#', Wall},
		{"Player", '@', Player},
		{"Monster", 'M', Monster},
		{"Gold", '$', Gold},
		{"Exit", 'E', Exit},
		{"Unknown", 'X', Empty}, // Should default to Empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := GetTileBySymbol(tt.symbol)
			if tile.Type != tt.expected {
				t.Errorf("GetTileBySymbol(%q) = %v, want %v", tt.symbol, tile.Type, tt.expected)
			}
		})
	}
}

func TestGetTileByType(t *testing.T) {
	tests := []struct {
		name     string
		tileType TileType
		expected rune
	}{
		{"Empty", Empty, ' '},
		{"Wall", Wall, '#'},
		{"Player", Player, '@'},
		{"Monster", Monster, 'M'},
		{"Gold", Gold, '$'},
		{"Exit", Exit, 'E'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := GetTileByType(tt.tileType)
			if tile.Symbol != tt.expected {
				t.Errorf("GetTileByType(%v) symbol = %q, want %q", tt.tileType, tile.Symbol, tt.expected)
			}
		})
	}
}

func TestRenderTile(t *testing.T) {
	// This is a simple test to ensure RenderTile doesn't crash
	// We can't easily test the actual styling output
	for tileType := range TileMap {
		result := RenderTile(tileType)
		if result == "" {
			t.Errorf("RenderTile(%v) returned empty string", tileType)
		}
	}
}

func TestRenderSymbol(t *testing.T) {
	// This is a simple test to ensure RenderSymbol doesn't crash
	symbols := []rune{' ', '#', '@', 'M', '$', 'E', 'X'}
	for _, symbol := range symbols {
		result := RenderSymbol(symbol)
		if result == "" {
			t.Errorf("RenderSymbol(%q) returned empty string", symbol)
		}
	}
}

func TestTileMapCompleteness(t *testing.T) {
	// Ensure all tile types have an entry in the map
	for i := TileType(0); i <= Lava; i++ {
		if _, ok := TileMap[i]; !ok {
			t.Errorf("TileType %d is not defined in TileMap", i)
		}
	}
}

func TestTileProperties(t *testing.T) {
	// Test that all tiles have the expected properties
	for tileType, tile := range TileMap {
		if tile.Type != tileType {
			t.Errorf("Tile %v has incorrect Type: got %v, want %v", tileType, tile.Type, tileType)
		}

		if tile.Symbol == 0 {
			t.Errorf("Tile %v has no Symbol", tileType)
		}

		if tile.Description == "" {
			t.Errorf("Tile %v has no Description", tileType)
		}

		// Skip style check as it's not easily testable
		// The style is initialized in the TileMap
	}
}
