# Tiny World

A tiny, slow-paced world and colony building game.

Made with [Ebitengine](https://github.com/hajimehoshi/ebiten) and the [Arche](https://github.com/mlange-42/arche) Entity Component System.

<div align="center" width="100%">
<img alt="Tiny World screenshot" src="https://github.com/mlange-42/tiny-world/assets/44003176/b3384739-af7c-4f44-996f-8f1cb5097fa3"></img>
</div>

## Usage

### Play in browser

You can play Tiny World [on itch.io](https://mlange-42.itch.io/tiny-world),
or the development version [here](https://mlange-42.github.io/tiny-world/).

### Precompiled binaries

You can download precompiled binaries for Linux, Windows and macOS from the [Releases](https://github.com/mlange-42/tiny-world/releases).

#### macOS version
For the macOS version, please right-click the app and select "Open" to bypass the security warning, as the binary is not signed.

In case you get the message `“tiny-world.app” is damaged and can’t be opened. You should move it to the Bin.`, please use the following command from the terminal:
```shell
xattr -c tiny-world.app
```
This will remove the quarantine attribute from the app. You can then open it as usual.

### Build from source

Clone the repository and build or run the game with [Go](https://go.dev):

```shell
git clone https://github.com/mlange-42/tiny-world.git
cd tiny-world
go run .
```

For building on Unix systems, `libgl1-mesa-dev` and `xorg-dev` are required.

## Playing

In the toolbar on the right, the top items are **buildings** that can be bought by the player for resources.
The **natural features** in the lower part appear randomly and are replenished when placed by the player.

* Pan: Arrows, WASD or middle mouse button
* Zoom: +/- or mouse wheel
* Pause/resume: Space
* Game speed: [/] (square brackets)
* Toggle fullscreen: F11

All UI controls have tooltips. Read them carefully!
