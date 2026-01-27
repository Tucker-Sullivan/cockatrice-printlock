package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

		cardsPath, err := promptPath(reader, "Path to Cockatrice cards.xml")
		if err != nil {
			return err
		}
		if err := validateFile(cardsPath); err != nil {
			return fmt.Errorf("cards.xml path invalid: %w", err)
		}

		deckDir, err := promptPath(reader, "Path to Cockatrice decks directory")
		if err != nil {
			return err
		}
		if err := validateDir(deckDir); err != nil {
			return fmt.Errorf("decks directory invalid: %w", err)
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

func validateFile(p string) error {
	p = strings.TrimSpace(p)
	if p == "" {
		return errors.New("empty path")
	}

	info, err := os.Stat(p)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("expected a file, got a directory")
	}
	return nil
}

func validateDir(p string) error {
	p = strings.TrimSpace(p)
	if p == "" {
		return errors.New("empty path")
	}

	info, err := os.Stat(p)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("expected a directory")
	}
	// Normalize (optional): ensure we store an absolute path
	abs, err := filepath.Abs(p)
	if err == nil {
		_ = abs
	}
	return nil
}
