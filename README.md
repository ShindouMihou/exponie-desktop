<div align="center">

![Exponie](https://exponie.me/banner.png)

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

Other than those, the desktop client downloads the dataset locally and checks for updates on the dataset every start-up, allowing you to play the client 
offline.
