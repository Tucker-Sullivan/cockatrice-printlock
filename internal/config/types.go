package config

import "strings"

type Config struct {
	DeckDir                string   `json:"deckDir"`
	CardsXMLPath           string   `json:"cardsXmlPath"`
	GlobalSetPriority      []string `json:"globalSetPriority"`
	GlobalSetsHavePriority bool     `json:"globalSetsHavePriority"`
}

func NewConfig(cardsPath, deckDir string) *Config {
	return &Config{
		CardsXMLPath:           cardsPath,
		DeckDir:                deckDir,
		GlobalSetPriority:      []string{},
		GlobalSetsHavePriority: false,
	}
}

func (c *Config) sanitize() {
	c.CardsXMLPath = strings.TrimSpace(c.CardsXMLPath)
	c.DeckDir = strings.TrimSpace(c.DeckDir)
	cleanedSets := make([]string, 0, len(c.GlobalSetPriority))
	for _, setCode := range c.GlobalSetPriority {
		tSetCode := strings.TrimSpace(setCode)
		if tSetCode == "" {
			continue
		}
		cleanedSets = append(cleanedSets, tSetCode)
	}
	c.GlobalSetPriority = cleanedSets
}
