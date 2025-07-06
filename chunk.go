package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

const ChunkSize = 32

type Voxel struct {
	position mgl32.Vec3
	albedo   mgl32.Vec3
}

type Chunk struct {
	voxels [ChunkSize * ChunkSize * ChunkSize]Voxel
}

func index(x, y, z int) int {
	return x + (y * ChunkSize) + (z * ChunkSize * ChunkSize)
}

func GenerateChunk() Chunk {
	var chunk Chunk

	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			for z := 0; z < ChunkSize; z++ {
				idx := index(x, y, z)
				chunk.voxels[idx] = Voxel{
					position: mgl32.Vec3{float32(x), float32(y), float32(z)},
					albedo:   mgl32.Vec3{1, 1, 1},
				}
			}
		}
	}

	return chunk
}
