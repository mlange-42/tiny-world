## Installation Guide

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
