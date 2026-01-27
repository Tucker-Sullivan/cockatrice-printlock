package deck

type Deck struct {
	Version string
	Header  []RawNode
	Zones   []Zone
}

type RawNode struct {
	Name string
	XML  []byte
}

type Zone struct {
	Name  string
	Cards []Card
}

type Card struct {
	Count           int
	Name            string
	SetShortName    string
	CollectorNumber string
	UUID            string
}

func NewDeck() *Deck {
	return &Deck{
		Version: "1",
		Header:  make([]RawNode, 0),
		Zones:   make([]Zone, 0),
	}
}
