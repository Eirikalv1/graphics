package main

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type GpuContext struct {
	instance  *wgpu.Instance
	device    *wgpu.Device
	queue     *wgpu.Queue
	surface   *wgpu.Surface
	swapChain *wgpu.SwapChain
	config    *wgpu.SwapChainDescriptor
}

func InitGpuContext(window *glfw.Window) (ctx *GpuContext, err error) {
	defer func() {
		if err != nil {
			ctx.Destroy()
			ctx = nil
		}
	}()

	ctx = &GpuContext{}

	ctx.instance = wgpu.CreateInstance(nil)

	ctx.surface = ctx.instance.CreateSurface(&wgpu.SurfaceDescriptor{
		XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
			Display: unsafe.Pointer(glfw.GetX11Display()),
			Window:  uint32(window.GetX11Window()),
		},
	})

	adapter, err := ctx.instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: ctx.surface,
	})
	if err != nil {
		return ctx, err
	}
	defer adapter.Release()

	ctx.device, err = adapter.RequestDevice(nil)
	if err != nil {
		return ctx, err
	}
	ctx.queue = ctx.device.GetQueue()

	surfaceCaps := ctx.surface.GetCapabilities(adapter)

	width, height := window.GetSize()
	ctx.config = &wgpu.SwapChainDescriptor{
		Usage:       wgpu.TextureUsage_RenderAttachment,
		Format:      surfaceCaps.Formats[1],
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentMode_Fifo,
		AlphaMode:   surfaceCaps.AlphaModes[0],
	}

	ctx.swapChain, err = ctx.device.CreateSwapChain(ctx.surface, ctx.config)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (ctx *GpuContext) Destroy() {
	if ctx.swapChain != nil {
		ctx.swapChain.Release()
		ctx.swapChain = nil
	}
	if ctx.config != nil {
		ctx.config = nil
	}
	if ctx.queue != nil {
		ctx.queue.Release()
		ctx.queue = nil
	}
	if ctx.device != nil {
		ctx.device.Release()
		ctx.device = nil
	}
	if ctx.surface != nil {
		ctx.surface.Release()
		ctx.surface = nil
	}
	if ctx.instance != nil {
		ctx.instance.Release()
		ctx.instance = nil
	}
}

func (ctx *GpuContext) Resize(width, height int) {
	if width > 0 && height > 0 {
		ctx.config.Width = uint32(width)
		ctx.config.Height = uint32(height)

		if ctx.swapChain != nil {
			ctx.swapChain.Release()
		}

		var err error
		ctx.swapChain, err = ctx.device.CreateSwapChain(ctx.surface, ctx.config)
		if err != nil {
			panic(err)
		}
	}
}
