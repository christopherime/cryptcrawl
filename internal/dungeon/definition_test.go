package dungeon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateExampleDungeon(t *testing.T) {
	def := CreateExampleDungeon()

	if def == nil {
		t.Fatal("CreateExampleDungeon returned nil")
	}

	if def.Name == "" {
		t.Error("Dungeon name is empty")
	}

	if len(def.Levels) == 0 {
		t.Error("Dungeon has no levels")
	}

	if len(def.Monsters) == 0 {
		t.Error("Dungeon has no monsters")
	}

	if len(def.Items) == 0 {
		t.Error("Dungeon has no items")
	}
}

func TestSaveAndLoadDungeonDefinition(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an example dungeon
	def := CreateExampleDungeon()

	// Save the dungeon definition
	path := filepath.Join(tempDir, "test-dungeon.json")
	err = SaveDungeonDefinition(def, path)
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}

	// Load the dungeon definition
	loadedDef, err := LoadDungeonDefinition(path)
	if err != nil {
		t.Fatalf("Failed to load dungeon definition: %v", err)
	}

	// Check that the loaded definition matches the original
	if loadedDef.Name != def.Name {
		t.Errorf("Loaded dungeon name does not match: got %q, want %q", loadedDef.Name, def.Name)
	}

	if len(loadedDef.Levels) != len(def.Levels) {
		t.Errorf("Loaded dungeon has different number of levels: got %d, want %d", len(loadedDef.Levels), len(def.Levels))
	}

	if len(loadedDef.Monsters) != len(def.Monsters) {
		t.Errorf("Loaded dungeon has different number of monsters: got %d, want %d", len(loadedDef.Monsters), len(def.Monsters))
	}

	if len(loadedDef.Items) != len(def.Items) {
		t.Errorf("Loaded dungeon has different number of items: got %d, want %d", len(loadedDef.Items), len(def.Items))
	}
}

func TestGenerateDungeonFromDefinition(t *testing.T) {
	def := CreateExampleDungeon()

	// Generate a dungeon from the definition
	dungeon, metadata, err := GenerateDungeonFromDefinition(def, 0)
	if err != nil {
		t.Fatalf("Failed to generate dungeon: %v", err)
	}

	// Check that the dungeon was generated correctly
	if len(dungeon) == 0 {
		t.Error("Generated dungeon is empty")
	}

	if metadata == nil {
		t.Error("Generated metadata is nil")
	}

	// Check that the metadata contains the expected fields
	if name, ok := metadata["name"].(string); !ok || name == "" {
		t.Error("Metadata does not contain a valid name")
	}

	if description, ok := metadata["description"].(string); !ok || description == "" {
		t.Error("Metadata does not contain a valid description")
	}

	// Check that rooms exist in metadata
	if _, ok := metadata["rooms"]; !ok {
		t.Error("Metadata does not contain rooms")
	}

	// Check that start position exists in metadata
	if _, ok := metadata["startPos"]; !ok {
		t.Error("Metadata does not contain a start position")
	}

	// Check that exit position exists in metadata
	if _, ok := metadata["exitPos"]; !ok {
		t.Error("Metadata does not contain an exit position")
	}
}

func TestLoadDungeonDefinitionFromDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create and save multiple dungeon definitions
	def1 := CreateExampleDungeon()
	def1.Name = "Dungeon 1"
	err = SaveDungeonDefinition(def1, filepath.Join(tempDir, "dungeon1.json"))
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}

	def2 := CreateExampleDungeon()
	def2.Name = "Dungeon 2"
	err = SaveDungeonDefinition(def2, filepath.Join(tempDir, "dungeon2.json"))
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}

	// Create a non-JSON file that should be ignored
	err = os.WriteFile(filepath.Join(tempDir, "not-a-dungeon.txt"), []byte("This is not a dungeon"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load all dungeon definitions from the directory
	defs, err := LoadDungeonDefinitionFromDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to load dungeon definitions: %v", err)
	}

	// Check that the correct number of definitions were loaded
	if len(defs) != 2 {
		t.Errorf("Expected 2 dungeon definitions, got %d", len(defs))
	}

	// Check that the definitions were loaded correctly
	foundDungeon1 := false
	foundDungeon2 := false
	for _, def := range defs {
		if def.Name == "Dungeon 1" {
			foundDungeon1 = true
		} else if def.Name == "Dungeon 2" {
			foundDungeon2 = true
		}
	}

	if !foundDungeon1 {
		t.Error("Dungeon 1 was not loaded")
	}

	if !foundDungeon2 {
		t.Error("Dungeon 2 was not loaded")
	}
}
