# Cataclysm: DDA Linux Launcher

### This is a Linux launcher for [Cataclysm Dark Days Ahead](https://github.com/CleverRaven/Cataclysm-DDA) written in Go. 

You may need to install CDDA dependencies based on your distribution.

Build: `$ go build`

Usage: 
- `$ ./cclauncher` will download the last build of Cdda, *tiles*, show the changelog, and extract it.
- `$ ./cclauncher -d` will only download the archive. (**optional**)
- `$ ./cclauncher -v [VERSION]` downloads a specific version. (default to **latest**)
- `$ ./cclauncher -g [t/c]` download tiles or ncurses version. (default to [**t**]iles)
- `$ ./cclauncher -no-c` this flag disable the changelog.

Example: `$ ./cclauncher -d -g c` downloads only the archive of the latest Ncurses build

Another example with `$ ./cclauncher -d -v 11476`
[![test.png](https://i.postimg.cc/44jmMDfV/test.png)](https://postimg.cc/pyYW51nL)

Help: `$ ./cclauncher -h` Show help.





#### Currently Implemented: 
- Download of the last build (tiles or curses)
- Extraction
- Barebones cli
- Download of a specific version
- Show changelog
- Updates
- Backups

#### TODO:
- Fix this mess
- Tileset backup

## Tileset are not backupped right now, be careful.
## This is not currently usable for daily use!
Stick to your package managers or manual downloads for now.