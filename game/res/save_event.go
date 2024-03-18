package res

// SaveEvent resource
type SaveEvent struct {
	// Whether the save button was clicked in this tick.
	ShouldSave bool
	// Whether the game should quit and show the main menu.
	ShouldQuit bool
	// Whether the save as map button was clicked in this tick.
	ShouldSaveMap bool
}
