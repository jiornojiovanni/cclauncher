# Cataclysm: DDA Linux Launcher

### This is a Linux launcher for [Cataclysm Dark Days Ahead](https://github.com/CleverRaven/Cataclysm-DDA) written in Go. 

You may need to install CDDA dependencies based on your distribution.

Updates are not implemented right now so you should delete CDDA folder before calling the program again.

## IT WON'T BACKUP YOUR SAVES! YOU ARE WARNED.

Build: `$ go build`

Usage: 
- `$ ./cclauncher` will download the last build of Cdda, *tiles*, and extract it.
- `$ ./cclauncher -d` will only download the archive. (**optional**)
- `$ ./cclauncher -v [VERSION]` downloads a specific version (**NOT IMPLEMENTED**) (default to **latest**)
- `$ ./cclauncher -g [t/c]` download tiles or ncurses version (default to [**t**]iles)

Example: `$ ./cclauncher -d -g c` downloads only the archive of the latest Ncurses build


Help: `$ ./cclauncher -h` Show help.








#### Currently Implemented: 
- Download of the last build (tiles or curses)
- Extraction
- Barebones cli

#### TODO:
- Download of a specific version.
- Updates
- Backups
- Show changelog
- Bright Nights download (Maybe)


## This is not currently usable for daily use!
Stick to your package managers or manual downloads for now.

###### This is probably bad Go code, feel free to help me correcting it :D !