package main

import (
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/msm"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

func main() {
	// Load backend using env path
	runtime.LoadBackendFromEnvOrDefault()
	// Set Cuda device to perform
	device := runtime.CreateDevice("CUDA", 0)
	runtime.SetDevice(&device)

	// Obtain the default MSM configuration.
	cfg := core.GetDefaultMSMConfig()

	// Define the size of the problem, here 2^18.
	size := 1 << 18

	// Generate scalars and points for the MSM operation.
	scalars := bn254.GenerateScalars(size)
	points := bn254.GenerateAffinePoints(size)

	// Create a CUDA stream for asynchronous operations.
	stream, _ := runtime.CreateStream()
	var p bn254.Projective

	// Allocate memory on the device for the result of the MSM operation.
	var out core.DeviceSlice
	_, e := out.MallocAsync(p.Size(), 1, stream)

	if e != runtime.Success {
		panic(e)
	}

	// Set the CUDA stream in the MSM configuration.
	cfg.StreamHandle = stream
	cfg.IsAsync = true

	// Perform the MSM operation.
	e = msm.Msm(scalars, points, &cfg, out)

	if e != runtime.Success {
		panic(e)
	}

	// Allocate host memory for the results and copy the results from the device.
	outHost := make(core.HostSlice[bn254.Projective], 1)
	runtime.SynchronizeStream(stream)
	runtime.DestroyStream(stream)
	outHost.CopyFromDevice(&out)

	// Free the device memory allocated for the results.
	out.Free()
}