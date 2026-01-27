package config

type Config struct {
	DeckDir           string   `json:"deckDir"`
	CardsXMLPath      string   `json:"cardsXmlPath"`
	GlobalSetPriority []string `json:"globalSetPriority"`
}
