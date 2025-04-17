package dungeon

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// DungeonLoader handles loading and managing dungeon definitions
type DungeonLoader struct {
	Dungeons     []*DungeonDefinition
	CurrentIndex int
	BasePath     string
}

// NewDungeonLoader creates a new dungeon loader
func NewDungeonLoader(basePath string) (*DungeonLoader, error) {
	// Create the dungeons directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create dungeons directory: %w", err)
	}

	// Create the examples directory if it doesn't exist
	examplesDir := filepath.Join(basePath, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create examples directory: %w", err)
	}

	// Create an example dungeon if none exist
	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dungeons directory: %w", err)
	}

	hasJsonFiles := false
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			hasJsonFiles = true
			break
		}
	}

	if !hasJsonFiles {
		// Check if there are examples
		exampleFiles, err := os.ReadDir(examplesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read examples directory: %w", err)
		}

		hasExamples := false
		for _, file := range exampleFiles {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
				hasExamples = true
				break
			}
		}

		if !hasExamples {
			// Create an example dungeon
			example := CreateExampleDungeon()
			examplePath := filepath.Join(examplesDir, "example_dungeon.json")
			if err := SaveDungeonDefinition(example, examplePath); err != nil {
				return nil, fmt.Errorf("failed to save example dungeon: %w", err)
			}
			log.Println("Created example dungeon at", examplePath)
		}
	}

	// Load all dungeon definitions
	dungeons, err := LoadDungeonDefinitionFromDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load dungeon definitions: %w", err)
	}

	// If no dungeons were found in the main directory, try loading from examples
	if len(dungeons) == 0 {
		dungeons, err = LoadDungeonDefinitionFromDir(examplesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load example dungeon definitions: %w", err)
		}
	}

	if len(dungeons) == 0 {
		return nil, fmt.Errorf("no dungeon definitions found")
	}

	return &DungeonLoader{
		Dungeons:     dungeons,
		CurrentIndex: 0,
		BasePath:     basePath,
	}, nil
}

// GetCurrentDungeon returns the current dungeon definition
func (dl *DungeonLoader) GetCurrentDungeon() *DungeonDefinition {
	if dl.CurrentIndex < 0 || dl.CurrentIndex >= len(dl.Dungeons) {
		return nil
	}
	return dl.Dungeons[dl.CurrentIndex]
}

// NextDungeon moves to the next dungeon definition
func (dl *DungeonLoader) NextDungeon() *DungeonDefinition {
	dl.CurrentIndex = (dl.CurrentIndex + 1) % len(dl.Dungeons)
	return dl.GetCurrentDungeon()
}

// PrevDungeon moves to the previous dungeon definition
func (dl *DungeonLoader) PrevDungeon() *DungeonDefinition {
	dl.CurrentIndex = (dl.CurrentIndex - 1 + len(dl.Dungeons)) % len(dl.Dungeons)
	return dl.GetCurrentDungeon()
}

// GetDungeonByName returns a dungeon definition by name
func (dl *DungeonLoader) GetDungeonByName(name string) *DungeonDefinition {
	for _, dungeon := range dl.Dungeons {
		if dungeon.Name == name {
			return dungeon
		}
	}
	return nil
}

// ReloadDungeons reloads all dungeon definitions
func (dl *DungeonLoader) ReloadDungeons() error {
	dungeons, err := LoadDungeonDefinitionFromDir(dl.BasePath)
	if err != nil {
		return fmt.Errorf("failed to reload dungeon definitions: %w", err)
	}

	// If no dungeons were found in the main directory, try loading from examples
	if len(dungeons) == 0 {
		dungeons, err = LoadDungeonDefinitionFromDir(filepath.Join(dl.BasePath, "examples"))
		if err != nil {
			return fmt.Errorf("failed to reload example dungeon definitions: %w", err)
		}
	}

	if len(dungeons) == 0 {
		return fmt.Errorf("no dungeon definitions found")
	}

	dl.Dungeons = dungeons
	dl.CurrentIndex = 0
	return nil
}

// GenerateCurrentLevel generates a dungeon from the current definition and level
func (dl *DungeonLoader) GenerateCurrentLevel(level int) ([][]rune, map[string]interface{}, error) {
	dungeon := dl.GetCurrentDungeon()
	if dungeon == nil {
		return nil, nil, fmt.Errorf("no current dungeon")
	}

	return GenerateDungeonFromDefinition(dungeon, level)
}
