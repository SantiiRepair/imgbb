package processor

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
)

// ConvertToWebP takes image bytes and returns its WebP version
func ConvertToWebP(inputData []byte, quality float32) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(inputData))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	options := &webp.Options{Lossless: false, Quality: quality}
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
