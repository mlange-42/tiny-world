# Tiny World Scenarios

Tiny World can be played on built-in as well as on user-provided maps or scenarios.
We call them scenarios, because each one can have its own probabilities for random terrain "cards".

**Built-in scenarios** are packaged with the executable, and can be found in the
[GitHub repository](https://github.com/mlange-42/tiny-world) under [`data/maps`](https://github.com/mlange-42/tiny-world/tree/main/data/maps).

**User scenarios** must be placed in a folder `maps` next to the executable.

## Creating Scenarios

Most of the scenario creation process can be done in the game.
To run the game in editor mode, start a new game with Shift+Click in the main menu.
Existing scenarios can also be started in editor mode with Shift+Click.

Controls in editor mode are the same as in a usual game.
However, all natural terrains are available and there are no costs or cost restrictions.
Further, terrain is not only placed by clicking, but also by dragging the mouse.

Saving a game in editor mode is like saving the editor session. Saving the map means to export it in the [Map Format](#map-format) for use as a scenario.

In a final step, random terrain frequencies for the scenario,
as well as required achievement and the map description can be tweaked by editing the exported scenario.

## Map Format

Maps are saved in JSON format. See the example below.

* `terrains` contains frequencies of random terrains.
* `initial_terrains` is the number of initially placable trains.
* `achievements` is a list of required achievements.
* `description` is the scenario description shown in the main menu tooltip.
* `center` is the relative starting position, from the top-left corner (0,0).
* `map` is the actual map.

Terrain characters are defined in [`data/json/terrain.json`](https://github.com/mlange-42/tiny-world/blob/main/data/json/terrain.json).
Achievements are defined in [`data/json/achievements.json`](https://github.com/mlange-42/tiny-world/blob/main/data/json/achievements.json)

```json
{
  "terrains": {
    "+": 1,
    "-": 20,
    "^": 4,
    "r": 1,
    "t": 6,
    "~": 6
  },
  "map": [
    "......---------~....",
    "....-----------~-...",
    "...------------~t-..",
    "..-----TTT----t~t--.",
    ".-----TTTT----t~~t-.",
    ".----rTTT-----rt~t--",
    "------Tt------rt~t--",
    "---------------t~~--",
    "--ttt---h------tt~--",
    "~~~~~t--------ttt~--",
    "---t~~t----t~~~~t~--",
    "----t~t---tt~tt~~~--",
    ".----~~~~~~~~t------",
    "..--------ttt-^^---.",
    "...---------^^----..",
    "....---tt--------...",
    ".....--tt-------....",
    "......---------....."
  ],
  "achievements": [
    "play-the-game"
  ],
  "description": [
    "A small (20x20) starting area with a river.",
    "Available random tiles: 500."
  ],
  "center": {
    "X": 8,
    "Y": 8
  },
  "initial_terrains": 500
}
```
