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
as well as required achievement, can be tweaked by editing the exported scenario.

## Map Format

A small example map is shown below.

* The 1st line contains frequencies of random terrains.
* The second line contains the number of initially placable trains.
* The 3rd line contains a list of required achievements, separated by spaces.
* The 4th line contains the relative coordinates of the starting position, from the top-left corner (0,0).
* All further lines are the actual map.

Terrain characters are defined in [`data/json/terrain.json`](https://github.com/mlange-42/tiny-world/blob/main/data/json/terrain.json).
Achievements are defined in [`data/json/achievements.json`](https://github.com/mlange-42/tiny-world/blob/main/data/json/achievements.json)

```
1r 20- 6^ 6~ 1+ 6t
500
play-the-game
8 8
......---------~....
....-----------~-...
...------------~t-..
..-----TTT----t~t--.
.-----TTTT----t~~t-.
.----rTTT-----rt~t--
------Tt------rt~t--
---------------t~~--
--ttt---h------tt~--
~~~~~t--------ttt~--
---t~~t----t~~~~t~--
----t~t---tt~tt~~~--
.----~~~~~~~~t------
..--------ttt-^^---.
...---------^^----..
....---tt--------...
.....--tt-------....
......---------.....
```
