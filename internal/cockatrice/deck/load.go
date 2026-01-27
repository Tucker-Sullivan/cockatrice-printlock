package deck

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type xmlCard struct {
	Number          int    `xml:"number,attr"`
	Name            string `xml:"name,attr"`
	SetShortName    string `xml:"setShortName,attr"`
	CollectorNumber string `xml:"collectorNumber,attr"`
	UUID            string `xml:"uuid,attr"`
}

func LoadDeck(deckFilePath string) (*Deck, error) {
	tDeckFilePath := strings.TrimSpace(deckFilePath)
	if tDeckFilePath == "" {
		return nil, errors.New("invalid deck file path")
	}

	file, err := os.Open(tDeckFilePath)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}
	defer cleanup()

	deck := NewDeck()
	var currentZone *Zone
	decoder := xml.NewDecoder(file)
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
			case "cockatrice_deck":
				v := strings.TrimSpace(getAttr(t, "version"))
				if v != "" {
					deck.Version = v
				}
			case "zone":
				zoneName := getAttr(t, "name")
				deck.Zones = append(deck.Zones, Zone{Name: strings.TrimSpace(zoneName)})
				currentZone = &deck.Zones[len(deck.Zones)-1]
			case "card":
				if currentZone == nil {
					return nil, errors.New("card outside of zone")
				}
				var item xmlCard
				if err := decoder.DecodeElement(&item, &t); err != nil {
					return nil, err
				}
				currentZone.Cards = append(currentZone.Cards, Card{
					Count:           item.Number,
					Name:            strings.TrimSpace(item.Name),
					SetShortName:    strings.TrimSpace(item.SetShortName),
					CollectorNumber: strings.TrimSpace(item.CollectorNumber),
					UUID:            strings.TrimSpace(item.UUID),
				})
			default:
				// Any other top-level element (deckname, tags, comments, bannerCard, etc.)
				// Capture and preserve it so the writer has it
				if currentZone == nil {
					raw, err := captureElement(decoder, t)
					if err != nil {
						return nil, err
					}
					deck.Header = append(deck.Header, RawNode{
						Name: t.Name.Local,
						XML:  raw,
					})
				}
			}

		case xml.EndElement:
			if t.Name.Local == "zone" {
				currentZone = nil
			}
		}
	}
	return deck, nil
}

func ResolveDeckPath(deckDir, deckName string) (string, error) {
	tDeckDir := strings.TrimSpace(deckDir)
	if tDeckDir == "" {
		return "", errors.New("invalid deck directory")
	}

	tDeckName := strings.TrimSpace(deckName)
	if tDeckName == "" {
		return "", errors.New("invalid deck name")
	}

	if filepath.IsAbs(tDeckName) {
		return "", errors.New("deck path must be relative to DeckDir")
	}

	if filepath.Ext(tDeckName) != ".cod" {
		tDeckName += ".cod"
	}

	deckFilePath, err := filepath.Abs(filepath.Join(tDeckDir, tDeckName))
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(tDeckDir, deckFilePath)
	if err != nil {
		return "", err
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", errors.New("deck path must be relative to DeckDir")
	}

	foundFiles := make([]string, 0)
	walk := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == tDeckName {
			foundFiles = append(foundFiles, strings.Replace(path, tDeckDir, "", 1))
		}
		return nil
	}

	if err := filepath.WalkDir(deckDir, walk); err != nil {
		return "", err
	}

	if len(foundFiles) == 0 {
		return "", errors.New("no decks found with the given name")
	} else if len(foundFiles) > 1 {
		errorString := "Multiple files found: try specifying a path to the file instead"
		for _, entry := range foundFiles {
			errorString += fmt.Sprintf("\n%v", entry)
		}
		return "", errors.New(errorString)
	} else {
		deckFilePath = filepath.Join(deckDir, foundFiles[0])
	}

	return deckFilePath, nil
}

func getAttr(element xml.StartElement, key string) string {
	value := ""
	for _, a := range element.Attr {
		if a.Name.Local == key {
			value = a.Value
			break
		}
	}
	return value
}

func captureElement(decoder *xml.Decoder, start xml.StartElement) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	depth := 1

	if err := enc.EncodeToken(start); err != nil {
		return nil, err
	}

	for depth > 0 {
		tok, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}

		if err := enc.EncodeToken(tok); err != nil {
			return nil, err
		}
	}

	if err := enc.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
