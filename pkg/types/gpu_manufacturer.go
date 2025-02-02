package types

type GPUManufacturer string

func (g GPUManufacturer) String() string {
	return string(g)
}

const (
	AMD    GPUManufacturer = "AMD"
	Intel  GPUManufacturer = "Intel"
	NVIDIA GPUManufacturer = "NVIDIA"
)
