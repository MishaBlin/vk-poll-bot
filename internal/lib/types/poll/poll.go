package poll

type Poll struct {
	ID      string
	Title   string
	OwnerID string
	Options []string
	Votes   []int
	Voters  map[string]int
	Active  bool
}
