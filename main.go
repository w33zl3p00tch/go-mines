package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"
	"github.com/w33zl3p00tch/go-mines/assets"
	"github.com/w33zl3p00tch/go-mines/fonts"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

var (
	dimX            = 256
	dimY            = 256
	tilePlaneCol    = color.RGBA{200, 200, 200, 255}
	tileBorderCol   = color.RGBA{120, 120, 120, 255}
	bgTilePlaneCol  = color.RGBA{140, 140, 140, 255}
	bgTileBorderCol = tileBorderCol
	hlPlaneCol      = color.RGBA{0, 0, 0, 0}
	hlBorderCol     = color.RGBA{255, 255, 255, 255}
	mineCol         = color.RGBA{255, 100, 100, 255}
	black           = color.RGBA{0, 0, 0, 255}
	tileFg          *ebiten.Image
	tileBg          *ebiten.Image
	highlight       *ebiten.Image
	bg              *ebiten.Image
	fg              *ebiten.Image
	gfont           font.Face
	bOne            *ebiten.Image
	bTwo            *ebiten.Image
	bThree          *ebiten.Image
	bFour           *ebiten.Image
	bFive           *ebiten.Image
	bombSlice       []*ebiten.Image
)

type tile struct {
	hasMine          bool
	isClicked        bool
	isFlagged        bool
	minX             int
	minY             int
	maxX             int
	maxY             int
	surroundingMines int
	mineImage        int
}

var field [][]tile

func init() {
	ebiten.SetMaxTPS(30)
	tt, err := truetype.Parse(fonts.Terminus_ttf)
	if err != nil {
		panic(err)
	}
	gfont = truetype.NewFace(tt, &truetype.Options{
		Size: 26,
		DPI:  72,
		//Hinting: font.HintingNone,
	})

	prepareBombs()

	field = prepareField()

	tileImg := generateTile("tile")
	tileFg, _ = ebiten.NewImageFromImage(tileImg, ebiten.FilterDefault)

	tileBgImg := generateTile("tileBg")
	tileBg, _ = ebiten.NewImageFromImage(tileBgImg, ebiten.FilterDefault)

	highlightImg := generateTile("highlight")
	highlight, _ = ebiten.NewImageFromImage(highlightImg, ebiten.FilterDefault)

	fg, _ = ebiten.NewImage(dimX, dimY, ebiten.FilterDefault)
	drawFg(0, 0)
}

// prepareBombs converts images to *ebiten.Image vars
func prepareBombs() {
	img, _, err := image.Decode(bytes.NewReader(assets.Bomb1))
	check(err)
	bOne, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(assets.Bomb2))
	check(err)
	bTwo, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	//img, _, err = image.Decode(bytes.NewReader(assets.Bomb3))
	//check(err)
	//bThree, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(assets.Bomb4))
	check(err)
	bFour, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(assets.Bomb5))
	check(err)
	bFive, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	bombSlice = []*ebiten.Image{bOne, bTwo, bFour, bFive}
}

// prepareField initializes the board and prepares a new game.
func prepareField() [][]tile {
	field := make([][]tile, 8)

	for i := 0; i < len(field); i++ {
		field[i] = make([]tile, 8)
	}

	for i := 0; i < len(field); i++ {
		for k := 0; k < len(field[0]); k++ {
			var t tile
			t.minX = k * 32
			t.maxX = t.minX + 31
			t.minY = i * 32
			t.maxY = t.minY + 31
			field[i][k] = t
		}
	}

	field = placeMines(field)
	field = countSurroundingMines(field)

	return field
}

func placeMines(field [][]tile) [][]tile {
	mines := 10
	minesPlaced := 0
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	maxX := len(field[0])
	maxY := len(field)

	for minesPlaced < mines {
		x := r.Intn(maxX)
		y := r.Intn(maxY)

		if field[y][x].hasMine {
			continue
		}

		field[y][x].hasMine = true
		field[y][x].mineImage = r.Intn(len(bombSlice))
		minesPlaced++
	}

	return field
}

// clearZeroTiles clears all tiles with zero surrounding mines as well as all
// tiles in the immediate vicinity. TODO: optimizations
func clearZeroTiles(x, y int) {
	type position struct {
		x int
		y int
	}

	var zeroTiles []position
	initial := position{x, y}
	zeroTiles = append(zeroTiles, initial)

	maxX := len(field[0]) - 1
	maxY := len(field) - 1

	for len(zeroTiles) > 0 {
		curX := zeroTiles[0].x
		curY := zeroTiles[0].y
		for i := 0; i < 3; i++ {
		INNER:
			for k := 0; k < 3; k++ {
				tX := curX - 1 + k
				tY := curY - 1 + i

				if tX < 0 || tY < 0 || tX > maxX || tY > maxY {
					continue INNER
				}

				t := field[tY][tX]

				if !t.isClicked {
					t.isClicked = true
					if t.surroundingMines == 0 {
						coord := position{tX, tY}
						zeroTiles = append(zeroTiles, coord)
					}
					field[tY][tX] = t
				}
			}
		}
		zeroTiles = zeroTiles[1:]
	}
}

func gameAction(x, y int) {
	t := field[y][x]
	if !t.isClicked {
		t.isClicked = true
		if t.surroundingMines == 0 {
			clearZeroTiles(x, y)
		}
	} else {
		return
	}
	if t.hasMine {
		uncoverAll()
	}
	field[y][x] = t
}

// uncoverAll shows the whole board in a losing situation when a mine was
// tripped
func uncoverAll() {
	for i := 0; i < len(field); i++ {
		for k := 0; k < len(field[0]); k++ {
			field[i][k].isClicked = true
		}
	}
}

func checkMouseAction(x, y int) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		gameAction(x, y)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		field[y][x].isFlagged = !field[y][x].isFlagged
	}
}

func drawFg(x, y int) {
	for i := 0; i < len(field); i++ {
		for k := 0; k < len(field[0]); k++ {
			t := field[i][k]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1.0, 1.0)
			op.GeoM.Translate(float64(0), float64(0))
			op.GeoM.Translate(float64(t.minX), float64(t.minY))
			if !t.isClicked {
				fg.DrawImage(tileFg, op)
				if t.isFlagged {
					text.Draw(fg, "f", gfont, t.minX+10, t.minY+24, black)
				}
			} else if !t.hasMine {
				fg.DrawImage(tileBg, op)
				sm := fmt.Sprint(t.surroundingMines)
				if sm == "0" {
					// don't draw zeros
					continue
				}
				text.Draw(fg, sm, gfont, t.minX+10, t.minY+24, hlBorderCol)
			} else if t.hasMine {
				fg.DrawImage(tileBg, op)
				fg.DrawImage(bombSlice[t.mineImage], op)
				//text.Draw(fg, "#", gfont, t.minX+10, t.minY+24, mineCol)
			}
			if x >= t.minX && y >= t.minY && x <= t.maxX && y <= t.maxY {
				if t.isClicked {
					// nothing more to do for clicked tiles
					continue
				}
				fg.DrawImage(highlight, op)
				checkMouseAction(k, i)
			}
		}
	}
}

func update(screen *ebiten.Image) error {
	x, y := ebiten.CursorPosition()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.0, 1.0)
	drawFg(x, y)
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.DrawImage(fg, op)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		field = prepareField()
	}

	// this delay reduces CPU usage. TODO: find better ways for this.
	// commented out, because this sometimes breaks the SPACE action :(
	//time.Sleep(time.Second / 30)

	return nil
}

func main() {
	if err := ebiten.Run(update, dimX, dimY, 1.0, "go-mines"); err != nil {
		check(err)
	}
}

// generateTile creates a little square image with a border.
func generateTile(tileType string) image.Image {
	w := 32
	h := 32
	tile := image.NewRGBA(image.Rect(0, 0, w, h))

	var borderCol, planeCol color.RGBA

	if tileType == "tile" {
		borderCol = tileBorderCol
		planeCol = tilePlaneCol
	} else if tileType == "highlight" {
		borderCol = hlBorderCol
		planeCol = hlPlaneCol
	} else if tileType == "tileBg" {
		borderCol = bgTileBorderCol
		planeCol = bgTilePlaneCol
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if x == 0 || x == w-1 || y == 0 || y == h-1 {
				tile.Set(x, y, borderCol)
			} else {
				tile.Set(x, y, planeCol)
			}
		}
	}

	return tile
}

// countSurroundingMines of a tile.
func countSurroundingMines(field [][]tile) [][]tile {
	mineX := len(field[0]) + 2
	mineY := len(field) + 2

	mineGrid := make([][]bool, mineY)

	for i := 0; i < len(mineGrid); i++ {
		mineGrid[i] = make([]bool, mineX)
	}

	for i := 0; i < len(field); i++ {
		for k := 0; k < len(field[0]); k++ {
			if field[i][k].hasMine {
				mineGrid[i+1][k+1] = true
			}
		}
	}

	for i := 0; i < len(field); i++ {
		for k := 0; k < len(field[0]); k++ {
			surrMines := &field[i][k].surroundingMines
			//  _______________________
			// |       |       |       |
			// |  y,x  | y,x+1 | y,x+2 |
			// |_______|_______|_______|
			// |       |       |       |
			// | y+1,x |   ?   |y+1,x+2|
			// |_______|_______|_______|
			// |       |       |       |
			// | y+2,x |y+2,x+1|y+2,x+2|
			// |_______|_______|_______|
			//
			for mY := 0; mY < 3; mY++ {
				for mX := 0; mX < 3; mX++ {
					if mX == 1 && mY == 1 {
						continue
					}
					if mineGrid[i+mY][k+mX] {
						*surrMines += 1
					}
				}
			}
		}
	}

	return field
}

// check errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}
