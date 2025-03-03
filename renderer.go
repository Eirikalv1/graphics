package main

import (
	"unsafe"

	"github.com/rajveermalviya/go-webgpu/wgpu"

	_ "embed"
)

//go:embed shader.wgsl
var shaderSource string

type Vertex struct {
	pos   [2]float32
	color [2]float32
}

var vertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepMode_Vertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Format:         wgpu.VertexFormat_Float32x2,
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         wgpu.VertexFormat_Float32x2,
			Offset:         4 * 2,
			ShaderLocation: 1,
		},
	},
}

var vertexData = [...]Vertex{
	{
		pos:   [2]float32{-1.0, 1.0},
		color: [2]float32{0.0, 1.0},
	},
	{
		pos:   [2]float32{1.0, 1.0},
		color: [2]float32{1.0, 1.0},
	},
	{
		pos:   [2]float32{1.0, -1.0},
		color: [2]float32{1.0, 0.0},
	},
	{
		pos:   [2]float32{-1.0, -1.0},
		color: [2]float32{0.0, 0.0},
	},
}

var indexData = [...]uint16{
	2, 1, 0,
	3, 2, 0,
}

type Renderer struct {
	pipeline              *wgpu.RenderPipeline
	pipelineLayout        *wgpu.PipelineLayout
	shader                *wgpu.ShaderModule
	vertexBuffer          *wgpu.Buffer
	indexBuffer           *wgpu.Buffer
	cameraBuffer          *wgpu.Buffer
	cameraBindGroup       *wgpu.BindGroup
	cameraBindGroupLayout *wgpu.BindGroupLayout
}

func InitRenderer(ctx *GpuContext, camera *Camera) (renderer *Renderer, err error) {
	defer func() {
		if err != nil {
			renderer.Destroy()
			renderer = nil
		}
	}()

	renderer = &Renderer{}

	renderer.shader, err = ctx.device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label:          "shader.wgsl",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shaderSource},
	})
	if err != nil {
		return renderer, err
	}

	renderer.vertexBuffer, err = ctx.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexData[:]),
		Usage:    wgpu.BufferUsage_Vertex,
	})
	if err != nil {
		return renderer, err
	}

	renderer.indexBuffer, err = ctx.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Index Buffer",
		Contents: wgpu.ToBytes(indexData[:]),
		Usage:    wgpu.BufferUsage_Index,
	})
	if err != nil {
		return renderer, err
	}

	renderer.cameraBuffer, err = ctx.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: "Camera Buffer",
		Size:  uint64(unsafe.Sizeof(camera.uniform)),
		Usage: wgpu.BufferUsage_Uniform | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return renderer, err
	}

	renderer.cameraBindGroupLayout, err = ctx.device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		Label: "Camera Bind Group Layout",
		Entries: []wgpu.BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: wgpu.ShaderStage_Vertex | wgpu.ShaderStage_Fragment,
				Buffer: wgpu.BufferBindingLayout{
					Type: wgpu.BufferBindingType_Uniform,
				},
			},
		},
	})
	if err != nil {
		return renderer, err
	}

	renderer.cameraBindGroup, err = ctx.device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  "Camera Bind Group",
		Layout: renderer.cameraBindGroupLayout,
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0,
				Buffer:  renderer.cameraBuffer,
				Offset:  0,
				Size:    uint64(unsafe.Sizeof(camera.uniform)),
			},
		},
	})
	if err != nil {
		return renderer, err
	}

	renderer.pipelineLayout, err = ctx.device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:            "Pipeline Layout",
		BindGroupLayouts: []*wgpu.BindGroupLayout{renderer.cameraBindGroupLayout},
	})
	if err != nil {
		return renderer, err
	}

	renderer.pipeline, err = ctx.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "Main Render Pipeline",
		Layout: renderer.pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     renderer.shader,
			EntryPoint: "vs_main",
			Buffers:    []wgpu.VertexBufferLayout{vertexBufferLayout},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:         wgpu.PrimitiveTopology_TriangleList,
			StripIndexFormat: wgpu.IndexFormat_Undefined,
			FrontFace:        wgpu.FrontFace_CCW,
			CullMode:         wgpu.CullMode_Back,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     renderer.shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    ctx.config.Format,
					Blend:     &wgpu.BlendState_Replace,
					WriteMask: wgpu.ColorWriteMask_All,
				},
			},
		},
	})
	if err != nil {
		return renderer, err
	}

	return renderer, nil
}

func (renderer *Renderer) Destroy() {
	if renderer.pipeline != nil {
		renderer.pipeline.Release()
	}

	if renderer.pipelineLayout != nil {
		renderer.pipelineLayout.Release()
	}

	if renderer.shader != nil {
		renderer.shader.Release()
	}

	if renderer.vertexBuffer != nil {
		renderer.vertexBuffer.Release()
	}

	if renderer.indexBuffer != nil {
		renderer.indexBuffer.Release()
	}

	if renderer.cameraBuffer != nil {
		renderer.cameraBuffer.Release()
	}

	if renderer.cameraBindGroup != nil {
		renderer.cameraBindGroup.Release()
	}

	if renderer.cameraBindGroupLayout != nil {
		renderer.cameraBindGroupLayout.Release()
	}
}

func (renderer *Renderer) Render(ctx *GpuContext, camera *Camera) error {
	nextTexture, err := ctx.swapChain.GetCurrentTextureView()
	if err != nil {
		return err
	}
	defer nextTexture.Release()

	encoder, err := ctx.device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Command Encoder",
	})
	if err != nil {
		return err
	}
	defer encoder.Release()

	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:       nextTexture,
				LoadOp:     wgpu.LoadOp_Clear,
				StoreOp:    wgpu.StoreOp_Store,
				ClearValue: wgpu.Color_Black,
			},
		},
	})
	defer renderPass.Release()

	cameraData := unsafe.Slice((*byte)(unsafe.Pointer(&camera.uniform)), unsafe.Sizeof(camera.uniform))
	ctx.queue.WriteBuffer(renderer.cameraBuffer, 0, cameraData)

	renderPass.SetPipeline(renderer.pipeline)
	renderPass.SetBindGroup(0, renderer.cameraBindGroup, nil)
	renderPass.SetVertexBuffer(0, renderer.vertexBuffer, 0, wgpu.WholeSize)
	renderPass.SetIndexBuffer(renderer.indexBuffer, wgpu.IndexFormat_Uint16, 0, wgpu.WholeSize)
	renderPass.DrawIndexed(uint32(len(indexData)), 1, 0, 0, 0)
	renderPass.End()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		return err
	}
	defer cmdBuffer.Release()

	ctx.queue.Submit(cmdBuffer)
	ctx.swapChain.Present()

	return nil
}
