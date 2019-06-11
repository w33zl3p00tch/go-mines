# go-mines
a Minesweeper clone written in Go using the ebiten game library


Left-Click: clear tile
Right-Click: flag tile
Space: new game


This started out as an exercise to test mouse interaction. While already playable, this is pretty much work-in-progress.

To build it you need a recent golang environment and ebiten with all its dependencies.
See [https://github.com/hajimehoshi/ebiten] and [https://ebiten.org/install.html] for
installation instructions.


## ToDo:
- automatically clear zero-valued tiles and surroundings when clicked
- find better ways to reduce ebiten's CPU usage
- inform the player when the game is won
- make board and number of mines configurable
- add graphics for mines and flagpoles

## License
go-mines is licensed under a BSD style license as stated in the LICENSE file.
The font used is "Terminus Font" which is released under the SIL Open Font License.
