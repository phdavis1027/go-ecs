package entity

func (ecs *ECS) CreateTile(layer int) (Entity, error) {
	tileEntity, err := ecs.CreateEntityOfType(TILE)
	if err != nil {
		return -1, err
	}

	error := ecs.AttachLayerComponent(tileEntity, TILE)
	if error != nil {
		return -1, error
	}

	return tileEntity, nil
}
