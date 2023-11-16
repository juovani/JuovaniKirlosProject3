package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"image"
	"log"
	"os"
	"path"
	//"golang.org/x/sys/unix"
	//"image"
	//"log"
	//"path"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	GUY_FRAME_WIDTH = 45
	GUY_HEIGHT      = 60
	FRAME_COUNT     = 4
	FRAME_PER_SHEET = 8
	map1Path        = "firstMap.tmx"
	map2Path        = "secondMap.tmx"
)
const (
	DOWN = iota
	UP
	LEFT
	RIGHT
)

type AnimatedSpriteDemo3 struct {
	spriteSheet *ebiten.Image
	playerXLoc  int
	playerYLoc  int
	direction   int
	frame       int
	frameDelay  int
	level       *tiled.Map
	levels      int
	tileHash    map[uint32]*ebiten.Image
}

func (demoGame *AnimatedSpriteDemo3) Update() error {
	demoGame.frameDelay += 1
	if demoGame.frameDelay%FRAME_COUNT == 0 {
		demoGame.frame += 1
		if demoGame.frame >= FRAME_PER_SHEET {
			demoGame.frame = 0
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && demoGame.playerXLoc > 0 {
			demoGame.direction = LEFT
			demoGame.playerXLoc -= 7
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && demoGame.playerXLoc < 1000-GUY_FRAME_WIDTH {
			demoGame.direction = RIGHT
			demoGame.playerXLoc += 10
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && demoGame.playerYLoc < 1000-GUY_HEIGHT {
			demoGame.direction = UP
			demoGame.playerYLoc -= 7
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && demoGame.playerYLoc > 0 {
			demoGame.direction = DOWN
			demoGame.playerYLoc += 7
		} else {
			demoGame.frame = 0
		}
	}
	if demoGame.playerXLoc >= 950 {
		if demoGame.levels == 0 {
			gameMap2, err := tiled.LoadFile(map2Path)
			if err != nil {
				fmt.Printf("error parsing map: %s", err.Error())
				return err
			}
			demoGame.levels = 1
			demoGame.playerXLoc = 0
			demoGame.playerYLoc = 128
			demoGame.level = gameMap2
		}
	}
	return nil
}

func (demoGame AnimatedSpriteDemo3) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	for _, layer := range demoGame.level.Layers {
		for tileY := 0; tileY < demoGame.level.Height; tileY++ {
			for tileX := 0; tileX < demoGame.level.Width; tileX++ {
				drawOptions.GeoM.Reset()
				tileXpos := float64(demoGame.level.TileWidth * tileX)
				tileYpos := float64(demoGame.level.TileHeight * tileY)
				drawOptions.GeoM.Translate(tileXpos, tileYpos)
				tileToDraw := layer.Tiles[tileY*demoGame.level.Width+tileX]
				if tileToDraw.ID <= 0 {
					continue
				}
				ebitenTileToDraw := demoGame.tileHash[tileToDraw.ID]
				if ebitenTileToDraw == nil {
					fmt.Printf("Nil tile image for tile ID: #{tileToDraw.ID}\n")
					continue
				}
				screen.DrawImage(ebitenTileToDraw, &drawOptions)

			}
		}
	}

	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(demoGame.playerXLoc), float64(demoGame.playerYLoc))
	screen.DrawImage(demoGame.spriteSheet.SubImage(image.Rect(demoGame.frame*GUY_FRAME_WIDTH,
		demoGame.direction*GUY_HEIGHT,
		demoGame.frame*GUY_FRAME_WIDTH+GUY_FRAME_WIDTH,
		demoGame.direction*GUY_HEIGHT+GUY_HEIGHT)).(*ebiten.Image), &drawOptions)

}

func (demoGame *AnimatedSpriteDemo3) GetCurrentLevel() int {
	return demoGame.levels
}

func main() {
	gameMap, err := tiled.LoadFile(map1Path)
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	ebitenImageMap := makeEbitenImagesFromMap(*gameMap)
	fmt.Println("tilesets:", gameMap.Tilesets[0].Tiles)
	fmt.Print("type:", fmt.Sprintf("%T", gameMap.Layers[0].Tiles[0]))

	animationGuy := LoadEmbeddedImage("", "player.png")

	oneLevelGame := AnimatedSpriteDemo3{
		levels:      0,
		spriteSheet: animationGuy,
		playerXLoc:  256,
		playerYLoc:  448,
		level:       gameMap,
		tileHash:    ebitenImageMap,
	}

	err = ebiten.RunGame(&oneLevelGame)
	if err != nil {
		fmt.Println("Couldn't run game:", err)
	}
}
func (demoGame AnimatedSpriteDemo3) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}
func makeEbitenImagesFromMap(tiledMap tiled.Map) map[uint32]*ebiten.Image {
	idToImage := make(map[uint32]*ebiten.Image)
	for _, tile := range tiledMap.Tilesets[0].Tiles {
		ebitenImageTile, _, err := ebitenutil.NewImageFromFile(tile.Image.Source)
		if err != nil {
			fmt.Println("Error loading tile image:", tile.Image.Source, err)
		}
		idToImage[tile.ID] = ebitenImageTile
	}
	return idToImage
}
func LoadEmbeddedImage(folderName string, imageName string) *ebiten.Image {
	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", folderName, imageName))
	if err != nil {
		log.Fatal("failed to load embedded image ", imageName, err)
	}
	ebitenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("Error loading tile image:", imageName, err)
	}
	return ebitenImage
}
