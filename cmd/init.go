package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize printlock configuration",
	Long:  "Interactively sets the Cockatrice cards.xml path and the decks directory, then saves config.json to your user config dir.",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Printlock init")
		fmt.Println("-------------")

		cardsPath, deckDir, err := searchDefaultDirectory()
		if err != nil {
			return err
		}
		if cardsPath == "" {
			cardsPath, err := promptPath(reader, "Path to Cockatrice cards.xml")
			if err != nil {
				return err
			}
			cardsPath, err = validateFile(cardsPath)
			if err != nil {
				return err
			}
		}
		if deckDir == "" {
			deckDir, err := promptPath(reader, "Path to Cockatrice decks directory")
			if err != nil {
				return err
			}
			deckDir, err = validateDir(deckDir)
			if err != nil {
				return err
			}
		}

		cfg := config.Config{
			CardsXMLPath:      cardsPath,
			DeckDir:           deckDir,
			GlobalSetPriority: []string{},
		}

		if err := config.SaveConfig(cfg); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("✔ Saved configuration")
		fmt.Printf("  cards.xml: %s\n", cardsPath)
		fmt.Printf("  decks dir: %s\n", deckDir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.SilenceUsage = true
}

func promptPath(r *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s: ", label)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func validateFile(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", errors.New("empty path")
	}

	info, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("expected a file, got a directory")
	}
	// Normalize (optional): ensure we store an absolute path
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func validateDir(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", errors.New("empty path")
	}

	info, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("expected a directory")
	}
	// Normalize (optional): ensure we store an absolute path
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func searchDefaultDirectory() (string, string, error) {
	cardsPath, deckDir := "", ""
	var err error
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			println("could not find LOCALAPPDATA")
			return "", "", nil
		}
		cardsPath, deckDir, err = searchCockatriceDirectory(filepath.Join(localAppData, "Cockatrice", "Cockatrice"))
		if err != nil {
			return "", "", err
		}
	case "darwin":
		applicationSupport, err := os.UserConfigDir()
		if err != nil {
			return "", "", err
		}
		if applicationSupport == "" {
			println("could not find Application Support")
			return "", "", nil
		}
		cardsPath, deckDir, err = searchCockatriceDirectory(filepath.Join(applicationSupport, "Cockatrice", "Cockatrice"))
		if err != nil {
			return "", "", err
		}
	}

	return cardsPath, deckDir, nil
}

func searchCockatriceDirectory(rootPath string) (string, string, error) {
	cardsPath, deckDir := "", ""
	stat, err := os.Stat(rootPath)
	if err != nil {
		return "", "", err
	}
	if !stat.IsDir() {
		return "", "", errors.New("could not find default Cockatrice directory")
	}
	cardsXmlPath := filepath.Join(rootPath, "cards.xml")
	if _, err := os.Stat(cardsXmlPath); err != nil {
		return "", "", nil
	}
	cardsPath = cardsXmlPath
	deckDirPath := filepath.Join(rootPath, "decks")
	deckDirStat, err := os.Stat(deckDirPath)
	if err != nil {
		return "", "", err
	}
	if !deckDirStat.IsDir() {
		return "", "", errors.New("could not find default decks directory")
	}
	entries, err := os.ReadDir(deckDirPath)
	if err != nil {
		return "", "", err
	}
	if len(entries) == 0 {
		return "", "", errors.New("empty default decks directory")
	}
	deckDir = deckDirPath
	return cardsPath, deckDir, nil
}
