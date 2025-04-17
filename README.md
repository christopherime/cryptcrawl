# CryptCrawl

A MUD-style dungeon crawler accessible via SSH, built with Go and the charm.sh libraries.

## Description

CryptCrawl is a terminal-based dungeon crawler game that users can connect to via SSH. It features:

- Procedurally generated dungeons
- Turn-based combat
- Multiple levels to explore
- Gold collection
- Monsters that get tougher as you progress

## Technologies Used

- Go programming language
- [charm.sh](https://charm.sh/) libraries:
  - [Wish](https://github.com/charmbracelet/wish) - SSH server
  - [BubbleTea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
  - [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
  - [LipGloss](https://github.com/charmbracelet/lipgloss) - Styling

## How to Play

1. Start the server:

   ```bash
   go run .
   ```

2. Connect via SSH:

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

## License

MIT
