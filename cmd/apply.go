package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/deck"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/config"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/printlock"
	"github.com/spf13/cobra"
)

var (
	applyDeckName string
	applySetsCSV  string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply deterministic printing UUIDs to a Cockatrice deck file",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := strings.TrimSpace(applyDeckName)
		if deckName == "" {
			return errors.New("--deck is required")
		}

		setPriority := parseCSV(applySetsCSV)

		// 1) load cfg
		cfg, err := config.LoadConfig()
		if err != nil {
			return errors.New("error loading config. make sure to run init command if config file has not been generated yet")
		}

		// 2) resolve effective setPriority (cli override vs config)
		if len(setPriority) == 0 {
			setPriority = cfg.GlobalSetPriority
		}

		// 3) load cardsdb (cards.xml)
		db, err := cardsdb.ParseCardsDB(cfg.CardsXMLPath)
		if err != nil {
			return err
		}

		// 4) resolve deck path
		deckFilePath, err := deck.ResolveDeckPath(cfg.DeckDir, deckName)
		if err != nil {
			return err
		}

		// 5) load deck
		d, err := deck.LoadDeck(deckFilePath)
		if err != nil {
			return err
		}

		// 6) apply printings (P1)
		opt := printlock.ApplyOptions{
			SetPriority: setPriority,
		}
		d, errs := printlock.ApplyPrintings(db, d, opt)
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}

		// 7) write deck (D3)
		err = d.WriteDeckFile(deckFilePath)
		if err != nil {
			return err
		}

		if len(errs) > 0 {
			return errors.New("apply completed with issues")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVar(&applyDeckName, "deck", "", "Deck name (without .cod extension)")
	applyCmd.Flags().StringVar(&applySetsCSV, "sets", "", "Comma-separated set priority override (e.g. DRC,DFT)")

	// Make --deck required in Cobra’s help/validation (still keep our check for older cobra behavior).
	_ = applyCmd.MarkFlagRequired("deck")
}
