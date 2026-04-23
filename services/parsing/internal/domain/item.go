package domain

type Item struct {
	ID    string
	Name  string
	Wears []string
}

type ItemRow struct {
	ID   string
	Name string
}

type ItemWearRow struct {
	ID   string
	Name string
}
