package goecs

type Optional[T any] struct {
  Val         T
  IsSome      bool
}
