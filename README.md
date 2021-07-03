# Cataclysm: DDA (Unofficial) Linux Launcher

### This is a Linux launcher for [Cataclysm Dark Days Ahead](https://github.com/CleverRaven/Cataclysm-DDA) written in Go. 

You may need to install CDDA dependencies based on your distribution.

Build: `$ go build`

Usage: 
- `$ ./cclauncher` show a generic usage of the tool.
- `$ ./cclauncher download` will download the latest Tiles build of C:DDA, backup the previous installation and restore any tilesets, saves, soundpacks etc.

Help: `$ ./cclauncher -h` or the `-h` with any subcommands will show additional flags and options.



*What will be saved across updates?* <br>
Everything inside *save*, *config* and every custom folder inside *gfx*, *data/sound*, *data/mods* (so if you modify in-repo mods they will be overwritten!)