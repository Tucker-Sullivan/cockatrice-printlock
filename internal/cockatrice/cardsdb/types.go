package cardsdb

type SetEntry struct {
	Code     string // <sets><set><name> e.g. "DRC"
	LongName string // <longname>
}

type Printing struct {
	UUID            string // attribute uuid=""
	SetCode         string // inner text of <set> e.g. "DRC"
	CollectorNumber string // attribute num="" (string; may be "249★", "CMD-261")
	MultiverseID    string // attribute muid="" (optional, but cheap to keep)
}

type CardEntry struct {
	Name      string
	Printings []Printing
}

type CardsDB struct {
	SetsByCode  map[string]SetEntry
	CardsByName map[string]CardEntry
}

func NewCardsDB() *CardsDB {
	return &CardsDB{
		SetsByCode:  make(map[string]SetEntry),
		CardsByName: make(map[string]CardEntry),
	}
}
