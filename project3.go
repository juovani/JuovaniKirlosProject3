package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/paths"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"image"
	"log"
	"math"
	"os"
	"path"
	"strings"
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
	ENEMY_FRAME_WIDTH     = 72
	ENEMY_HEIGHT          = 96
	ENEMY_FRAME_COUNT     = 4
	ENEMY_FRAME_PER_SHEET = 3
)
const (
	NPC_FRAME_WIDTH = 60
	NPC_HEIGHT      = 80
	//NPC_FRAME_COUNT     = 2
	//NPC_FRAME_PER_SHEET = 3
)
const (
	NPC_DOWN = 2
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
const (
	COIN_FRAME_WIDTH     = 27
	COIN_HEIGHT          = 27
	COIN_FRAME_COUNT     = 1
	COIN_FRAME_PER_SHEET = 8
)
const (
	COIN_RIGHT = 0
)

type AnimatedSpriteDemo3 struct {
	spriteSheet    *ebiten.Image
	bullet         *ebiten.Image
	playerXLoc     int
	playerYLoc     int
	direction      int
	frame          int
	frameDelay     int
	damage         int
	temp           bool
	level          *tiled.Map
	levels         int
	tileHash       map[uint32]*ebiten.Image
	pathFindingMap []string
	pathMap        *paths.Grid
	path           *paths.Path
	enemy1         enemy
	enemy2         enemy
	enemy3         enemy
	enemy4         enemy
	npc1           npc
	shot           []shots
	msg            bool
	textFont       font.Face
	coin1          coins
	coin2          coins
	coin3          coins
	coin4          coins
}
type coins struct {
	sprite     *ebiten.Image
	frame      int
	frameDelay int
	coinXLoc   int
	coinYLoc   int
	direction  int
	pickedUp   bool
}
type shots struct {
	sprite     *ebiten.Image
	direction  int
	bulletXLoc int
	bulletYLoc int
}
type enemy struct {
	sprite     *ebiten.Image
	xLocNnemy  int
	yLocEnemy  int
	frame      int
	frameDelay int
	direction  int
	health     int
	alive      bool
}
type npc struct {
	sprite     *ebiten.Image
	xLocNPC    int
	yLocNPC    int
	frame      int
	frameDelay int
	direction  int
}

func newShots(image *ebiten.Image) shots {
	return shots{
		sprite: image,
	}
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		demoGame.temp = true

		newBullet := newShots(demoGame.bullet)
		newBullet.bulletXLoc = demoGame.playerXLoc
		newBullet.bulletYLoc = demoGame.playerYLoc
		newBullet.direction = demoGame.direction
		demoGame.shot = append(demoGame.shot, newBullet)
	} else {
		for i := range demoGame.shot {
			if demoGame.shot[i].direction == RIGHT {
				demoGame.shot[i].bulletXLoc += 4
			} else if demoGame.shot[i].direction == LEFT {
				demoGame.shot[i].bulletXLoc -= 4
			} else if demoGame.shot[i].direction == UP {
				demoGame.shot[i].bulletYLoc -= 4
			} else if demoGame.shot[i].direction == DOWN {
				demoGame.shot[i].bulletYLoc += 4
			}
		}
	}
	//if demoGame.enemy1.alive == true {
	//	demoGame.enemy1.frameDelay += 1
	//	if demoGame.enemy1.frameDelay%ENEMY_FRAME_COUNT == 0 {
	//		demoGame.enemy1.frame += 1
	//		if demoGame.enemy1.frame >= ENEMY_FRAME_PER_SHEET {
	//			demoGame.enemy1.frame = 0
	//		}
	//		if demoGame.enemy1.direction == ENEMY_RIGHT {
	//			demoGame.enemy1.xLocNnemy += 3
	//			if demoGame.enemy1.xLocNnemy >= 350 {
	//				demoGame.enemy1.direction = ENEMY_LEFT
	//			}
	//		} else if demoGame.enemy1.direction == ENEMY_LEFT {
	//			demoGame.enemy1.xLocNnemy -= 3
	//			if demoGame.enemy1.xLocNnemy <= 200 {
	//				demoGame.enemy1.direction = ENEMY_RIGHT
	//			}
	//		}
	//	}
	//
	//}
	//if demoGame.enemy2.alive == true {
	//	demoGame.enemy2.frameDelay += 1
	//	if demoGame.enemy2.frameDelay%ENEMY_FRAME_COUNT == 0 {
	//		demoGame.enemy2.frame += 1
	//		if demoGame.enemy2.frame >= ENEMY_FRAME_PER_SHEET {
	//			demoGame.enemy2.frame = 0
	//		}
	//		if demoGame.enemy2.direction == ENEMY_RIGHT {
	//			demoGame.enemy2.xLocNnemy += 3
	//			if demoGame.enemy2.xLocNnemy >= 750 {
	//				demoGame.enemy2.direction = ENEMY_LEFT
	//			}
	//		} else if demoGame.enemy2.direction == ENEMY_LEFT {
	//			demoGame.enemy2.xLocNnemy -= 3
	//			if demoGame.enemy2.xLocNnemy <= 570 {
	//				demoGame.enemy2.direction = ENEMY_RIGHT
	//			}
	//		}
	//	}
	//}
	//if demoGame.enemy3.alive == true {
	//	demoGame.enemy3.frameDelay += 1
	//	if demoGame.enemy3.frameDelay%ENEMY_FRAME_COUNT == 0 {
	//		demoGame.enemy3.frame += 1
	//		if demoGame.enemy3.frame >= ENEMY_FRAME_PER_SHEET {
	//			demoGame.enemy3.frame = 0
	//		}
	//		if demoGame.enemy3.direction == ENEMY_RIGHT {
	//			demoGame.enemy3.xLocNnemy += 3
	//			if demoGame.enemy3.xLocNnemy >= 750 {
	//				demoGame.enemy3.direction = ENEMY_LEFT
	//			}
	//		} else if demoGame.enemy3.direction == ENEMY_LEFT {
	//			demoGame.enemy3.xLocNnemy -= 3
	//			if demoGame.enemy3.xLocNnemy <= 570 {
	//				demoGame.enemy3.direction = ENEMY_RIGHT
	//			}
	//		}
	//	}
	//}
	//if demoGame.enemy4.alive == true {
	//	demoGame.enemy4.frameDelay += 1
	//	if demoGame.enemy4.frameDelay%ENEMY_FRAME_COUNT == 0 {
	//		demoGame.enemy4.frame += 1
	//		if demoGame.enemy4.frame >= ENEMY_FRAME_PER_SHEET {
	//			demoGame.enemy4.frame = 0
	//		}
	//		if demoGame.enemy4.direction == ENEMY_RIGHT {
	//			demoGame.enemy4.xLocNnemy += 3
	//			if demoGame.enemy4.xLocNnemy >= 350 {
	//				demoGame.enemy4.direction = ENEMY_LEFT
	//			}
	//		} else if demoGame.enemy4.direction == ENEMY_LEFT {
	//			demoGame.enemy4.xLocNnemy -= 3
	//			if demoGame.enemy4.xLocNnemy <= 200 {
	//				demoGame.enemy4.direction = ENEMY_RIGHT
	//			}
	//		}
	//	}
	//}

	// map switching
	if demoGame.playerXLoc >= 950 && demoGame.levels == 0 && demoGame.direction == RIGHT {
		gameMap, err := tiled.LoadFile(map2Path)
		if err != nil {
			fmt.Printf("error parsing map: %s", err.Error())
			return err
		}
		demoGame.levels = 1
		demoGame.playerXLoc = 0
		demoGame.playerYLoc = 128
		demoGame.level = gameMap
	}
	if demoGame.direction == LEFT && demoGame.levels == 1 && demoGame.playerXLoc <= 0 {
		gameMap, err := tiled.LoadFile(map1Path)
		if err != nil {
			fmt.Printf("error parsing map: %s", err.Error())
			return err
		}
		demoGame.levels = 0
		demoGame.playerXLoc = 950
		demoGame.playerYLoc = 448
		demoGame.level = gameMap
	}

	if demoGame.playerXLoc >= 950 && demoGame.levels == 1 && demoGame.direction == RIGHT {
		gameMap, err := tiled.LoadFile(map3Path)
		if err != nil {
			fmt.Printf("error parsing map: %s", err.Error())
			return err
		}
		demoGame.levels = 2
		demoGame.playerXLoc = 0
		demoGame.playerYLoc = 448
		demoGame.level = gameMap
	}
	if demoGame.direction == LEFT && demoGame.levels == 2 && demoGame.playerXLoc <= 0 {
		gameMap, err := tiled.LoadFile(map2Path)
		if err != nil {
			fmt.Printf("error parsing map: %s", err.Error())
			return err
		}
		demoGame.levels = 1
		demoGame.playerXLoc = 950
		demoGame.playerYLoc = 704
		demoGame.level = gameMap
	}

	demoGame.msg = false
	nonPCOne := image.Rect(demoGame.npc1.xLocNPC, demoGame.npc1.yLocNPC, demoGame.npc1.xLocNPC+NPC_FRAME_WIDTH, demoGame.npc1.yLocNPC+NPC_HEIGHT)
	playerChar := image.Rect(demoGame.playerXLoc, demoGame.playerYLoc, demoGame.playerXLoc+GUY_FRAME_WIDTH, demoGame.playerYLoc+GUY_HEIGHT)
	if playerChar.Overlaps(nonPCOne) {
		demoGame.msg = true
	} else {
		demoGame.msg = false
	}
	if demoGame.levels == 1 {
		enemy1 := image.Rect(demoGame.enemy1.xLocNnemy, demoGame.enemy1.yLocEnemy, demoGame.enemy1.xLocNnemy+ENEMY_FRAME_WIDTH, demoGame.enemy1.yLocEnemy+ENEMY_HEIGHT)
		if playerChar.Overlaps(enemy1) {
			demoGame.enemy1.alive = false
		}
		enemy2 := image.Rect(demoGame.enemy2.xLocNnemy, demoGame.enemy2.yLocEnemy, demoGame.enemy2.xLocNnemy+ENEMY_FRAME_WIDTH, demoGame.enemy2.yLocEnemy+ENEMY_HEIGHT)
		if playerChar.Overlaps(enemy2) {
			demoGame.enemy2.alive = false
		}
		enemy3 := image.Rect(demoGame.enemy3.xLocNnemy, demoGame.enemy3.yLocEnemy, demoGame.enemy3.xLocNnemy+ENEMY_FRAME_WIDTH, demoGame.enemy3.yLocEnemy+ENEMY_HEIGHT)
		if playerChar.Overlaps(enemy3) {
			demoGame.enemy3.alive = false
		}
		enemy4 := image.Rect(demoGame.enemy4.xLocNnemy, demoGame.enemy4.yLocEnemy, demoGame.enemy4.xLocNnemy+ENEMY_FRAME_WIDTH, demoGame.enemy4.yLocEnemy+ENEMY_HEIGHT)
		if playerChar.Overlaps(enemy4) {
			demoGame.enemy4.alive = false
		}
	}
	//enemy3 := image.Rect(demoGame.enemy3.xLocNnemy, demoGame.enemy3.yLocEnemy, demoGame.enemy3.xLocNnemy+ENEMY_FRAME_WIDTH, demoGame.enemy3.yLocEnemy+ENEMY_HEIGHT)
	//if playerChar.Overlaps(enemy3) {
	//	demoGame.enemy3.alive = false
	//}

	// tiles collision for level 0
	if demoGame.levels == 0 {
		for tileY := 0; tileY < demoGame.level.Height; tileY += 1 {
			for tileX := 0; tileX < demoGame.level.Width; tileX += 1 {
				tileID := demoGame.level.Layers[1].Tiles[tileY*demoGame.level.Width+tileX].ID

				if tileID == 3 || tileID == 10 {
					block4 := image.Rect(tileX*64, tileY*64, (tileX*64)+64, (tileY*64)+64)
					player := image.Rect(demoGame.playerXLoc, demoGame.playerYLoc, demoGame.playerXLoc+GUY_FRAME_WIDTH, demoGame.playerYLoc+GUY_HEIGHT)
					if player.Overlaps(block4) && demoGame.direction == DOWN {
						demoGame.playerYLoc -= 2
					} else if player.Overlaps(block4) && demoGame.direction == RIGHT {
						demoGame.playerXLoc -= 2
					} else if player.Overlaps(block4) && demoGame.direction == UP {
						demoGame.playerYLoc += 2
					} else if player.Overlaps(block4) && demoGame.direction == LEFT {
						demoGame.playerXLoc += 2
					}
				}
			}
		}
	} else if demoGame.levels == 1 || demoGame.levels == 2 {
		// collision for level 1 and 2
		for tileY := 0; tileY < demoGame.level.Height; tileY += 1 {
			for tileX := 0; tileX < demoGame.level.Width; tileX += 1 {
				tileID := demoGame.level.Layers[1].Tiles[tileY*demoGame.level.Width+tileX].ID

				if tileID == 1 || tileID == 8 {
					block4 := image.Rect(tileX*64, tileY*64, (tileX*64)+64, (tileY*64)+64)
					player := image.Rect(demoGame.playerXLoc, demoGame.playerYLoc, demoGame.playerXLoc+GUY_FRAME_WIDTH, demoGame.playerYLoc+GUY_HEIGHT)
					if player.Overlaps(block4) && demoGame.direction == DOWN {
						demoGame.playerYLoc -= 2
					} else if player.Overlaps(block4) && demoGame.direction == RIGHT {
						demoGame.playerXLoc -= 2
					} else if player.Overlaps(block4) && demoGame.direction == UP {
						demoGame.playerYLoc += 2
					} else if player.Overlaps(block4) && demoGame.direction == LEFT {
						demoGame.playerXLoc += 2
					}
				}
			}
		}
	}
	if demoGame.coin1.pickedUp == false && demoGame.levels == 1 {
		demoGame.coin1.frameDelay += 1
		if demoGame.coin1.frameDelay%FRAME_COUNT == 0 {
			demoGame.coin1.frame += 1
			if demoGame.coin1.frame >= COIN_FRAME_PER_SHEET {
				demoGame.coin1.frame = 0
			}
		}
	}
	if demoGame.coin2.pickedUp == false && demoGame.levels == 1 {
		demoGame.coin2.frameDelay += 1
		if demoGame.coin2.frameDelay%FRAME_COUNT == 0 {
			demoGame.coin2.frame += 1
			if demoGame.coin2.frame >= COIN_FRAME_PER_SHEET {
				demoGame.coin2.frame = 0
			}
		}
	}
	if demoGame.coin3.pickedUp == false && demoGame.levels == 1 {
		demoGame.coin3.frameDelay += 1
		if demoGame.coin3.frameDelay%FRAME_COUNT == 0 {
			demoGame.coin3.frame += 1
			if demoGame.coin3.frame >= COIN_FRAME_PER_SHEET {
				demoGame.coin3.frame = 0
			}
		}
	}
	if demoGame.coin4.pickedUp == false && demoGame.levels == 1 {
		demoGame.coin4.frameDelay += 1
		if demoGame.coin4.frameDelay%FRAME_COUNT == 0 {
			demoGame.coin4.frame += 1
			if demoGame.coin4.frame >= COIN_FRAME_PER_SHEET {
				demoGame.coin4.frame = 0
			}
		}
	}
	if demoGame.levels == 1 {
		startRow := int(demoGame.enemy3.yLocEnemy) / demoGame.level.TileHeight
		startCol := int(demoGame.enemy3.xLocNnemy) / demoGame.level.TileWidth
		startCell := demoGame.pathMap.Get(startCol, startRow)
		endCell := demoGame.pathMap.Get(demoGame.playerXLoc/demoGame.level.TileWidth, demoGame.playerYLoc/demoGame.level.TileHeight)
		demoGame.path = demoGame.pathMap.GetPathFromCells(startCell, endCell, false, false)
		if demoGame.path != nil {
			pathCell := demoGame.path.Current()
			if math.Abs(float64(pathCell.X*demoGame.level.TileWidth)-(float64(demoGame.enemy3.xLocNnemy))) <= 2 &&
				math.Abs(float64(pathCell.Y*demoGame.level.TileHeight)-(float64(demoGame.enemy3.yLocEnemy))) <= 2 { //if we are now on the tile we need to be on
				demoGame.path.Advance()
			}
			direction := 0
			if pathCell.X*demoGame.level.TileWidth > int(demoGame.enemy3.xLocNnemy) {
				direction = 1
			} else if pathCell.X*demoGame.level.TileWidth < int(demoGame.enemy3.xLocNnemy) {
				direction = -1
			}
			Ydirection := 0
			if pathCell.Y*demoGame.level.TileHeight > int(demoGame.enemy3.yLocEnemy) {
				Ydirection = 1
			} else if pathCell.Y*demoGame.level.TileHeight < int(demoGame.enemy3.yLocEnemy) {
				Ydirection = -1
			}
			demoGame.enemy3.xLocNnemy += direction * 2
			demoGame.enemy3.yLocEnemy += Ydirection * 2
		}
	}

	return nil
}

func (demoGame AnimatedSpriteDemo3) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	// drawing all the layers of the map
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
	// drawing the player
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(demoGame.playerXLoc), float64(demoGame.playerYLoc))
	screen.DrawImage(demoGame.spriteSheet.SubImage(image.Rect(demoGame.frame*GUY_FRAME_WIDTH,
		demoGame.direction*GUY_HEIGHT,
		demoGame.frame*GUY_FRAME_WIDTH+GUY_FRAME_WIDTH,
		demoGame.direction*GUY_HEIGHT+GUY_HEIGHT)).(*ebiten.Image), &drawOptions)
	// drawing NPC1 in level 0
	if demoGame.levels == 0 {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.npc1.xLocNPC), float64(demoGame.npc1.yLocNPC))
		screen.DrawImage(demoGame.npc1.sprite.SubImage(image.Rect(demoGame.npc1.frame*NPC_FRAME_WIDTH,
			demoGame.npc1.direction*NPC_HEIGHT,
			demoGame.npc1.frame*NPC_FRAME_WIDTH+NPC_FRAME_WIDTH,
			demoGame.npc1.direction*NPC_HEIGHT+NPC_HEIGHT)).(*ebiten.Image), &drawOptions)
	}
	// drawing the enemies in level 1
	if demoGame.levels == 1 {
		if demoGame.enemy1.alive == true {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(demoGame.enemy1.xLocNnemy), float64(demoGame.enemy1.yLocEnemy))
			screen.DrawImage(demoGame.enemy1.sprite.SubImage(image.Rect(demoGame.enemy1.frame*ENEMY_FRAME_WIDTH,
				demoGame.enemy1.direction*ENEMY_HEIGHT,
				demoGame.enemy1.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
				demoGame.enemy1.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)
		} else {
			demoGame.damage++
		}
		if demoGame.enemy2.alive == true {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(demoGame.enemy2.xLocNnemy), float64(demoGame.enemy2.yLocEnemy))
			screen.DrawImage(demoGame.enemy2.sprite.SubImage(image.Rect(demoGame.enemy2.frame*ENEMY_FRAME_WIDTH,
				demoGame.enemy2.direction*ENEMY_HEIGHT,
				demoGame.enemy2.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
				demoGame.enemy2.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)
		} else {
			demoGame.damage++
		}
		if demoGame.enemy3.alive == true {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(demoGame.enemy3.xLocNnemy), float64(demoGame.enemy3.yLocEnemy))
			screen.DrawImage(demoGame.enemy3.sprite.SubImage(image.Rect(demoGame.enemy3.frame*ENEMY_FRAME_WIDTH,
				demoGame.enemy3.direction*ENEMY_HEIGHT,
				demoGame.enemy3.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
				demoGame.enemy3.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)
		} else {
			demoGame.damage++
		}
		if demoGame.enemy4.alive == true {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(demoGame.enemy4.xLocNnemy), float64(demoGame.enemy4.yLocEnemy))
			screen.DrawImage(demoGame.enemy4.sprite.SubImage(image.Rect(demoGame.enemy4.frame*ENEMY_FRAME_WIDTH,
				demoGame.enemy4.direction*ENEMY_HEIGHT,
				demoGame.enemy4.frame*ENEMY_FRAME_WIDTH+ENEMY_FRAME_WIDTH,
				demoGame.enemy4.direction*ENEMY_HEIGHT+ENEMY_HEIGHT)).(*ebiten.Image), &drawOptions)
		} else {
			demoGame.damage++
		}
	}
	DrawCenteredText(screen, demoGame.textFont, fmt.Sprintf("Damage: %d", demoGame.damage), 65, 30)
	if demoGame.msg == true && demoGame.levels == 0 {
		DrawCenteredText(screen, basicfont.Face7x13, fmt.Sprintf("Hi Player, you should check the next room \n"+
			"head over to the end of that dirt bridge\n"+
			"if you find any enemies please try to kill\n"+
			"them. they are trying to take over"), 400, 200)
	}
	if demoGame.temp == true {
		for _, shot := range demoGame.shot {
			drawOptions.GeoM.Reset()
			drawOptions.GeoM.Translate(float64(shot.bulletXLoc), float64(shot.bulletYLoc))
			screen.DrawImage(shot.sprite, &drawOptions)
		}
	}
	if demoGame.levels == 1 && demoGame.enemy1.alive == false {
		demoGame.coin1.coinXLoc = demoGame.enemy1.xLocNnemy
		demoGame.coin1.coinYLoc = demoGame.enemy1.yLocEnemy
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.coin1.coinXLoc), float64(demoGame.coin1.coinYLoc))
		screen.DrawImage(demoGame.coin1.sprite.SubImage(image.Rect(demoGame.coin1.frame*COIN_FRAME_WIDTH,
			demoGame.coin1.direction*COIN_HEIGHT,
			demoGame.coin1.frame*COIN_FRAME_WIDTH+COIN_FRAME_WIDTH,
			demoGame.coin1.direction*COIN_HEIGHT+COIN_HEIGHT)).(*ebiten.Image), &drawOptions)
	}
	if demoGame.levels == 1 && demoGame.enemy2.alive == false {
		demoGame.coin2.coinXLoc = demoGame.enemy2.xLocNnemy
		demoGame.coin2.coinYLoc = demoGame.enemy2.yLocEnemy
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.coin2.coinXLoc), float64(demoGame.coin2.coinYLoc))
		screen.DrawImage(demoGame.coin2.sprite.SubImage(image.Rect(demoGame.coin2.frame*COIN_FRAME_WIDTH,
			demoGame.coin2.direction*COIN_HEIGHT,
			demoGame.coin2.frame*COIN_FRAME_WIDTH+COIN_FRAME_WIDTH,
			demoGame.coin2.direction*COIN_HEIGHT+COIN_HEIGHT)).(*ebiten.Image), &drawOptions)
	}
	if demoGame.levels == 1 && demoGame.enemy3.alive == false {
		demoGame.coin3.coinXLoc = demoGame.enemy3.xLocNnemy
		demoGame.coin3.coinYLoc = demoGame.enemy3.yLocEnemy
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.coin3.coinXLoc), float64(demoGame.coin3.coinYLoc))
		screen.DrawImage(demoGame.coin3.sprite.SubImage(image.Rect(demoGame.coin3.frame*COIN_FRAME_WIDTH,
			demoGame.coin3.direction*COIN_HEIGHT,
			demoGame.coin3.frame*COIN_FRAME_WIDTH+COIN_FRAME_WIDTH,
			demoGame.coin3.direction*COIN_HEIGHT+COIN_HEIGHT)).(*ebiten.Image), &drawOptions)
	}
	if demoGame.levels == 1 && demoGame.enemy4.alive == false {
		demoGame.coin4.coinXLoc = demoGame.enemy4.xLocNnemy
		demoGame.coin4.coinYLoc = demoGame.enemy4.yLocEnemy
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(demoGame.coin4.coinXLoc), float64(demoGame.coin4.coinYLoc))
		screen.DrawImage(demoGame.coin4.sprite.SubImage(image.Rect(demoGame.coin4.frame*COIN_FRAME_WIDTH,
			demoGame.coin4.direction*COIN_HEIGHT,
			demoGame.coin4.frame*COIN_FRAME_WIDTH+COIN_FRAME_WIDTH,
			demoGame.coin4.direction*COIN_HEIGHT+COIN_HEIGHT)).(*ebiten.Image), &drawOptions)
	}

}

func main() {

	gameMap, err := tiled.LoadFile(map1Path)
	pathMap := makeSearchMap(gameMap)
	searchablePathMap := paths.NewGridFromStringArrays(pathMap, gameMap.TileWidth, gameMap.TileHeight)
	searchablePathMap.SetWalkable('8', false)
	searchablePathMap.SetWalkable('1', false)
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
	enemyAnimation := LoadEmbeddedImage("", "earthEnemy.png")
	npc1Animation := LoadEmbeddedImage("", "npc2.png")
	shotAnimation := LoadEmbeddedImage("", "laser.png")
	coinAnimation := LoadEmbeddedImage("", "coin.png")

	drawFont := LoadScoreFont()
	allShots := make([]shots, 0, 20)

	oneLevelGame := AnimatedSpriteDemo3{
		pathFindingMap: pathMap,
		pathMap:        searchablePathMap,
		levels:         0,
		spriteSheet:    animationGuy,
		playerXLoc:     256,
		playerYLoc:     448,
		level:          gameMap,
		tileHash:       ebitenImageMap,
		damage:         1,
		textFont:       drawFont,
		shot:           allShots,
		bullet:         shotAnimation,
		coin1: coins{
			sprite:    coinAnimation,
			direction: COIN_RIGHT,
			pickedUp:  false,
		},
		coin2: coins{
			sprite:    coinAnimation,
			direction: COIN_RIGHT,
			pickedUp:  false,
		},
		coin3: coins{
			sprite:    coinAnimation,
			direction: COIN_RIGHT,
			pickedUp:  false,
		},
		coin4: coins{
			sprite:    coinAnimation,
			direction: COIN_RIGHT,
			pickedUp:  false,
		},
		enemy1: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 200,
			yLocEnemy: 608,
			direction: ENEMY_RIGHT,
			alive:     true,
		},
		enemy2: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 570,
			yLocEnemy: 384,
			direction: ENEMY_RIGHT,
			alive:     true,
		},
		enemy3: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 570,
			yLocEnemy: 300,
			//direction: ENEMY_LEFT,
			alive: true,
		},
		enemy4: enemy{
			sprite:    enemyAnimation,
			xLocNnemy: 200,
			yLocEnemy: 384,
			direction: ENEMY_LEFT,
			alive:     true,
		},
		npc1: npc{
			sprite:    npc1Animation,
			xLocNPC:   384,
			yLocNPC:   236,
			direction: NPC_DOWN,
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
func LoadScoreFont() font.Face {
	//originally inspired by https://www.fatoldyeti.com/posts/roguelike16/
	trueTypeFont, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		fmt.Println("Error loading font for score:", err)
	}
	fontFace, err := opentype.NewFace(trueTypeFont, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println("Error loading font of correct size for score:", err)
	}
	return fontFace
}
func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy int) {
	//from https://github.com/sedyh/ebitengine-cheatsheet

	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, colornames.Black)
}

func makeSearchMap(tiledMap *tiled.Map) []string {
	mapAsStringSlice := make([]string, 0, tiledMap.Height) //each row will be its own string
	row := strings.Builder{}
	for position, tile := range tiledMap.Layers[0].Tiles {
		if position%tiledMap.Width == 0 && position > 0 { // we get the 2d array as an unrolled one-d array
			mapAsStringSlice = append(mapAsStringSlice, row.String())
			row = strings.Builder{}
		}
		row.WriteString(fmt.Sprintf("%d", tile.ID))
	}
	mapAsStringSlice = append(mapAsStringSlice, row.String())
	return mapAsStringSlice
}

//func check(demoGame AnimatedSpriteDemo3) {
//
//	startRow := int(demoGame.enemy1.yLocEnemy) / demoGame.level.TileHeight
//	startCol := int(demoGame.enemy1.xLocNnemy) / demoGame.level.TileWidth
//	startCell := demoGame.pathMap.Get(startCol, startRow)
//	endCell := demoGame.pathMap.Get(demoGame.playerXLoc, demoGame.playerYLoc)
//	demoGame.path = demoGame.pathMap.GetPathFromCells(startCell, endCell, false, false)
//}
