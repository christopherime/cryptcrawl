# CryptCrawl

A MUD-style dungeon crawler accessible via SSH, built with Go and the charm.sh libraries. Create custom dungeons and share them with friends!

## Description

CryptCrawl is a terminal-based dungeon crawler game that users can connect to via SSH. It features:

- Procedurally generated dungeons
- Custom dungeon creation through JSON configuration files
- Colorful terminal UI with visibility system
- Turn-based combat
- Multiple levels to explore
- Gold collection and items
- Monsters that get tougher as you progress
- Traps, chests, and other interactive elements

## Table of Contents

- [CryptCrawl](#cryptcrawl)
  - [Description](#description)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Building from Source](#building-from-source)
  - [How to Play](#how-to-play)
  - [Controls](#controls)
  - [Creating Custom Dungeons](#creating-custom-dungeons)
    - [Getting Started with Custom Dungeons](#getting-started-with-custom-dungeons)
    - [Dungeon Structure](#dungeon-structure)
    - [Level Layout](#level-layout)
    - [Rooms](#rooms)
    - [Monsters](#monsters)
    - [Items](#items)
    - [Events](#events)
  - [Development](#development)
    - [Project Structure](#project-structure)
    - [Building from Source](#building-from-source-1)
    - [Running Tests](#running-tests)
    - [Environment Variables](#environment-variables)
  - [License](#license)


## Installation

### Prerequisites

- Go 1.20 or higher
- SSH client

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/cryptcrawl.git
   cd cryptcrawl
   ```

2. Build the application:

   ```bash
   make build
   ```

3. Generate SSH keys (if not already done):

   ```bash
   make ssh-key
   ```

## How to Play

1. Start the server:

   ```bash
   make run
   ```

2. Connect via SSH:

   ```bash
   make ssh-connect
   ```

   Or directly:

   ```bash
   ssh localhost -p 23234
   ```

3. Use the arrow keys or WASD to move around the dungeon.
4. Press space to attack monsters adjacent to you.
5. Collect gold and find the exit to progress to the next level.
6. Escape from the third level to win the game!

## Controls

- Arrow keys / WASD / HJKL: Move
- Space: Attack adjacent monsters
- ?: Toggle help
- Q / Ctrl+C: Quit

## Creating Custom Dungeons

CryptCrawl allows you to create custom dungeons using JSON configuration files. These files define the layout, monsters, items, and events in your dungeon.

### Getting Started with Custom Dungeons

1. Create a new dungeon template:

   ```bash
   make new-dungeon
   ```

   This will create a file called `my_dungeon.json` in the `dungeons` directory.

2. Edit the file to customize your dungeon.
3. Start the game, and it will automatically load your custom dungeon.

### Dungeon Structure

A dungeon definition consists of the following main components:

- **Metadata**: Name, description, author, and version
- **Levels**: One or more levels that make up the dungeon
- **Monsters**: Definitions of monster types that can appear in the dungeon
- **Items**: Definitions of item types that can appear in the dungeon
- **Events**: Special events that can be triggered during gameplay

Here's a basic example of a dungeon definition:

```json
{
  "name": "My Custom Dungeon",
  "description": "A dangerous dungeon filled with traps and treasures.",
  "author": "Your Name",
  "version": "1.0.0",
  "levels": [...],
  "monsters": [...],
  "items": [...],
  "events": [...]
}
```

### Level Layout

Each level has a layout defined as an array of strings, where each character represents a tile:

- `#`: Wall
- `.`: Empty space
- `+`: Door
- `@`: Player starting position (use `S` in the layout)
- `E`: Exit to the next level
- `M`, `S`, `Z`, `W`: Monster (different types)
- `$`: Gold
- `?`: Chest
- `^`: Trap
- `~`: Water or lava

Example level layout:

```json
"layout": [
  "####################",
  "#........#.........#",
  "#........#.........#",
  "#........+.........#",
  "#........#.........#",
  "#........#.........#",
  "#........#####.....#",
  "#................E.#",
  "#.S................#",
  "####################"
]
```

### Rooms

Rooms are defined areas within a level. They have a name, description, position, size, and doors:

```json
"rooms": [
  {
    "id": "entrance",
    "name": "Entrance",
    "description": "The entrance to the dungeon.",
    "x": 1,
    "y": 1,
    "width": 8,
    "height": 6,
    "doors": [
      {
        "x": 9,
        "y": 3
      }
    ]
  }
]
```

### Monsters

Monsters are defined with their stats, appearance, and loot tables:

```json
"monsters": [
  {
    "id": "skeleton",
    "name": "Skeleton",
    "description": "A reanimated skeleton wielding a rusty sword.",
    "symbol": "S",
    "color": "#ffffff",
    "health": 5,
    "damage": 2,
    "levelScale": 1.5,
    "abilities": [],
    "lootTable": [
      {
        "itemId": "gold",
        "chance": 0.7,
        "minCount": 1,
        "maxCount": 5
      }
    ]
  }
]
```

Monster placement is defined in the level's `encounters` section:

```json
"encounters": [
  {
    "monsterId": "skeleton",
    "count": 2,
    "minLevel": 1,
    "maxLevel": 1,
    "roomId": "main_hall"
  },
  {
    "monsterId": "zombie",
    "count": 1,
    "minLevel": 1,
    "maxLevel": 2,
    "position": {
      "x": 15,
      "y": 5
    }
  }
]
```

### Items

Items are defined with their properties and effects:

```json
"items": [
  {
    "id": "health_potion",
    "name": "Health Potion",
    "description": "A potion that restores health.",
    "symbol": "!",
    "color": "#ff0000",
    "type": "consumable",
    "value": 10,
    "effects": [
      {
        "type": "heal",
        "value": 5,
        "duration": 0
      }
    ]
  }
]
```

Item placement is defined in the level's `items` section:

```json
"items": [
  {
    "itemId": "gold",
    "roomId": "entrance",
    "chance": 0.8
  },
  {
    "itemId": "health_potion",
    "position": {
      "x": 12,
      "y": 2
    },
    "chance": 1.0
  }
]
```

### Events

Events are special occurrences that can be triggered during gameplay:

```json
"events": [
  {
    "id": "entrance_event",
    "name": "Entrance Event",
    "description": "An event that triggers when the player enters the dungeon.",
    "trigger": "level_start",
    "actions": [
      {
        "type": "message",
        "target": "player",
        "value": "You enter the forgotten crypt. The air is stale and cold."
      },
      {
        "type": "sound",
        "target": "global",
        "value": "door_creak"
      }
    ]
  }
]
```

## Development

### Project Structure

The project follows a standard Go project layout:

```
cryptcrawl/
├── bin/                  # Compiled binaries
├── cmd/                  # Application entry points
│   └── cryptcrawl/       # Main application
├── internal/             # Private application code
│   └── dungeon/          # Dungeon generation and loading
├── dungeons/             # User-created dungeon definitions
├── .ssh/                 # SSH keys
├── .vscode/              # VS Code configuration
├── Makefile              # Build automation
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # This file
```

### Building from Source

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Clean build artifacts
make clean
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...
```

### Environment Variables

- `HOST`: The host to bind to (default: localhost)
- `PORT`: The port to listen on (default: 23234)
- `DEBUG`: Enable debug mode (default: false)
- `DUNGEON_DIR`: Directory containing dungeon definitions (default: dungeons)

## License

MIT
