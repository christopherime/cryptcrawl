package dungeon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDungeonLoader(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-loader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a new dungeon loader
	loader, err := NewDungeonLoader(tempDir)
	if err != nil {
		t.Fatalf("Failed to create dungeon loader: %v", err)
	}
	
	// Check that the loader was initialized correctly
	if loader == nil {
		t.Fatal("NewDungeonLoader returned nil")
	}
	
	if loader.BasePath != tempDir {
		t.Errorf("Loader has incorrect base path: got %q, want %q", loader.BasePath, tempDir)
	}
	
	if len(loader.Dungeons) == 0 {
		t.Error("Loader has no dungeons")
	}
	
	// Check that the examples directory was created
	examplesDir := filepath.Join(tempDir, "examples")
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Error("Examples directory was not created")
	}
	
	// Check that an example dungeon was created
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}
	
	hasExampleDungeon := false
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			hasExampleDungeon = true
			break
		}
	}
	
	if !hasExampleDungeon {
		t.Error("No example dungeon was created")
	}
}

func TestDungeonLoaderNavigation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-loader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create multiple dungeon definitions
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
	
	// Create a new dungeon loader
	loader, err := NewDungeonLoader(tempDir)
	if err != nil {
		t.Fatalf("Failed to create dungeon loader: %v", err)
	}
	
	// Check that the loader has the correct number of dungeons
	if len(loader.Dungeons) != 2 {
		t.Errorf("Loader has incorrect number of dungeons: got %d, want 2", len(loader.Dungeons))
	}
	
	// Check that the current dungeon is the first one
	current := loader.GetCurrentDungeon()
	if current == nil {
		t.Fatal("GetCurrentDungeon returned nil")
	}
	
	// Navigate to the next dungeon
	next := loader.NextDungeon()
	if next == nil {
		t.Fatal("NextDungeon returned nil")
	}
	
	if next.Name == current.Name {
		t.Error("NextDungeon did not change the current dungeon")
	}
	
	// Navigate back to the first dungeon
	prev := loader.PrevDungeon()
	if prev == nil {
		t.Fatal("PrevDungeon returned nil")
	}
	
	if prev.Name != current.Name {
		t.Error("PrevDungeon did not return to the first dungeon")
	}
}

func TestGetDungeonByName(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-loader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create multiple dungeon definitions
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
	
	// Create a new dungeon loader
	loader, err := NewDungeonLoader(tempDir)
	if err != nil {
		t.Fatalf("Failed to create dungeon loader: %v", err)
	}
	
	// Get a dungeon by name
	dungeon := loader.GetDungeonByName("Dungeon 1")
	if dungeon == nil {
		t.Fatal("GetDungeonByName returned nil")
	}
	
	if dungeon.Name != "Dungeon 1" {
		t.Errorf("GetDungeonByName returned incorrect dungeon: got %q, want %q", dungeon.Name, "Dungeon 1")
	}
	
	// Try to get a non-existent dungeon
	dungeon = loader.GetDungeonByName("Non-existent Dungeon")
	if dungeon != nil {
		t.Error("GetDungeonByName returned a dungeon for a non-existent name")
	}
}

func TestReloadDungeons(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-loader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a dungeon definition
	def1 := CreateExampleDungeon()
	def1.Name = "Dungeon 1"
	err = SaveDungeonDefinition(def1, filepath.Join(tempDir, "dungeon1.json"))
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}
	
	// Create a new dungeon loader
	loader, err := NewDungeonLoader(tempDir)
	if err != nil {
		t.Fatalf("Failed to create dungeon loader: %v", err)
	}
	
	// Check that the loader has one dungeon
	if len(loader.Dungeons) != 1 {
		t.Errorf("Loader has incorrect number of dungeons: got %d, want 1", len(loader.Dungeons))
	}
	
	// Create another dungeon definition
	def2 := CreateExampleDungeon()
	def2.Name = "Dungeon 2"
	err = SaveDungeonDefinition(def2, filepath.Join(tempDir, "dungeon2.json"))
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}
	
	// Reload the dungeons
	err = loader.ReloadDungeons()
	if err != nil {
		t.Fatalf("Failed to reload dungeons: %v", err)
	}
	
	// Check that the loader now has two dungeons
	if len(loader.Dungeons) != 2 {
		t.Errorf("Loader has incorrect number of dungeons after reload: got %d, want 2", len(loader.Dungeons))
	}
}

func TestGenerateCurrentLevel(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dungeon-loader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a dungeon definition
	def := CreateExampleDungeon()
	err = SaveDungeonDefinition(def, filepath.Join(tempDir, "dungeon.json"))
	if err != nil {
		t.Fatalf("Failed to save dungeon definition: %v", err)
	}
	
	// Create a new dungeon loader
	loader, err := NewDungeonLoader(tempDir)
	if err != nil {
		t.Fatalf("Failed to create dungeon loader: %v", err)
	}
	
	// Generate a level
	dungeon, metadata, err := loader.GenerateCurrentLevel(0)
	if err != nil {
		t.Fatalf("Failed to generate level: %v", err)
	}
	
	// Check that the dungeon was generated correctly
	if len(dungeon) == 0 {
		t.Error("Generated dungeon is empty")
	}
	
	if metadata == nil {
		t.Error("Generated metadata is nil")
	}
}
