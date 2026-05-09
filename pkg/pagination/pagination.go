package pagination

// Page holds the input parameters for a paginated query.
// PageState is the opaque cursor returned by Cassandra — pass nil on the first request
// and the value from the previous PagedResult to fetch the next page.
type Page struct {
	Size      int    // number of rows per page (required, > 0)
	PageState []byte // cursor from the previous response; nil = first page
}

// PagedResult wraps a slice of results together with the next page cursor.
// When NextPageState is nil the caller has reached the last page.
type PagedResult[T any] struct {
	Items         []T
	NextPageState []byte // nil when there are no more pages
}
