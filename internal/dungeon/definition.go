package dungeon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
)

// DungeonDefinition represents a custom dungeon definition
type DungeonDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Author      string             `json:"author"`
	Version     string             `json:"version"`
	Levels      []LevelDefinition  `json:"levels"`
	Monsters    []MonsterTemplate  `json:"monsters"`
	Items       []ItemTemplate     `json:"items"`
	Events      []EventDefinition  `json:"events"`
}

// LevelDefinition represents a single level in a dungeon
type LevelDefinition struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Width       int                `json:"width"`
	Height      int                `json:"height"`
	Layout      []string           `json:"layout"`
	Rooms       []RoomDefinition   `json:"rooms"`
	Encounters  []EncounterSpawn   `json:"encounters"`
	Items       []ItemSpawn        `json:"items"`
	StartPos    Position           `json:"startPos"`
	ExitPos     Position           `json:"exitPos"`
}

// RoomDefinition represents a room in a level
type RoomDefinition struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	X           int                `json:"x"`
	Y           int                `json:"y"`
	Width       int                `json:"width"`
	Height      int                `json:"height"`
	Doors       []Position         `json:"doors"`
}

// EncounterSpawn defines where monsters spawn
type EncounterSpawn struct {
	MonsterID   string             `json:"monsterId"`
	Count       int                `json:"count"`
	MinLevel    int                `json:"minLevel"`
	MaxLevel    int                `json:"maxLevel"`
	Position    *Position          `json:"position"`
	RoomID      string             `json:"roomId"`
}

// ItemSpawn defines where items spawn
type ItemSpawn struct {
	ItemID      string             `json:"itemId"`
	Position    *Position          `json:"position"`
	RoomID      string             `json:"roomId"`
	Chance      float64            `json:"chance"`
}

// MonsterTemplate defines a monster type
type MonsterTemplate struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Symbol      string             `json:"symbol"`
	Color       string             `json:"color"`
	Health      int                `json:"health"`
	Damage      int                `json:"damage"`
	LevelScale  float64            `json:"levelScale"`
	Abilities   []string           `json:"abilities"`
	LootTable   []LootEntry        `json:"lootTable"`
}

// ItemTemplate defines an item type
type ItemTemplate struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Symbol      string             `json:"symbol"`
	Color       string             `json:"color"`
	Type        string             `json:"type"`
	Value       int                `json:"value"`
	Effects     []ItemEffect       `json:"effects"`
}

// ItemEffect defines an effect that an item can have
type ItemEffect struct {
	Type        string             `json:"type"`
	Value       int                `json:"value"`
	Duration    int                `json:"duration"`
}

// LootEntry defines an item that can be dropped by a monster
type LootEntry struct {
	ItemID      string             `json:"itemId"`
	Chance      float64            `json:"chance"`
	MinCount    int                `json:"minCount"`
	MaxCount    int                `json:"maxCount"`
}

// EventDefinition defines a game event
type EventDefinition struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Trigger     string             `json:"trigger"`
	Actions     []EventAction      `json:"actions"`
}

// EventAction defines an action that happens during an event
type EventAction struct {
	Type        string             `json:"type"`
	Target      string             `json:"target"`
	Value       interface{}        `json:"value"`
}

// Position represents a 2D position
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// LoadDungeonDefinition loads a dungeon definition from a file
func LoadDungeonDefinition(path string) (*DungeonDefinition, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dungeon definition file: %w", err)
	}

	var def DungeonDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse dungeon definition: %w", err)
	}

	return &def, nil
}

// LoadDungeonDefinitionFromDir loads all dungeon definitions from a directory
func LoadDungeonDefinitionFromDir(dir string) ([]*DungeonDefinition, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dungeon directory: %w", err)
	}

	var defs []*DungeonDefinition
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, file.Name())
		def, err := LoadDungeonDefinition(path)
		if err != nil {
			fmt.Printf("Warning: failed to load dungeon definition %s: %v\n", path, err)
			continue
		}

		defs = append(defs, def)
	}

	return defs, nil
}

// SaveDungeonDefinition saves a dungeon definition to a file
func SaveDungeonDefinition(def *DungeonDefinition, path string) error {
	data, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dungeon definition: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write dungeon definition file: %w", err)
	}

	return nil
}

// CreateExampleDungeon creates an example dungeon definition
func CreateExampleDungeon() *DungeonDefinition {
	return &DungeonDefinition{
		Name:        "The Forgotten Crypt",
		Description: "A dark and dangerous crypt filled with undead monsters and ancient treasures.",
		Author:      "CryptCrawl",
		Version:     "1.0.0",
		Levels: []LevelDefinition{
			{
				ID:          "level1",
				Name:        "Entrance Hall",
				Description: "The entrance to the crypt. Dusty and abandoned.",
				Width:       20,
				Height:      10,
				Layout: []string{
					"####################",
					"#........#.........#",
					"#........#.........#",
					"#........+.........#",
					"#........#.........#",
					"#........#.........#",
					"#........#####.....#",
					"#................E.#",
					"#.S................#",
					"####################",
				},
				Rooms: []RoomDefinition{
					{
						ID:          "entrance",
						Name:        "Entrance",
						Description: "The entrance to the crypt.",
						X:           1,
						Y:           1,
						Width:       8,
						Height:      6,
						Doors: []Position{
							{X: 9, Y: 3},
						},
					},
					{
						ID:          "main_hall",
						Name:        "Main Hall",
						Description: "The main hall of the crypt.",
						X:           10,
						Y:           1,
						Width:       9,
						Height:      8,
						Doors: []Position{
							{X: 9, Y: 3},
						},
					},
				},
				Encounters: []EncounterSpawn{
					{
						MonsterID: "skeleton",
						Count:     2,
						MinLevel:  1,
						MaxLevel:  1,
						RoomID:    "main_hall",
					},
					{
						MonsterID: "zombie",
						Count:     1,
						MinLevel:  1,
						MaxLevel:  2,
						Position:  &Position{X: 15, Y: 5},
					},
				},
				Items: []ItemSpawn{
					{
						ItemID:   "gold",
						RoomID:   "entrance",
						Chance:   0.8,
					},
					{
						ItemID:   "health_potion",
						Position: &Position{X: 12, Y: 2},
						Chance:   1.0,
					},
					{
						ItemID:   "rusty_sword",
						RoomID:   "main_hall",
						Chance:   0.5,
					},
				},
				StartPos: Position{X: 2, Y: 8},
				ExitPos:  Position{X: 17, Y: 7},
			},
		},
		Monsters: []MonsterTemplate{
			{
				ID:          "skeleton",
				Name:        "Skeleton",
				Description: "A reanimated skeleton wielding a rusty sword.",
				Symbol:      "S",
				Color:       "#ffffff",
				Health:      5,
				Damage:      2,
				LevelScale:  1.5,
				LootTable: []LootEntry{
					{
						ItemID:    "gold",
						Chance:    0.7,
						MinCount:  1,
						MaxCount:  5,
					},
					{
						ItemID:    "bone_shard",
						Chance:    0.3,
						MinCount:  1,
						MaxCount:  3,
					},
				},
			},
			{
				ID:          "zombie",
				Name:        "Zombie",
				Description: "A shambling corpse with rotting flesh.",
				Symbol:      "Z",
				Color:       "#00ff00",
				Health:      8,
				Damage:      1,
				LevelScale:  1.2,
				LootTable: []LootEntry{
					{
						ItemID:    "gold",
						Chance:    0.5,
						MinCount:  1,
						MaxCount:  3,
					},
					{
						ItemID:    "rotten_flesh",
						Chance:    0.6,
						MinCount:  1,
						MaxCount:  2,
					},
				},
			},
		},
		Items: []ItemTemplate{
			{
				ID:          "gold",
				Name:        "Gold",
				Description: "Shiny gold coins.",
				Symbol:      "$",
				Color:       "#ffff00",
				Type:        "currency",
				Value:       1,
			},
			{
				ID:          "health_potion",
				Name:        "Health Potion",
				Description: "A potion that restores health.",
				Symbol:      "!",
				Color:       "#ff0000",
				Type:        "consumable",
				Value:       10,
				Effects: []ItemEffect{
					{
						Type:  "heal",
						Value: 5,
					},
				},
			},
			{
				ID:          "rusty_sword",
				Name:        "Rusty Sword",
				Description: "An old, rusty sword. Still sharp enough to cut.",
				Symbol:      "/",
				Color:       "#aaaaaa",
				Type:        "weapon",
				Value:       5,
				Effects: []ItemEffect{
					{
						Type:  "damage",
						Value: 2,
					},
				},
			},
			{
				ID:          "bone_shard",
				Name:        "Bone Shard",
				Description: "A sharp shard of bone.",
				Symbol:      "*",
				Color:       "#ffffff",
				Type:        "material",
				Value:       2,
			},
			{
				ID:          "rotten_flesh",
				Name:        "Rotten Flesh",
				Description: "A piece of rotten flesh. Smells terrible.",
				Symbol:      "%",
				Color:       "#00aa00",
				Type:        "material",
				Value:       1,
			},
		},
		Events: []EventDefinition{
			{
				ID:          "entrance_event",
				Name:        "Entrance Event",
				Description: "An event that triggers when the player enters the dungeon.",
				Trigger:     "level_start",
				Actions: []EventAction{
					{
						Type:   "message",
						Value:  "You enter the forgotten crypt. The air is stale and cold.",
					},
					{
						Type:   "sound",
						Value:  "door_creak",
					},
				},
			},
			{
				ID:          "skeleton_death",
				Name:        "Skeleton Death",
				Description: "An event that triggers when a skeleton dies.",
				Trigger:     "monster_death",
				Actions: []EventAction{
					{
						Type:   "message",
						Value:  "The skeleton crumbles to dust!",
					},
					{
						Type:   "sound",
						Value:  "bone_crunch",
					},
				},
			},
		},
	}
}

// GenerateDungeonFromDefinition generates a dungeon from a definition
func GenerateDungeonFromDefinition(def *DungeonDefinition, level int) ([][]rune, map[string]interface{}, error) {
	if level < 0 || level >= len(def.Levels) {
		return nil, nil, fmt.Errorf("invalid level index: %d", level)
	}

	levelDef := def.Levels[level]
	
	// Create the dungeon grid
	dungeon := make([][]rune, levelDef.Height)
	for i := range dungeon {
		dungeon[i] = make([]rune, levelDef.Width)
		for j := range dungeon[i] {
			dungeon[i][j] = '#' // Default to walls
		}
	}
	
	// Parse the layout
	for y, row := range levelDef.Layout {
		if y >= levelDef.Height {
			break
		}
		for x, char := range row {
			if x >= levelDef.Width {
				break
			}
			dungeon[y][x] = char
		}
	}
	
	// Create metadata for the dungeon
	metadata := map[string]interface{}{
		"name":        levelDef.Name,
		"description": levelDef.Description,
		"rooms":       levelDef.Rooms,
		"startPos":    levelDef.StartPos,
		"exitPos":     levelDef.ExitPos,
		"monsters":    make([]map[string]interface{}, 0),
		"items":       make([]map[string]interface{}, 0),
	}
	
	// Place monsters
	monsterMap := make(map[string]MonsterTemplate)
	for _, monster := range def.Monsters {
		monsterMap[monster.ID] = monster
	}
	
	for _, encounter := range levelDef.Encounters {
		monster, ok := monsterMap[encounter.MonsterID]
		if !ok {
			continue
		}
		
		count := encounter.Count
		for i := 0; i < count; i++ {
			var x, y int
			
			if encounter.Position != nil {
				// Fixed position
				x, y = encounter.Position.X, encounter.Position.Y
			} else if encounter.RoomID != "" {
				// Random position in a room
				var room *RoomDefinition
				for _, r := range levelDef.Rooms {
					if r.ID == encounter.RoomID {
						room = &r
						break
					}
				}
				
				if room == nil {
					continue
				}
				
				// Find a random empty position in the room
				attempts := 0
				for attempts < 100 {
					rx := rand.Intn(room.Width-2) + room.X + 1
					ry := rand.Intn(room.Height-2) + room.Y + 1
					
					if rx < 0 || rx >= levelDef.Width || ry < 0 || ry >= levelDef.Height {
						attempts++
						continue
					}
					
					if dungeon[ry][rx] == '.' {
						x, y = rx, ry
						break
					}
					
					attempts++
				}
				
				if attempts >= 100 {
					continue
				}
			} else {
				// Random position anywhere in the dungeon
				attempts := 0
				for attempts < 100 {
					rx := rand.Intn(levelDef.Width)
					ry := rand.Intn(levelDef.Height)
					
					if dungeon[ry][rx] == '.' {
						x, y = rx, ry
						break
					}
					
					attempts++
				}
				
				if attempts >= 100 {
					continue
				}
			}
			
			// Place the monster
			if x >= 0 && x < levelDef.Width && y >= 0 && y < levelDef.Height {
				dungeon[y][x] = []rune(monster.Symbol)[0]
				
				// Add monster to metadata
				monsterLevel := encounter.MinLevel
				if encounter.MaxLevel > encounter.MinLevel {
					monsterLevel = encounter.MinLevel + rand.Intn(encounter.MaxLevel-encounter.MinLevel+1)
				}
				
				monsterData := map[string]interface{}{
					"id":          monster.ID,
					"name":        monster.Name,
					"description": monster.Description,
					"symbol":      monster.Symbol,
					"color":       monster.Color,
					"health":      int(float64(monster.Health) * (1.0 + float64(monsterLevel-1)*monster.LevelScale)),
					"damage":      int(float64(monster.Damage) * (1.0 + float64(monsterLevel-1)*monster.LevelScale*0.5)),
					"level":       monsterLevel,
					"position":    map[string]int{"x": x, "y": y},
				}
				
				monsters := metadata["monsters"].([]map[string]interface{})
				metadata["monsters"] = append(monsters, monsterData)
			}
		}
	}
	
	// Place items
	itemMap := make(map[string]ItemTemplate)
	for _, item := range def.Items {
		itemMap[item.ID] = item
	}
	
	for _, itemSpawn := range levelDef.Items {
		// Check if the item should spawn based on chance
		if rand.Float64() > itemSpawn.Chance {
			continue
		}
		
		item, ok := itemMap[itemSpawn.ItemID]
		if !ok {
			continue
		}
		
		var x, y int
		
		if itemSpawn.Position != nil {
			// Fixed position
			x, y = itemSpawn.Position.X, itemSpawn.Position.Y
		} else if itemSpawn.RoomID != "" {
			// Random position in a room
			var room *RoomDefinition
			for _, r := range levelDef.Rooms {
				if r.ID == itemSpawn.RoomID {
					room = &r
					break
				}
			}
			
			if room == nil {
				continue
			}
			
			// Find a random empty position in the room
			attempts := 0
			for attempts < 100 {
				rx := rand.Intn(room.Width-2) + room.X + 1
				ry := rand.Intn(room.Height-2) + room.Y + 1
				
				if rx < 0 || rx >= levelDef.Width || ry < 0 || ry >= levelDef.Height {
					attempts++
					continue
				}
				
				if dungeon[ry][rx] == '.' {
					x, y = rx, ry
					break
				}
				
				attempts++
			}
			
			if attempts >= 100 {
				continue
			}
		} else {
			// Random position anywhere in the dungeon
			attempts := 0
			for attempts < 100 {
				rx := rand.Intn(levelDef.Width)
				ry := rand.Intn(levelDef.Height)
				
				if dungeon[ry][rx] == '.' {
					x, y = rx, ry
					break
				}
				
				attempts++
			}
			
			if attempts >= 100 {
				continue
			}
		}
		
		// Place the item
		if x >= 0 && x < levelDef.Width && y >= 0 && y < levelDef.Height {
			dungeon[y][x] = []rune(item.Symbol)[0]
			
			// Add item to metadata
			itemData := map[string]interface{}{
				"id":          item.ID,
				"name":        item.Name,
				"description": item.Description,
				"symbol":      item.Symbol,
				"color":       item.Color,
				"type":        item.Type,
				"value":       item.Value,
				"position":    map[string]int{"x": x, "y": y},
			}
			
			items := metadata["items"].([]map[string]interface{})
			metadata["items"] = append(items, itemData)
		}
	}
	
	// Place player and exit
	dungeon[levelDef.StartPos.Y][levelDef.StartPos.X] = '@'
	dungeon[levelDef.ExitPos.Y][levelDef.ExitPos.X] = 'E'
	
	return dungeon, metadata, nil
}
