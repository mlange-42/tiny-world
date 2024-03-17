# Changelog

## [[unpublished]](https://github.com/mlange-42/tiny-world/compare/v0.1.0...main)

### Game features

* The game can be played on pre-defined, embedded maps/scenarios (#168)
* The game can be played on user-created maps/scenarios (#171)
* Adds six scenarios (#168, #207, #221, #226)
* Adds achievements that are required to unlock maps/scenarios (#187, #194, #195, #196, #203)
* Adds a map editor mode for easier scenario creation (#223)
* Adds hills terrain, which allows no farming (#218)
* Adds shepherds and pastures for food production on hills (#218)
* Limits the total number of placable terrain tiles / natural features (#231)
* Adds churches and monasteries that unlock more tiles (#231)

### Game rules

* Haulers stay in buildings for a second before continuing (#163)
* Watermill consumes 1 wood, like windmill, instead of 1 stone (#33)

### Usability

* Adds a menu to save/quit to main menu from in the game (#179)
* Shows the build radius around the cursor when building a castle or tower (#183)
* Show information on why something can't be built in button tooltips (#185)
* Highlight resources with negative net production in the info bar (#186)
* Adds a "Save map" button to the in-game menu (#190)
* Adds scroll containers for main menu lists (#191)
* Fixes the in-game UI layout for small screens (#193)
* Adds a status bar for information like why a building has a warning sign (#213, #214)
* Use square brackets instead of PageUp/PageDown to control game speed (#215)

### Graphics

* Adds random variants for grass tiles (#158)
* Tweak info labels to accommodate 3-digit values without size changes (#164)
* Complete rework and themed styling of the main menu (#176, #177)
* Text highlight color is a tileset property (#186)
* Adds the possibility for multiple terrains below a terrain (used for shepherds) (#220)

### Documentation

* Adds doc-strings for all resources (#159)
* Adds documentation on scenario editing and creation under [`docs/SCENARIOS.md`](https://github.com/mlange-42/tiny-world/blob/main/docs/SCENARIOS.md) (#224)

### Other

* Adds precompiled binaries for MacOS/Darwin to release builds (#165 by [Ecostack](https://github.com/Ecostack))
* Adds a command line tool under `cmd/stats` to print terrain frequencies of maps, useful for deriving random terrain probabilities (#223)

## [[v0.1.0]](https://github.com/mlange-42/tiny-world/tree/v0.1.0)

First release of Tiny World.