package printlock

import "fmt"

type ApplyOptions struct {
	PreferHigherCollectionNumber bool
	SetPriority                  []string
}

type ApplyError struct {
	Overwritten bool
	Zone        string
	Card        string
	Reason      string
}

func (e ApplyError) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Card, e.Zone, e.Reason)
}
