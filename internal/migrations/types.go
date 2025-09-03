package migrations

// Migration represents a database migration
type Migration struct {
	Version    int64
	Concurrent bool
	Path       string
	Name       string
}
