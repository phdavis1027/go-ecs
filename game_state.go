package goecs

type GIndex struct {
  index int;
  generation int;
}

type GArrayEntry[T any] struct {
  value T;
  generation int;
}

type GArray[T any] struct {
  data []GArrayEntry[T];
}

func (a *GArray[T]) Get(index GIndex) Optional[T] {
  entry := a.data[index.index];

  if entry.generation == index.generation {
    return Optional[T]{some: true, value: entry.value};
  }

  return Optional[T]{some: false};
}

func (a *GArray[T]) Set(index GIndex, value T) bool {
  if (index.generation > a.data[index.index].generation) {
    a.data[index.index] = GArrayEntry[T]{value: value, generation: index.generation};

    return true;
  }

  return false;
}

type Optional[T any] struct {
  some bool;
  value T;
}
