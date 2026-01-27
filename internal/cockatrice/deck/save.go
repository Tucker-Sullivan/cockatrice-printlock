package deck

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func (d *Deck) WriteDeckFile(deckFilePath string) (err error) {
	if d == nil {
		return errors.New("nil deck")
	}

	tDeckFilePath := strings.TrimSpace(deckFilePath)
	if tDeckFilePath == "" {
		return errors.New("invalid path")
	}

	tmpPath := tDeckFilePath + ".tmp"

	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	closed, success := false, false
	defer func() {
		if !closed {
			if cerr := file.Close(); err == nil && cerr != nil {
				err = cerr
			}
		}

		if !success {
			if rerr := os.Remove(tmpPath); err == nil && rerr != nil {
				err = rerr
			}
		}
	}()

	if _, werr := file.WriteString(xml.Header); werr != nil {
		return werr
	}

	version := strings.TrimSpace(d.Version)
	if version == "" {
		version = "1"
	}

	if _, werr := fmt.Fprintf(file, "<cockatrice_deck version=%q>\n", version); werr != nil {
		return werr
	}

	for _, h := range d.Header {
		if len(h.XML) == 0 {
			continue
		}
		if _, werr := file.Write(h.XML); werr != nil {
			return werr
		}
		if _, werr := file.WriteString("\n"); werr != nil {
			return werr
		}
	}

	enc := xml.NewEncoder(file)
	enc.Indent("    ", "    ")

	for _, z := range d.Zones {
		zStart := xml.StartElement{
			Name: xml.Name{Local: "zone"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "name"}, Value: z.Name},
			},
		}

		if eerr := enc.EncodeToken(zStart); eerr != nil {
			return eerr
		}

		for _, c := range z.Cards {
			cardStart := xml.StartElement{
				Name: xml.Name{Local: "card"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "number"}, Value: fmt.Sprint(c.Count)},
					{Name: xml.Name{Local: "name"}, Value: c.Name},
				},
			}

			if c.SetShortName != "" {
				cardStart.Attr = append(cardStart.Attr, xml.Attr{
					Name:  xml.Name{Local: "setShortName"},
					Value: c.SetShortName,
				})
			}
			if c.CollectorNumber != "" {
				cardStart.Attr = append(cardStart.Attr, xml.Attr{
					Name:  xml.Name{Local: "collectorNumber"},
					Value: c.CollectorNumber,
				})
			}
			if c.UUID != "" {
				cardStart.Attr = append(cardStart.Attr, xml.Attr{
					Name:  xml.Name{Local: "uuid"},
					Value: c.UUID,
				})
			}

			if eerr := enc.EncodeToken(cardStart); eerr != nil {
				return eerr
			}
			if eerr := enc.EncodeToken(xml.EndElement{Name: cardStart.Name}); eerr != nil {
				return eerr
			}
		}

		if eerr := enc.EncodeToken(xml.EndElement{Name: zStart.Name}); eerr != nil {
			return eerr
		}
	}

	if eerr := enc.Flush(); eerr != nil {
		return eerr
	}

	if _, werr := file.WriteString("\n</cockatrice_deck>\n"); werr != nil {
		return werr
	}

	if serr := file.Sync(); serr != nil {
		return serr
	}

	if cerr := file.Close(); cerr != nil {
		return cerr
	}
	closed = true

	// Windows: renaming over an existing file often fails
	if runtime.GOOS == "windows" {
		_ = os.Remove(tDeckFilePath)
	}

	if rerr := os.Rename(tmpPath, tDeckFilePath); rerr != nil {
		return rerr
	}

	success = true
	return nil
}
