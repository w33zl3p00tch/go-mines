package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"
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
	field = calculateSurroundingMines(field)

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
		minesPlaced++
	}

	return field
}

func gameAction(x, y int) {
	t := field[y][x]
	if !t.isClicked {
		t.isClicked = true
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
				text.Draw(fg, sm, gfont, t.minX+10, t.minY+24, hlBorderCol)
			} else if t.hasMine {
				fg.DrawImage(tileBg, op)
				text.Draw(fg, "#", gfont, t.minX+10, t.minY+24, mineCol)
			}
			if x >= t.minX && y >= t.minY && x <= t.maxX && y <= t.maxY {
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
	time.Sleep(time.Second / 30)

	return nil
}

func main() {
	if err := ebiten.Run(update, dimX, dimY, 1.0, "go-mines"); err != nil {
		log.Fatal(err)
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

// calculateSurroundingMines is a horrendous abomination. Every tile's
// neighboring fields are checked for mines and added to tile.surroundingMines.
// While it works, this should be rewritten to be more legible and maybe even
// elegant.
func calculateSurroundingMines(field [][]tile) [][]tile {
	for y := 0; y < len(field); y++ {
		for x := 0; x < len(field[0]); x++ {
			if x == 0 && y == 0 {
				// upper left corner
				if field[0][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][0].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x == len(field[0])-1 && y == 0 {
				// upper right corner
				if field[0][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][x].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x == 0 && y == len(field)-1 {
				// lower left corner
				if field[len(field)-2][0].hasMine {
					field[y][x].surroundingMines++
				}
				if field[len(field)-2][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[len(field)-1][1].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x == len(field[0])-1 && y == len(field)-1 {
				// lower right corner
				if field[len(field)-1][len(field[0])-2].hasMine {
					field[y][x].surroundingMines++
				}
				if field[len(field)-2][len(field[0])-2].hasMine {
					field[y][x].surroundingMines++
				}
				if field[len(field)-2][len(field[0])-1].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x > 0 && x < len(field[0])-1 && y == 0 {
				// upper edge
				if field[0][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][x].hasMine {
					field[y][x].surroundingMines++
				}
				if field[1][x+1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[0][x+1].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x == 0 && y > 0 && y < len(field)-1 {
				// left edge
				if field[y-1][0].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][0].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x == len(field[0])-1 && y > 0 && y < len(field)-1 {
				// right edge
				if field[y-1][x].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][x].hasMine {
					field[y][x].surroundingMines++
				}
			} else if x > 0 && x < len(field[0])-1 && y == len(field)-1 {
				// lower edge
				if field[y][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x+1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y][x+1].hasMine {
					field[y][x].surroundingMines++
				}
			} else {
				// finally we're back in sanity-land
				if field[y-1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y-1][x+1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y][x+1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][x+1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][x].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y+1][x-1].hasMine {
					field[y][x].surroundingMines++
				}
				if field[y][x-1].hasMine {
					field[y][x].surroundingMines++
				}
			}
		}
	}

	return field
}
