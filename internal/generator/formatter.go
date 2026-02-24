package generator

import (
	"diagram-gen/internal/model"
)

// Formatter defines the interface for diagram generators.
type Formatter interface {
	Generate(diagram *model.Diagram) ([]byte, error)
	Format() string
}
