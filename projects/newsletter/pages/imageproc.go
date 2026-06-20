package pages

import (
	"bytes"
	"image"
	_ "image/gif" // register GIF decoding
	"image/jpeg"
	_ "image/png" // register PNG decoding
	"mime/multipart"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // register WebP decoding
)

// maxUploadImageSize is the largest original file we accept before
// downscaling — generous, since the stored result is always re-encoded
// down to maxImageDimension at processedImageQuality.
const maxUploadImageSize = 25 << 20 // 25MB

const (
	maxImageDimension    = 1080
	processedImageQuality = 80
)

// processUploadedImage decodes a multipart file upload, downscales it so
// neither side exceeds maxImageDimension (skipped if already smaller), and
// re-encodes it as a JPEG at processedImageQuality. This keeps stored answer
// images small regardless of what the original phone/camera photo weighed.
// Non-image or undecodable uploads (shouldn't happen given the answer_images
// MIME whitelist, but FormFile content isn't re-verified here) are returned
// as an error so the caller can surface it instead of silently storing junk.
func processUploadedImage(file multipart.File) (data []byte, filename string, err error) {
	src, _, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w > maxImageDimension || h > maxImageDimension {
		if w >= h {
			h = h * maxImageDimension / w
			w = maxImageDimension
		} else {
			w = w * maxImageDimension / h
			h = maxImageDimension
		}
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
		src = dst
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, src, &jpeg.Options{Quality: processedImageQuality}); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "answer.jpg", nil
}
