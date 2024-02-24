# Tiny World

A tiny, slow-paced world and colony building game.

Made with [Arche](https://github.com/mlange-42/arche) and [Ebitengine](https://github.com/hajimehoshi/ebiten).
Very early work in progress!

![Tiny World screenshot](https://github.com/mlange-42/tiny-world/assets/44003176/9d7b0314-3c29-4773-8670-7e3a9b0df74b)

## Usage

Currently, you need to clone the repository and run the game with [Go](https://go.dev):

```shell
git clone https://github.com/mlange-42/tiny-world.git
cd tiny-world
go run .
```

## Controls

* Middle mouse button / mouse wheel: pan and zoom.
* Left click with selected terrain or land use: place it.
* Right click with selected terrain or land use: remove it.
* Ctrl+S: saves the game to `save/autosave.json`

Load a saved game by running with the file as argument:

```shell
go run . save/autosave.json
```
