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

		// load cfg
		cfg, err := config.LoadConfig()
		if err != nil {
			return errors.New("error loading config. make sure to run init command if config file has not been generated yet")
		}

		// load cardsdb
		db, err := cardsdb.ParseCardsDB(cfg.CardsXMLPath)
		if err != nil {
			return err
		}

		// resolve deck path
		deckFilePath, err := deck.ResolveDeckPath(cfg.DeckDir, deckName)
		if err != nil {
			return err
		}

		// load deck
		d, err := deck.LoadDeck(deckFilePath)
		if err != nil {
			return err
		}

		// apply printings
		cliSetPriority := parseCSV(applySetsCSV)
		completedSetPriority := make([]string, 0, len(cliSetPriority)+len(cfg.GlobalSetPriority))

		if cfg.GlobalSetsHavePriority {
			completedSetPriority = append(completedSetPriority, cfg.GlobalSetPriority...)
			completedSetPriority = append(completedSetPriority, cliSetPriority...)
		} else {
			completedSetPriority = append(completedSetPriority, cliSetPriority...)
			completedSetPriority = append(completedSetPriority, cfg.GlobalSetPriority...)
		}
		opt := printlock.ApplyOptions{
			PreferHigherCollectionNumber: cfg.PreferHigherCollectionNumber,
			SetPriority:                  completedSetPriority,
		}
		d, errs := printlock.ApplyPrintings(db, d, opt)
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}

		// write deck
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
