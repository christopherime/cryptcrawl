package main

import (
	"github.com/charmbracelet/lipgloss"
)

// TileType represents a type of dungeon tile
type TileType int

// Tile types
const (
	Empty TileType = iota
	Wall
	Player
	Monster
	Gold
	Exit
	Trap
	Chest
	Door
	Water
	Lava
)

// Tile represents a dungeon tile with a type and visual representation
type Tile struct {
	Type        TileType
	Symbol      rune
	Style       lipgloss.Style
	Walkable    bool
	Description string
}

// TileMap maps tile types to their visual representation
var TileMap = map[TileType]Tile{
	Empty: {
		Type:        Empty,
		Symbol:      ' ',
		Style:       lipgloss.NewStyle(),
		Walkable:    true,
		Description: "An empty floor tile.",
	},
	Wall: {
		Type:        Wall,
		Symbol:      '#',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Background(lipgloss.Color("#333333")),
		Walkable:    false,
		Description: "A solid stone wall.",
	},
	Player: {
		Type:        Player,
		Symbol:      '@',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff")).Bold(true),
		Walkable:    false,
		Description: "That's you!",
	},
	Monster: {
		Type:        Monster,
		Symbol:      'M',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true),
		Walkable:    false,
		Description: "A dangerous monster.",
	},
	Gold: {
		Type:        Gold,
		Symbol:      '$',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00")).Bold(true),
		Walkable:    true,
		Description: "Shiny gold coins.",
	},
	Exit: {
		Type:        Exit,
		Symbol:      'E',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true),
		Walkable:    true,
		Description: "An exit to the next level.",
	},
	Trap: {
		Type:        Trap,
		Symbol:      '^',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ff00ff")),
		Walkable:    true,
		Description: "A dangerous trap.",
	},
	Chest: {
		Type:        Chest,
		Symbol:      '?',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaa00")).Bold(true),
		Walkable:    true,
		Description: "A mysterious chest.",
	},
	Door: {
		Type:        Door,
		Symbol:      '+',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#aa5500")),
		Walkable:    true,
		Description: "A door.",
	},
	Water: {
		Type:        Water,
		Symbol:      '~',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")),
		Walkable:    false,
		Description: "Deep water.",
	},
	Lava: {
		Type:        Lava,
		Symbol:      '~',
		Style:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5500")).Background(lipgloss.Color("#aa0000")),
		Walkable:    false,
		Description: "Deadly lava.",
	},
}

// GetTileBySymbol returns a tile by its symbol
func GetTileBySymbol(symbol rune) Tile {
	for _, tile := range TileMap {
		if tile.Symbol == symbol {
			return tile
		}
	}
	return TileMap[Empty]
}

// GetTileByType returns a tile by its type
func GetTileByType(tileType TileType) Tile {
	return TileMap[tileType]
}

// RenderTile returns a styled string representation of a tile
func RenderTile(tileType TileType) string {
	tile := TileMap[tileType]
	return tile.Style.Render(string(tile.Symbol))
}

// RenderSymbol returns a styled string representation of a symbol
func RenderSymbol(symbol rune) string {
	for _, tile := range TileMap {
		if tile.Symbol == symbol {
			return tile.Style.Render(string(tile.Symbol))
		}
	}
	return string(symbol)
}
