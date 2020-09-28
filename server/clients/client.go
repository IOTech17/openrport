package clients

import "sort"

// Client represents rport clients.
type Client struct {
	ID          string `json:"id"`
	Password    string `json:"password"`
	lockedBySID string
}

func SortByID(a []*Client, desc bool) {
	sort.Slice(a, func(i, j int) bool {
		less := a[i].ID < a[j].ID
		if desc {
			return !less
		}
		return less
	})
}
