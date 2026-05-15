package imagenormalization

import (
	"bytes"
	"fmt"

	"github.com/disintegration/imaging"
)

type Normalization struct{}

func NewNormalizationService() *Normalization {
	return &Normalization{}
}

func (n Normalization) Normalize(imageBytes []byte) ([]byte, error) {
	src, err := imaging.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	dst := imaging.Grayscale(src)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, dst, imaging.PNG); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}
