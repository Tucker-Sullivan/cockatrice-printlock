package printlock

import (
	"fmt"

	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/deck"
)

func ApplyPrintings(db *cardsdb.CardsDB, d *deck.Deck, opt ApplyOptions) (*deck.Deck, []error) {
	applyErrors := make([]error, 0)
	for zi := range d.Zones {
		zone := &d.Zones[zi]
		for ci := range zone.Cards {
			card := &zone.Cards[ci]
			if card.UUID != "" {
				// validate UUID
				printing, found := db.FindPrintingByUUID(card.Name, card.UUID)
				if !found {
					applyErrors = append(applyErrors, ApplyError{
						Zone:   zone.Name,
						Card:   card.Name,
						Reason: fmt.Sprintf("invalid UUID - %v", card.UUID),
					})
					continue
				}
				updateCard(card, printing)
				continue
			} else if card.SetShortName != "" && card.CollectorNumber != "" {
				// check for specific set/collectorNumber match printing
				printing, found := db.FindPrinting(card.Name, card.SetShortName, card.CollectorNumber)
				if !found {
					applyErrors = append(applyErrors, ApplyError{
						Zone:   zone.Name,
						Card:   card.Name,
						Reason: fmt.Sprintf("invalid set/collectorNumber %v - %v", card.SetShortName, card.CollectorNumber),
					})
					continue
				}
				updateCard(card, printing)
				continue
			} else if card.SetShortName != "" {
				// check for specific set and match first printing found
				printings := db.FindPrintingsInSet(card.Name, card.SetShortName)
				if len(printings) > 0 {
					updateCard(card, printings[0])
					continue
				}
			}
			// check for printing in order of set priority
			printing, found := db.FindPrintingBySetPriority(card.Name, opt.SetPriority)
			if !found {
				applyErrors = append(applyErrors, ApplyError{
					Zone:   zone.Name,
					Card:   card.Name,
					Reason: "could not find printing",
				})
				continue
			}
			updateCard(card, printing)
		}
	}
	return d, applyErrors
}

func updateCard(card *deck.Card, printing cardsdb.Printing) {
	card.UUID = printing.UUID
	card.CollectorNumber = printing.CollectorNumber
	card.SetShortName = printing.SetCode
}
