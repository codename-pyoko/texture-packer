# Texture Packer

CLI tool for packing textures in deeply nested folder structures.

The resulting texture and atlas are usable with [Phaser](https://github.com/photonstorm/phaser).

Only supports png input and output for now, but can easily be extended if necessary.

## Usage

```
texture-packer -root root-folder/ -out packed-texture.png -atlas packed-texture.json
```

Flag  | Description | Required
---|---|---
root  | The root folder to be traversed for `png` files to add to the packed texture | Yes
out   | The png file output for the packed texture | Yes
atlas | The JSON formatted atlas describing where each texture ended up | No
