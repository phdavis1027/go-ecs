package tile

import (
	"github.com/phdavis1027/goecs/entity"
	"github.com/phdavis1027/goecs/util"
)

type TileManager struct {
  Width  int
  Height int
  Tiles  []entity.Entity
}

func NewTileManager(width, height int) *TileManager {
  return &TileManager{
    Width:  width,
    Height: height,
    Tiles:  make([]entity.Entity, 0),
  }
}

func (tm *TileManager) AddTile(ecs *entity.ECS) error {
  tile, err := ecs.CreateEntityOfType(entity.Tile) 
  if err != nil {
    return err
  }

  err = ecs.AttachRectComponent(tile, entity.Tile, entity.NewRect(util.Vec2{}, 0, 0))
  if err != nil {
    return err
  }

  err = ecs.AttachLayerComponent(tile, entity.Tile, -1)
  if err != nil {
    return err
  }

  tm.Tiles = append(tm.Tiles, tile)

  return nil
}
