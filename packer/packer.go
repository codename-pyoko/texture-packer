package packer

import (
	"fmt"
	"image"
	"os"
	"sort"

	"github.com/codename-pyoko/texture-packer/binpack"
	"github.com/sirupsen/logrus"
)

type placedImage struct {
	x   int
	y   int
	img image.Image
}

type imageEntry struct {
	path  string
	image image.Image
}

type PackerOptions struct {
}

type Packer struct {
	images       []imageEntry
	atlas        Atlas
	placedImages []placedImage
}

func NewPacker() *Packer {
	return &Packer{
		atlas: NewAtlas(),
	}
}

func (p *Packer) Len() int {
	return len(p.images)
}

func (p *Packer) Size(n int) (int, int) {
	b := p.images[n].image.Bounds()
	return b.Dx(), b.Dy()
}

func (p *Packer) Place(n, x, y int) {
	logrus.Debugf("placing image %d at %d,%d", n, x, y)
	p.atlas.Add(x, y, p.images[n].path, p.images[n].image)
	p.placedImages = append(p.placedImages, placedImage{x, y, p.images[n].image})
}

func (p *Packer) UseImage(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}

	img, format, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode image %s: %w", path, err)
	}

	logrus.Debugf("decoded image %s (%s): %v", path, format, img.Bounds())

	p.images = append(p.images, imageEntry{path, img})

	return nil
}

func (p *Packer) interpose(target *image.NRGBA, img image.Image, targetX, targetY int) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			target.Set(targetX+x, targetY+y, img.At(x, y))
		}
	}
}

func intMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (p *Packer) Pack(opts PackerOptions) (image.Image, error) {
	logrus.Debugf("sorting images")
	sort.Slice(p.images, func(i, j int) bool {
		lhs, rhs := p.images[i].image.Bounds(), p.images[j].image.Bounds()
		return intMax(rhs.Dx(), rhs.Dy()) < intMax(lhs.Dx(), lhs.Dy())
	})

	logrus.Debugf("packing images")
	w, h := binpack.Pack(p)

	logrus.Debugf("creating image if dimensions %d,%d", w, h)
	packed := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})

	logrus.Debugf("interposing images")
	for _, pimg := range p.placedImages {
		p.interpose(packed, pimg.img, pimg.x, pimg.y)
	}

	return packed, nil
}

func (p *Packer) Atlas() Atlas {
	return p.atlas
}
