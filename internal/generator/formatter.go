package generator

import (
	"diagram-gen/internal/model"
)

type Formatter interface {
	Generate(diagram *model.Diagram) ([]byte, error)
	Format() string
}
