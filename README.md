# FM
A bare bones cross-platform terminal based file manager written in Go

# Screenshot
![IMAGE ALT TEXT HERE](https://raw.githubusercontent.com/maksimKorzh/fm/main/assets/fm.png)

# Features
 - Old school 2 panels look
 - bulk insert files
 - bulk copy
 - bulk delete
 - make directory
 - file/dir highlighting

# Key bindigns
    Arrows: move cursor
    Insert: insert files to copy/delete
    Delete: delete inserted files
      Home: copy inserted files
       End: create new directory
    PgDown: select next drive (windows only)
      PgUp: select previous drive (windows only)

# GNU Win32 Core Utils
    This file manager relies on linux core utils to perform
    commands "cp", "rm" and "mkdir", so if you're on windows
    you need to install GNU Win32 Core Utils that can be found
    in "core-utils" folder, also make sure you update the PATH
    variable so that core utils became available system wide.

# Latest Release
https://github.com/maksimKorzh/fm/releases/

# Build from sources
```bash
cd src
go mod init fm
go build fm.go
```