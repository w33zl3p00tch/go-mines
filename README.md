# go-mines
a Minesweeper clone written in Go using the ebiten game library

```
Left-Click: clear tile
Right-Click: flag tile
Space: new game
```

![Alt text](/images/screenshot.png?raw=true "Screenshot")

This started out as an exercise to test mouse interaction. While already playable, this is pretty much work-in-progress.

Precompiled binaries can be found under "Releases". 

To build it you need a recent golang environment and ebiten with all its dependencies.
See [https://github.com/hajimehoshi/ebiten] and [https://ebiten.org/install.html] for
installation instructions.


## ToDo:
- [x] automatically clear zero-valued tiles and surroundings when clicked
- [ ] find better ways to reduce ebiten's CPU usage
- [ ] inform the player when the game is won
- [x] make board and number of mines configurable
- [ ] add graphics for mines and flagpoles
- [ ] when calling ebiten.SetScreenSize() the screen blacks out for a split second. - Is there a way around this, eg. setting a solid color other than black?

## License
go-mines is licensed under a BSD style license as stated in the LICENSE file.
The font used is "Terminus Font" which is released under the SIL Open Font License.
