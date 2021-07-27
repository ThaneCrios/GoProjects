package functions

import (
	guuid "github.com/google/uuid"
)

type IDs string
type Structs string

type Functions struct {
	IDs IDs
}

func (ids *IDs) GenUUID() string {
	id := guuid.New()
	return id.String()
}

func (ids *IDs) SliceUUID(uuid string) string {
	return uuid[:5]
}
