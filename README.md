<div align="center">

![Exponie](https://exponie.mihou.pw/banner.png)

a minimalist's choice of spelling test.
</div>

##

Exponentia is a minimal, straight-to-the-point spelling test that works both online and limited offline, this repository 
is dedicated to the desktop client which is written with [`wails.io`](https://wails.io), allowing users to play the spelling 
test offline.

You can also play the web-version, which tends to be more updated, via the [`exponie.me`](https://exponie.me) link.

## Demo
![image](https://github.com/ShindouMihou/exponie-desktop/assets/69381903/4dbc7d95-cb55-4ff4-80b9-99cfdf77946f)

## Differences

The desktop client was designed to be a simpler port of Exponentia, therefore, there were changes made:
1. **Non-configurable settings**: As most of the opt-able settings were mostly unavailable in the desktop client, we decided to disable settings for now.
2. **Text-to-speech**: We disabled text-to-speech because the desktop client doesn't support any known ways for text-to-speech so far.
3. **Offline Definitions**: The desktop client has definitions available even when offline, this is because we can store large files over the filesystem, which is not possible with browser.

Other than those, the desktop client downloads the dataset locally and checks for updates on the dataset every start-up, allowing you to play the client 
offline.

## Installation

Windows users can download the installer, or the standalone executable on [`GitHub Releases`](https://github.com/ShindouMihou/exponie-desktop/releases). We 
recommend using the installer if you want to have the application findable using the Start Menu, and other related. You can use the standalone executable 
to run the application directly, if you want.

To compile for your own platform, you can clone the repository and install [`wails.io`](https://wails.io) before running the following command:
```shell
wails build -upx -clean
```
The above command should compile the application in less than 20 seconds and create a binary for your platform under `build/bin`. We recommend having `upx` installed 
to compress the binary.
