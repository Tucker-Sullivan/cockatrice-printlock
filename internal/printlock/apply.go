package printlock

import (
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/cardsdb"
	"github.com/Tucker-Sullivan/cockatrice-printlock/internal/cockatrice/deck"
)

func ApplyPrintings(db *cardsdb.CardsDB, d *deck.Deck, opt ApplyOptions) (*deck.Deck, []error) {
	applyErrors := make([]error, 0)
	for zi := range d.Zones {
		zone := &d.Zones[zi]
		for ci := range zone.Cards {
			card := &zone.Cards[ci]
			// check for printing in order of set priority
			printing, found := db.FindPrintingBySetPriority(card.Name, opt.SetPriority, opt.PreferHigherCollectionNumber)
			if !found {
				applyErrors = append(applyErrors, ApplyError{
					Zone:   zone.Name,
					Card:   card.Name,
					Reason: "could not find printing",
				})
				continue
			}
			card.UUID = printing.UUID
			card.CollectorNumber = printing.CollectorNumber
			card.SetShortName = printing.SetCode
		}
	}
	return d, applyErrors
}
