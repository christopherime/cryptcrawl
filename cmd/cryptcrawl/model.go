package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define key mappings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Attack key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Attack},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "w", "k"),
		key.WithHelp("↑/w/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "s", "j"),
		key.WithHelp("↓/s/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "a", "h"),
		key.WithHelp("←/a/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d", "l"),
		key.WithHelp("→/d/l", "move right"),
	),
	Attack: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "attack"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

// Legacy tile symbols for backward compatibility
const (
	TileEmpty   = ' '
	TileWall    = '#'
	TilePlayer  = '@'
	TileMonster = 'M'
	TileGold    = '$'
	TileExit    = 'E'
)

// Position represents a 2D position
type Position struct {
	X, Y int
}

// Entity represents a game entity
type Entity struct {
	Pos       Position
	Symbol    rune
	Health    int
	MaxHealth int
	Damage    int
	Name      string
}

// Model represents the game state
type model struct {
	dungeon   [][]TileType // Using TileType instead of rune
	player    Entity
	monsters  []Entity
	width     int
	height    int
	viewport  viewport.Model
	help      help.Model
	keys      keyMap
	showHelp  bool
	messages  []string
	gold      int
	level     int
	gameOver  bool
	gameWon   bool
	revealMap bool // Debug option to reveal the entire map
}

// Initialize the model
func initialModel() model {
	// Set up the help model
	h := help.New()
	h.ShowAll = false

	// Check if debug mode is enabled
	debugMode := os.Getenv("DEBUG") == "true"

	// Create a new model
	m := model{
		width:     97, // Width based on the provided string length
		height:    30, // Increased height for better visibility
		help:      h,
		keys:      keys,
		showHelp:  false,
		messages:  []string{"Welcome to CryptCrawl! Use arrow keys to move."},
		gold:      0,
		level:     1,
		gameOver:  false,
		gameWon:   false,
		revealMap: debugMode, // Reveal the entire map in debug mode
	}

	// Try to load a dungeon from the dungeon loader
	if dungeonLoader != nil {
		def := dungeonLoader.GetCurrentDungeon()
		if def != nil {
			m.addMessage(fmt.Sprintf("Loaded dungeon: %s", def.Name))
			m.addMessage(def.Description)

			// Generate the dungeon from the definition
			grid, _, err := dungeonLoader.GenerateCurrentLevel(m.level - 1)
			if err == nil && grid != nil {
				// Convert the rune grid to TileType grid
				m.dungeon = make([][]TileType, len(grid))
				for y := range grid {
					m.dungeon[y] = make([]TileType, len(grid[y]))
					for x, r := range grid[y] {
						switch r {
						case '#':
							m.dungeon[y][x] = Wall
						case '@':
							m.dungeon[y][x] = Player
							// Set player position
							m.player.Pos.X = x
							m.player.Pos.Y = y
						case 'E':
							m.dungeon[y][x] = Exit
						case '$':
							m.dungeon[y][x] = Gold
						case 'M', 'S', 'Z', 'W':
							m.dungeon[y][x] = Monster
							// Add monster
							monster := Entity{
								Pos:       Position{X: x, Y: y},
								Symbol:    r,
								Health:    5,
								MaxHealth: 5,
								Damage:    2,
								Name:      "Monster",
							}
							m.monsters = append(m.monsters, monster)
						case '?':
							m.dungeon[y][x] = Chest
						case '^':
							m.dungeon[y][x] = Trap
						case '+':
							m.dungeon[y][x] = Door
						case '~':
							// Could be water or lava
							if rand.Intn(2) == 0 {
								m.dungeon[y][x] = Water
							} else {
								m.dungeon[y][x] = Lava
							}
						default:
							m.dungeon[y][x] = Empty
						}
					}
				}
				return m
			}
		}
	}

	// Generate the dungeon
	m.generateDungeon()

	// Set up the viewport
	vp := viewport.New(m.width, m.height-5) // Leave room for messages and status
	vp.SetContent(m.dungeonToString())
	m.viewport = vp

	return m
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
		case key.Matches(msg, m.keys.Up):
			if !m.gameOver && !m.gameWon {
				m.movePlayer(0, -1)
			}
		case key.Matches(msg, m.keys.Down):
			if !m.gameOver && !m.gameWon {
				m.movePlayer(0, 1)
			}
		case key.Matches(msg, m.keys.Left):
			if !m.gameOver && !m.gameWon {
				m.movePlayer(-1, 0)
			}
		case key.Matches(msg, m.keys.Right):
			if !m.gameOver && !m.gameWon {
				m.movePlayer(1, 0)
			}
		case key.Matches(msg, m.keys.Attack):
			if !m.gameOver && !m.gameWon {
				m.attackNearbyMonsters()
			}
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 5 // Leave room for messages and status
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update viewport
	m.viewport.SetContent(m.dungeonToString())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m model) View() string {
	if m.gameOver {
		return fmt.Sprintf("\n\n  GAME OVER\n\n  You reached level %d and collected %d gold.\n\n  Press q to quit.", m.level, m.gold)
	}

	if m.gameWon {
		return fmt.Sprintf("\n\n  VICTORY!\n\n  You escaped the dungeon with %d gold!\n\n  Press q to quit.", m.gold)
	}

	// Render the dungeon
	dungeonView := m.viewport.View()

	// Render the status bar
	healthStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	goldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))
	levelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))

	healthBar := fmt.Sprintf("❤️ %s%d/%d", healthStyle.Render(""), m.player.Health, m.player.MaxHealth)
	goldBar := fmt.Sprintf("💰 %s%d", goldStyle.Render(""), m.gold)
	levelBar := fmt.Sprintf("📜 %sLevel %d", levelStyle.Render(""), m.level)

	statusBar := fmt.Sprintf("%s | %s | %s", healthBar, goldBar, levelBar)

	// Render the message log (last 3 messages)
	messageLog := ""
	startIdx := 0
	if len(m.messages) > 3 {
		startIdx = len(m.messages) - 3
	}
	for i := startIdx; i < len(m.messages); i++ {
		messageLog += fmt.Sprintf("  %s\n", m.messages[i])
	}

	// Render help if needed
	helpView := ""
	if m.showHelp {
		helpView = "\n" + m.help.View(m.keys)
	}

	// Combine all views
	return fmt.Sprintf("%s\n\n%s\n%s%s", dungeonView, statusBar, messageLog, helpView)
}

// Generate a random dungeon
func (m *model) generateDungeon() {
	// No need to seed in Go 1.20+
	// rand.Seed is deprecated

	// Create an empty dungeon filled with walls
	m.dungeon = make([][]TileType, m.height)
	for i := range m.dungeon {
		m.dungeon[i] = make([]TileType, m.width)
		for j := range m.dungeon[i] {
			m.dungeon[i][j] = Wall
		}
	}

	// Create rooms
	numRooms := rand.Intn(5) + 5 // 5-10 rooms
	rooms := make([]struct{ x, y, w, h int }, 0, numRooms)

	for i := 0; i < numRooms; i++ {
		roomW := rand.Intn(8) + 5 // 5-12 width
		roomH := rand.Intn(5) + 3 // 3-7 height
		roomX := rand.Intn(m.width-roomW-2) + 1
		roomY := rand.Intn(m.height-roomH-2) + 1

		// Check for overlap with existing rooms
		overlap := false
		for _, r := range rooms {
			if roomX <= r.x+r.w+1 && roomX+roomW+1 >= r.x &&
				roomY <= r.y+r.h+1 && roomY+roomH+1 >= r.y {
				overlap = true
				break
			}
		}

		if !overlap {
			// Carve out the room
			for y := roomY; y < roomY+roomH; y++ {
				for x := roomX; x < roomX+roomW; x++ {
					m.dungeon[y][x] = Empty
				}
			}

			// Add the room to our list
			rooms = append(rooms, struct{ x, y, w, h int }{roomX, roomY, roomW, roomH})
		}
	}

	// Connect rooms with corridors
	for i := 0; i < len(rooms)-1; i++ {
		startX := rooms[i].x + rooms[i].w/2
		startY := rooms[i].y + rooms[i].h/2
		endX := rooms[i+1].x + rooms[i+1].w/2
		endY := rooms[i+1].y + rooms[i+1].h/2

		// Horizontal corridor
		for x := min(startX, endX); x <= max(startX, endX); x++ {
			m.dungeon[startY][x] = Empty
		}

		// Vertical corridor
		for y := min(startY, endY); y <= max(startY, endY); y++ {
			m.dungeon[y][endX] = Empty
		}
	}

	// Place player in the first room
	playerX := rooms[0].x + rooms[0].w/2
	playerY := rooms[0].y + rooms[0].h/2
	m.player = Entity{
		Pos:       Position{X: playerX, Y: playerY},
		Symbol:    TilePlayer,
		Health:    10,
		MaxHealth: 10,
		Damage:    2,
		Name:      "Player",
	}
	m.dungeon[playerY][playerX] = Player

	// Place exit in the last room
	exitX := rooms[len(rooms)-1].x + rooms[len(rooms)-1].w/2
	exitY := rooms[len(rooms)-1].y + rooms[len(rooms)-1].h/2
	m.dungeon[exitY][exitX] = Exit

	// Place monsters and gold
	m.monsters = []Entity{}
	for i := 1; i < len(rooms)-1; i++ {
		// Add 1-3 monsters per room
		numMonsters := rand.Intn(3) + 1
		for j := 0; j < numMonsters; j++ {
			monsterX := rooms[i].x + rand.Intn(rooms[i].w)
			monsterY := rooms[i].y + rand.Intn(rooms[i].h)

			// Make sure the position is empty
			if m.dungeon[monsterY][monsterX] == Empty {
				monster := Entity{
					Pos:       Position{X: monsterX, Y: monsterY},
					Symbol:    TileMonster,
					Health:    3 + m.level,
					MaxHealth: 3 + m.level,
					Damage:    1 + m.level/2,
					Name:      "Monster",
				}
				m.monsters = append(m.monsters, monster)
				m.dungeon[monsterY][monsterX] = Monster
			}
		}

		// Add 1-5 gold piles per room
		numGold := rand.Intn(5) + 1
		for j := 0; j < numGold; j++ {
			goldX := rooms[i].x + rand.Intn(rooms[i].w)
			goldY := rooms[i].y + rand.Intn(rooms[i].h)

			// Make sure the position is empty
			if m.dungeon[goldY][goldX] == Empty {
				m.dungeon[goldY][goldX] = Gold
			}
		}
	}
}

// Convert the dungeon to a string for display
func (m model) dungeonToString() string {
	var result string
	for y := 0; y < len(m.dungeon); y++ {
		for x := 0; x < len(m.dungeon[y]); x++ {
			// If the tile is not visible to the player and revealMap is false, show a blank space
			if !m.revealMap && !m.isVisible(x, y) {
				result += " "
			} else {
				result += RenderTile(m.dungeon[y][x])
			}
		}
		result += "\n"
	}
	return result
}

// isVisible determines if a tile is visible to the player
func (m model) isVisible(x, y int) bool {
	// Simple visibility: if it's within 5 tiles of the player, it's visible
	dx := abs(x - m.player.Pos.X)
	dy := abs(y - m.player.Pos.Y)
	return dx <= 5 && dy <= 5
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Move the player in the given direction
func (m *model) movePlayer(dx, dy int) {
	newX := m.player.Pos.X + dx
	newY := m.player.Pos.Y + dy

	// Check if the new position is valid
	if newX < 0 || newX >= m.width || newY < 0 || newY >= m.height {
		return
	}

	// Check what's at the new position
	switch m.dungeon[newY][newX] {
	case Wall, Water, Lava:
		// Can't move through walls or hazards
		return
	case Monster:
		// Attack the monster
		for i, monster := range m.monsters {
			if monster.Pos.X == newX && monster.Pos.Y == newY {
				// Player attacks monster
				damage := m.player.Damage
				m.monsters[i].Health -= damage
				m.addMessage(fmt.Sprintf("You hit the monster for %d damage!", damage))

				// Check if monster is dead
				if m.monsters[i].Health <= 0 {
					m.addMessage("You killed the monster!")
					m.dungeon[newY][newX] = Empty
					// Remove the monster from the list
					m.monsters = append(m.monsters[:i], m.monsters[i+1:]...)
				} else {
					// Monster attacks back
					damage := m.monsters[i].Damage
					m.player.Health -= damage
					m.addMessage(fmt.Sprintf("The monster hits you for %d damage!", damage))

					// Check if player is dead
					if m.player.Health <= 0 {
						m.gameOver = true
						m.addMessage("You died!")
					}
				}
				break
			}
		}
	case Gold:
		// Collect gold
		goldAmount := rand.Intn(10) + 1
		m.gold += goldAmount
		m.addMessage(fmt.Sprintf("You found %d gold!", goldAmount))
		m.dungeon[newY][newX] = Empty
		// Move player
		m.dungeon[m.player.Pos.Y][m.player.Pos.X] = Empty
		m.player.Pos.X = newX
		m.player.Pos.Y = newY
		m.dungeon[newY][newX] = Player
	case Chest:
		// Open chest
		itemType := rand.Intn(3)
		switch itemType {
		case 0: // Gold
			goldAmount := rand.Intn(20) + 10
			m.gold += goldAmount
			m.addMessage(fmt.Sprintf("You found %d gold in the chest!", goldAmount))
		case 1: // Health potion
			healthAmount := rand.Intn(5) + 3
			m.player.Health = min(m.player.Health+healthAmount, m.player.MaxHealth)
			m.addMessage(fmt.Sprintf("You found a health potion! +%d HP", healthAmount))
		case 2: // Damage boost
			m.player.Damage++
			m.addMessage("You found a weapon upgrade! +1 damage")
		}
		m.dungeon[newY][newX] = Empty
		// Move player
		m.dungeon[m.player.Pos.Y][m.player.Pos.X] = Empty
		m.player.Pos.X = newX
		m.player.Pos.Y = newY
		m.dungeon[newY][newX] = Player
	case Trap:
		// Trigger trap
		damage := rand.Intn(3) + 1
		m.player.Health -= damage
		m.addMessage(fmt.Sprintf("You triggered a trap! -%d HP", damage))
		// Check if player is dead
		if m.player.Health <= 0 {
			m.gameOver = true
			m.addMessage("You died!")
			return
		}
		// Move player
		m.dungeon[m.player.Pos.Y][m.player.Pos.X] = Empty
		m.player.Pos.X = newX
		m.player.Pos.Y = newY
		m.dungeon[newY][newX] = Player
	case Exit:
		// Go to next level or win the game
		if m.level < 3 {
			m.level++
			m.addMessage(fmt.Sprintf("You descend to level %d...", m.level))
			m.generateDungeon()
		} else {
			m.gameWon = true
			m.addMessage("You escaped the dungeon!")
		}
	case Empty, Door:
		// Move player
		m.dungeon[m.player.Pos.Y][m.player.Pos.X] = Empty
		m.player.Pos.X = newX
		m.player.Pos.Y = newY
		m.dungeon[newY][newX] = Player
	}

	// Move monsters after player's turn
	m.moveMonsters()
}

// Move all monsters
func (m *model) moveMonsters() {
	for i := range m.monsters {
		// Skip dead monsters
		if m.monsters[i].Health <= 0 {
			continue
		}

		// 50% chance to move
		if rand.Intn(2) == 0 {
			continue
		}

		// Get current position
		oldX := m.monsters[i].Pos.X
		oldY := m.monsters[i].Pos.Y

		// Determine direction towards player
		dx := 0
		dy := 0
		if m.monsters[i].Pos.X < m.player.Pos.X {
			dx = 1
		} else if m.monsters[i].Pos.X > m.player.Pos.X {
			dx = -1
		}
		if m.monsters[i].Pos.Y < m.player.Pos.Y {
			dy = 1
		} else if m.monsters[i].Pos.Y > m.player.Pos.Y {
			dy = -1
		}

		// Randomly choose to move in x or y direction
		if rand.Intn(2) == 0 && dx != 0 {
			dy = 0
		} else if dy != 0 {
			dx = 0
		}

		// Calculate new position
		newX := oldX + dx
		newY := oldY + dy

		// Check if the new position is valid
		if newX < 0 || newX >= m.width || newY < 0 || newY >= m.height {
			continue
		}

		// Check what's at the new position
		switch m.dungeon[newY][newX] {
		case Empty, Gold, Trap, Chest, Door:
			// Move monster
			m.dungeon[oldY][oldX] = Empty
			m.monsters[i].Pos.X = newX
			m.monsters[i].Pos.Y = newY
			m.dungeon[newY][newX] = Monster
		case Player:
			// Attack player
			damage := m.monsters[i].Damage
			m.player.Health -= damage
			m.addMessage(fmt.Sprintf("The monster hits you for %d damage!", damage))

			// Check if player is dead
			if m.player.Health <= 0 {
				m.gameOver = true
				m.addMessage("You died!")
			}
		}
	}
}

// Attack all monsters adjacent to the player
func (m *model) attackNearbyMonsters() {
	attacked := false

	// Check all adjacent positions
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			// Skip the player's position
			if dx == 0 && dy == 0 {
				continue
			}

			newX := m.player.Pos.X + dx
			newY := m.player.Pos.Y + dy

			// Check bounds
			if newX < 0 || newX >= m.width || newY < 0 || newY >= m.height {
				continue
			}

			// Check if there's a monster
			if m.dungeon[newY][newX] == Monster {
				attacked = true

				// Find the monster
				for i, monster := range m.monsters {
					if monster.Pos.X == newX && monster.Pos.Y == newY {
						// Player attacks monster
						damage := m.player.Damage
						m.monsters[i].Health -= damage
						m.addMessage(fmt.Sprintf("You hit the monster for %d damage!", damage))

						// Check if monster is dead
						if m.monsters[i].Health <= 0 {
							m.addMessage("You killed the monster!")
							m.dungeon[newY][newX] = Empty
							// Remove the monster from the list
							m.monsters = append(m.monsters[:i], m.monsters[i+1:]...)
						}
						break
					}
				}
			}
		}
	}

	if !attacked {
		m.addMessage("You swing at the air!")
	}
}

// Add a message to the message log
func (m *model) addMessage(msg string) {
	m.messages = append(m.messages, msg)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
