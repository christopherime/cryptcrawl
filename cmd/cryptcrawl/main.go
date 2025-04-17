package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"

	"cryptcrawl/internal/dungeon"
)

// Default configuration
const (
	defaultHost = "localhost"
	defaultPort = 23234
)

// Global dungeon loader
var dungeonLoader *dungeon.DungeonLoader

func main() {
	// Get configuration from environment variables
	host := getEnv("HOST", defaultHost)
	portStr := getEnv("PORT", fmt.Sprintf("%d", defaultPort))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid PORT value: %s, using default: %d", portStr, defaultPort)
		port = defaultPort
	}

	// Setup debug logging if enabled
	debugMode := getEnv("DEBUG", "false") == "true"
	if debugMode {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Println("Failed to create debug log file:", err)
		} else {
			defer f.Close()
			log.Println("Debug logging enabled to debug.log")
		}
	}

	// Initialize dungeon loader
	dungeonDir := getEnv("DUNGEON_DIR", "dungeons")
	if err := os.MkdirAll(dungeonDir, 0755); err != nil {
		log.Fatalf("Failed to create dungeons directory: %v", err)
	}

	loader, err := dungeon.NewDungeonLoader(dungeonDir)
	if err != nil {
		log.Fatalf("Failed to initialize dungeon loader: %v", err)
	}
	dungeonLoader = loader

	// Create SSH server
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/cryptcrawl_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && err != ssh.ErrServerClosed {
		log.Fatalln(err)
	}
}

// teaHandler creates a new BubbleTea program for each SSH session
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	_, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, bye!")
		return nil, nil
	}

	m := initialModel()
	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
