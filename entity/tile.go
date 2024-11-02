package entity

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

type TileEntity = Entity

type TileMapComponent struct {
	Height, Width int
	GapX          int
	GapY		  int
	TileHeight	  int
	TileWidth	  int
	Tiles         []TileEntity
}

func NewTileMapComponent(width, height int) *TileMapComponent {
	return &TileMapComponent{
		Width:  width,
		Height: height,
		GapX:   0,
		GapY:   0,
		Tiles:  make([]TileEntity, width*height),
	}
}

func (self *TileMapComponent) GetTile(x, y int) (error, TileEntity) {
	if x < 0 || x >= self.Width || y < 0 || y >= self.Height {
		return errors.New("Tile out of bounds"), 0
	}
	return nil, self.Tiles[x+y*self.Width]
}

func (self *TileMapComponent) SetTile(x, y int, tile TileEntity) error {
	if x < 0 || x >= self.Width || y < 0 || y >= self.Height {
		return errors.New("Tile out of bounds")
	}

	self.Tiles[x+y*self.Width] = tile

	return nil
}

func (self *ECS) CreateTile() (TileEntity, error) {
	tile , err := self.CreateEntityOfType(TILE)
	if err != nil {
		return 0, err
	}

	err = self.AttachLayerComponent(tile, TILE)	
	if err != nil {
		return 0, err
	}

	err = self.AttachRenderableQuadComponent(tile, TILE)
	if err != nil {
		return 0, err
	}

	err = self.AttachPositionComponent(tile, TILE)
	if err != nil {
		return 0, err
	}

	return tile, nil
}

func (self *ECS) LoadLevel(level LevelEntity) error {
	tileMap, err := self.GetTileMapComponent(level)
	if err != nil {
		return err
	}

	for y := 0; y < tileMap.Height; y++ {
		for x := 0; x < tileMap.Width; x++ {
			tile, err := self.CreateTile()
			if err != nil {
				return err
			}

			// Leave tiles with layer = 0 

			// Set the tile's position
			position, err := self.GetPositionComponent(tile)
			if err != nil {
				return err
			}
			*position  = PositionComponent(mgl32.Vec3{
				float32(x * tileMap.TileWidth + tileMap.GapX),
				float32(y * tileMap.TileHeight + tileMap.GapY),
				float32(0),
			})

			// Set the tile's renderableQuadComponent
			rQuad, err := self.GetRenderableQuadComponent(tile)
			if err != nil {
				return err
			}

			*rQuad = RenderableQuadComponent{
				Width:  float32(tileMap.TileWidth),
				Height: float32(tileMap.TileHeight),
				Color: mgl32.Vec4{0.5, 0.25, 0.175, 1},
			}

			err = tileMap.SetTile(x, y, tile)
			if err != nil {
				return err
			}
		}
	}
}
