package cardsdb

import (
	"strings"
)

func (c *CardsDB) FindCard(name string) (CardEntry, bool) {
	entry, ok := c.CardsByName[strings.TrimSpace(name)]
	return entry, ok
}

func (c *CardsDB) FindPrintingsInSet(cardName, setCode string) []Printing {
	printings := make([]Printing, 0)
	tCardName := strings.TrimSpace(cardName)
	tSetCode := strings.TrimSpace(setCode)
	card, ok := c.CardsByName[tCardName]
	if !ok {
		return printings
	}
	for _, printing := range card.Printings {
		if printing.SetCode == tSetCode {
			printings = append(printings, printing)
		}
	}

	return printings
}

func (c *CardsDB) FindPrinting(cardName, setCode, collectionNumber string) (Printing, bool) {
	printings := c.FindPrintingsInSet(cardName, setCode)
	if len(printings) == 0 {
		return Printing{}, false
	}

	tCollectionNumber := strings.TrimSpace(collectionNumber)
	if tCollectionNumber == "" {
		return printings[0], true
	}

	for _, printing := range printings {
		if printing.CollectorNumber == tCollectionNumber {
			return printing, true
		}
	}

	return Printing{}, false
}

func (c *CardsDB) FindPrintingByUUID(cardName, uuid string) (Printing, bool) {
	card, exists := c.FindCard(cardName)
	if !exists {
		return Printing{}, false
	}

	tUUID := strings.TrimSpace(uuid)
	if tUUID == "" {
		return Printing{}, false
	}

	for _, printing := range card.Printings {
		if printing.UUID == tUUID {
			return printing, true
		}
	}

	return Printing{}, false
}

func (c *CardsDB) FindPrintingBySetPriority(cardName string, setPriority []string) (Printing, bool) {
	if len(setPriority) == 0 {
		return Printing{}, false
	}

	tCardName := strings.TrimSpace(cardName)
	if tCardName == "" {
		return Printing{}, false
	}

	for _, setCode := range setPriority {
		tSetCode := strings.TrimSpace(setCode)
		printings := c.FindPrintingsInSet(tCardName, tSetCode)
		if len(printings) == 0 {
			continue
		}
		return printings[0], true
	}

	return Printing{}, false
}
