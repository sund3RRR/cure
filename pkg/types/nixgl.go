package types

type NixGL string

func (n NixGL) String() string {
	return string(n)
}

const (
	NixGLAuto NixGL = "auto"
	NixGLAll  NixGL = "all"
	NixGLNone NixGL = "no"
)

var NixGLValues = []NixGL{NixGLAuto, NixGLAll, NixGLNone}

type NixGLPackage string

func (n NixGLPackage) String() string {
	return string(n)
}

const (
	NixGLPackageAuto NixGLPackage = "auto"
	NixGLMesa        NixGLPackage = "mesa"
	NixVulkanMesa    NixGLPackage = "mesa-vl"
	NixGLNvidia      NixGLPackage = "nvidia"
	NixVulkanNvidia  NixGLPackage = "nvidia-vl"
)

var NixGLPackageValues = []NixGLPackage{NixGLPackageAuto, NixGLMesa, NixVulkanMesa, NixGLNvidia, NixVulkanNvidia}
