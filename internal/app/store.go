package app

// SetStore injects a storage implementation into the App and returns the App for chaining.
func (a *App) SetStore(s Store) *App {
	a.store = s
	return a
}
