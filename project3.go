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
	map3Path        = "thirdMap.tmx"
)
const (
	ENEMY_FRAME_WIDTH     = 90
	ENEMY_HEIGHT          = 90
	ENEMY_FRAME_COUNT     = 4
	ENEMY_FRAME_PER_SHEET = 3
)
const (
	ENEMY_UP = iota
	ENEMY_RIGHT
	ENEMY_DOWN
	ENEMY_LEFT
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
	enemy       enemy
	enemy2      enemy
}
type enemy struct {
	sprite     *ebiten.Image
	xLocNnemy  int
	yLocEnemy  int
	frame      int
	frameDelay int
	direction  int
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

	demoGame.enemy.frameDelay += 1
	if demoGame.enemy.frameDelay%ENEMY_FRAME_COUNT == 0 {
		demoGame.enemy.frame += 1
		if demoGame.enemy.frame >= ENEMY_FRAME_PER_SHEET {
			demoGame.enemy.frame = 0
		}
		if demoGame.enemy.direction == ENEMY_RIGHT {
			demoGame.enemy.xLocNnemy += 3
			if demoGame.enemy.xLocNnemy >= 350 {
				demoGame.enemy.direction = ENEMY_LEFT
			}
		} else if demoGame.enemy.direction == ENEMY_LEFT {
			demoGame.enemy.xLocNnemy -= 3
			if demoGame.enemy.xLocNnemy <= 200 {
				demoGame.enemy.direction = ENEMY_RIGHT
			}
		}
	}

	demoGame.enemy2.frameDelay += 1
	if demoGame.enemy2.frameDelay%ENEMY_FRAME_COUNT == 0 {
		demoGame.enemy2.frame += 1
		if demoGame.enemy2.frame >= ENEMY_FRAME_PER_SHEET {
			demoGame.enemy2.frame = 0
		}
		if demoGame.enemy2.direction == ENEMY_RIGHT {
			demoGame.enemy2.xLocNnemy += 3
			if demoGame.enemy2.xLocNnemy >= 750 {
				demoGame.enemy2.direction = ENEMY_LEFT
			}
		} else if demoGame.enemy2.direction == ENEMY_LEFT {
			demoGame.enemy2.xLocNnemy -= 3
			if demoGame.enemy2.xLocNnemy <= 570 {
				demoGame.enemy2.direction = ENEMY_RIGHT
			}
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
	if demoGame.playerXLoc >= 950 {
		if demoGame.levels == 1 {
			gameMap3, err := tiled.LoadFile(map3Path)
			if err != nil {
				fmt.Printf("error parsing map: %s", err.Error())
				return err
			}
			demoGame.levels = 2
			demoGame.playerXLoc = 0
			demoGame.playerYLoc = 448
			demoGame.level = gameMap3
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

	if demoGame.levels == 1 {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.enemy.xLocNnemy), float64(demoGame.enemy.yLocEnemy))
		screen.DrawImage(demoGame.enemy.sprite.SubImage(image.Rect(demoGame.enemy.frame*ENEMY_FRAME_WIDTH,
			demoGame.enemy.direction*ENEMY_HEIGHT,
			demoGame.enemy.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
			demoGame.enemy.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)

		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.enemy2.xLocNnemy), float64(demoGame.enemy2.yLocEnemy))
		screen.DrawImage(demoGame.enemy2.sprite.SubImage(image.Rect(demoGame.enemy2.frame*ENEMY_FRAME_WIDTH,
			demoGame.enemy2.direction*ENEMY_HEIGHT,
			demoGame.enemy2.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
			demoGame.enemy2.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)
	}

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
	enemyAnimation := LoadEmbeddedImage("", "skelly.png")

	oneLevelGame := AnimatedSpriteDemo3{
		levels:      0,
		spriteSheet: animationGuy,
		playerXLoc:  256,
		playerYLoc:  448,
		level:       gameMap,
		tileHash:    ebitenImageMap,
		enemy: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 200,
			yLocEnemy: 608,
			direction: ENEMY_RIGHT,
		},
		enemy2: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 570,
			yLocEnemy: 384,
			direction: ENEMY_LEFT,
		},
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
