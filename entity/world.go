package entity

type LevelEntity = Entity 

func (ecs *ECS) CreateLevel(title string) (LevelEntity, error) {
	level, err := ecs.CreateEntityOfType(LEVEL)
	if err != nil {
		return level, err
	}

	// Right now, Levels only have tilemap components
	ecs.AttachTileMapComponent(level, TILE)

	return level, nil
}
