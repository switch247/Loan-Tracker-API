package Domain

type LoanCollections struct {
	Users         Collection
	RefreshTokens Collection
	ResetTokens   Collection
	Loans         Collection
	Logs          Collection
}

type Filter struct {
	Title      string
	Slug       string
	LoanerName string
	Page       int
	Limit      int
	Tags       []string
	SortBy     string
	OrderBy    int
}

type PaginationMetaData struct {
	TotalRecords int
	TotalPages   int
	PageSize     int
	CurrentPage  int
}
