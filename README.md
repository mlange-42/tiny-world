# Tiny World

A tiny, slow-paced world and colony building game.

Made with [Arche](https://github.com/mlange-42/arche) and [Ebitengine](https://github.com/hajimehoshi/ebiten).
Very early work in progress!

![Tiny World screenshot](https://github.com/mlange-42/tiny-world/assets/44003176/222bc13d-1ea1-43d0-9698-8f9d94feeade)

## Usage

### Play in browser

You can play Tiny World in a web browser [here](https://mlange-42.github.io/tiny-world/).
However, the browser version does not yet support saving a game.

### Play locally

Currently, you need to clone the repository and run the game with [Go](https://go.dev):

```shell
git clone https://github.com/mlange-42/tiny-world.git
cd tiny-world
go run .
```

## Playing

In the toolbar on the right, the top items are **buildings** that can be bought by the player for resources.
The **natural features** in the lower part appear randomly and are replenished when placed by the player.

* Middle mouse button / mouse wheel: pan and zoom.
* Space: pause/resume
* Left click with selected terrain or buildable: place it.
* Right click with selected buildable: remove it.
* Ctrl+S: saves the game to `save/autosave.json`


Load a saved game by running with the `-s` option:

```shell
go run . -s save/autosave.json
```
