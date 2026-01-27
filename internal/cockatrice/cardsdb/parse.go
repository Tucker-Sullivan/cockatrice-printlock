package cardsdb

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
)

type xmlSetEntry struct {
	Name     string `xml:"name"`
	LongName string `xml:"longname"`
}

type xmlCard struct {
	Name string        `xml:"name"`
	Sets []xmlPrinting `xml:"set"`
}

type xmlPrinting struct {
	UUID            string `xml:"uuid,attr"`
	CollectorNumber string `xml:"num,attr"`
	MultiverseID    string `xml:"muid,attr"`
	SetCode         string `xml:",chardata"`
}

func ParseCardsDB(path string) (*CardsDB, error) {
	cardsDB := NewCardsDB()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}
	defer cleanup()

	decoder := xml.NewDecoder(file)
	inSets, inCards := false, false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "sets":
				inSets = true
			case "cards":
				inCards = true
			case "set":
				if !inSets {
					continue
				}
				var item xmlSetEntry
				if err := decoder.DecodeElement(&item, &t); err != nil {
					return nil, err
				}
				if _, ok := cardsDB.SetsByCode[item.Name]; !ok {
					cardsDB.SetsByCode[item.Name] = SetEntry{
						Code:     strings.TrimSpace(item.Name),
						LongName: strings.TrimSpace(item.LongName),
					}
				}
			case "card":
				if !inCards {
					continue
				}
				var item xmlCard
				if err := decoder.DecodeElement(&item, &t); err != nil {
					return nil, err
				}
				if _, ok := cardsDB.CardsByName[item.Name]; !ok {
					printings := make([]Printing, 0, len(item.Sets))
					for _, xmlPrint := range item.Sets {
						printings = append(printings, Printing{
							UUID:            strings.TrimSpace(xmlPrint.UUID),
							SetCode:         strings.TrimSpace(xmlPrint.SetCode),
							MultiverseID:    strings.TrimSpace(xmlPrint.MultiverseID),
							CollectorNumber: strings.TrimSpace(xmlPrint.CollectorNumber),
						})
					}
					cardsDB.CardsByName[item.Name] = CardEntry{
						Name:      strings.TrimSpace(item.Name),
						Printings: printings,
					}
				}
			}
		// end of StartElement
		case xml.EndElement:
			switch t.Name.Local {
			case "sets":
				inSets = false
			case "cards":
				inCards = false
			}
		// end of EndElement
		default:
		}
	}

	return cardsDB, nil
}
