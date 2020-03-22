package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	png "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/codename-pyoko/texture-packer/packer"
	"github.com/sirupsen/logrus"
)

func saveImage(path string, image image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}

	if err := png.Encode(f, image); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

func saveAtlas(path string, atlas packer.Atlas) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}

	if err := json.NewEncoder(f).Encode(&atlas); err != nil {
		return fmt.Errorf("failed to encode atlas to json: %w", err)
	}

	return nil
}

func main() {
	root := flag.String("root", "", "the root path to find images for texture packing")
	out := flag.String("out", "", "the output file - png")
	atlas := flag.String("atlas", "", "the atlas output file - json")

	flag.Parse()

	if *root == "" {
		logrus.Fatalf("argument 'root' missing")
	}

	if *out == "" {
		logrus.Fatalf("argument 'out' is missing")
	}

	p := packer.NewPacker()

	logrus.Infof("walking %s", *root)
	filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".png") {
			if err := p.UseImage(path); err != nil {
				logrus.Errorf("could not use image: %v", err)
			}
		}

		return err
	})

	if p.Len() == 0 {
		logrus.Infof("no images found - exiting")
	}

	logrus.Infof("found %d images", p.Len())

	image, err := p.Pack(packer.PackerOptions{})
	if err != nil {
		logrus.Fatalf("failed to pack: %v", err)
	}

	if err := saveImage(*out, image); err != nil {
		logrus.Fatalf("failed to save image: %v", err)
	}

	if *atlas != "" {
		if err := saveAtlas(*atlas, p.Atlas()); err != nil {
			logrus.Fatalf("failed to save atlas: %v", err)
		}
	}

}
