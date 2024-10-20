package main

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/phdavis1027/goecs/entity"
)

func printA(ecs *entity.ECS, queries []entity.EntityType, entities []roaring64.Bitmap, queriesMut []entity.EntityType, entitiesMut []roaring64.Bitmap) {
	for i := 0; i < 100; i++ {
		// Do something with the entities
		println("A")
	}
}

func printB(ecs *entity.ECS, queries []entity.EntityType, entities []roaring64.Bitmap, queriesMut []entity.EntityType, entitiesMut []roaring64.Bitmap) {
	for i := 0; i < 100; i++ {
		// Do something with the entities
		println("B")
	}
}

func printC(ecs *entity.ECS, queries []entity.EntityType, entities []roaring64.Bitmap, queriesMut []entity.EntityType, entitiesMut []roaring64.Bitmap) {
	for i := 0; i < 100; i++ {
		// Do something with the entities
		println("C")
	}
}

func main() {
	ecs := entity.CreateEcsOfCapacity(10)

	ecs.RegisterSystem("A", printA)
	ecs.RegisterSystem("B", printB)
	ecs.RegisterSystem("C", printC)

	ecs.RegisterQueries("A", entity.Zero)
	ecs.RegisterMutQueries("A", entity.One)

	ecs.RegisterQueries("B", entity.One)
	ecs.RegisterMutQueries("B", entity.Two)

	ecs.RegisterQueries("C", entity.Two)

	ecs.CompileSchedule(true)
	ecs.RunSchedule()
}
