package resource

import (
	"encoding/json"
	"runtime"

	"github.com/EngineerBetter/control-tower/resource/internal/file"
	"github.com/EngineerBetter/control-tower/util/bincache"
)

//go:generate go-bindata -o internal/file/file.go -ignore (\.go$)|(\.git)|(bosh/assets) -nometadata -pkg file ../../control-tower-ops/... assets/...

// Resource safely exposes the json parameters of a resource
type Resource struct {
	URL     string `json:"url"`
	Version string `json:"version"`
	SHA1    string `json:"sha1"`
}

var resources map[string]Resource

var (
	// DirectorManifest statically defines director-manifest.yml contents
	DirectorManifest = mustAssetString("assets/manifest.yml")

	// AWSDirectorCloudConfig statically defines aws cloud-config.yml
	AWSDirectorCloudConfig = mustAssetString("assets/aws/cloud-config.yml")
	// AWSCPIOps statically defines aws-cpi.yml contents
	AWSCPIOps = mustAssetString("assets/aws/cpi.yml")
	//GCPJumpboxUserOps statically defines gcp jumpbox-user.yml
	GCPJumpboxUserOps = mustAssetString("assets/gcp/jumpbox-user.yml")
	// GCPDirectorCloudConfig statically defines gcp cloud-config.yml
	GCPDirectorCloudConfig = mustAssetString("assets/gcp/cloud-config.yml")
	// GCPCPIOps statically defines gcp-cpi.yml contents
	GCPCPIOps = mustAssetString("assets/gcp/cpi.yml")
	// GCPExternalIPOps statically defines external-ip.yml contents
	GCPExternalIPOps = mustAssetString("assets/gcp/external-ip.yml")
	// GCPDirectorCustomOps statically defines custom-ops.yml contents
	GCPDirectorCustomOps = mustAssetString("assets/gcp/custom-ops.yml")

	// AWSTerraformConfig holds the terraform conf for AWS
	AWSTerraformConfig = mustAssetString("assets/aws/infrastructure.tf")

	// GCPTerraformConfig holds the terraform conf for GCP
	GCPTerraformConfig = mustAssetString("assets/gcp/infrastructure.tf")

	// ExternalIPOps statically defines external-ip.yml contents
	ExternalIPOps = mustAssetString("assets/external-ip.yml")
	// AWSDirectorCustomOps statically defines custom-ops.yml contents
	AWSDirectorCustomOps = mustAssetString("assets/aws/custom-ops.yml")

	// AWSReleaseVersions carries all versions of releases
	AWSReleaseVersions = mustAssetString("../../control-tower-ops/ops/versions-aws.json")

	// GCPReleaseVersions carries all versions of releases
	GCPReleaseVersions = mustAssetString("../../control-tower-ops/ops/versions-gcp.json")

	// AddNewCa carries the ops file that adds a new CA required for cert rotation
	AddNewCa = mustAssetString("assets/maintenance/add-new-ca.yml")

	// RemoveOldCa carries the ops file that removes the old CA required for cert rotation
	RemoveOldCa = mustAssetString("assets/maintenance/remove-old-ca.yml")

	// CleanupCerts moves renewed values of certs to old keys in director vars store
	CleanupCerts = mustAssetString("assets/maintenance/cleanup-certs.yml")
)

// NOTE(px) remove this in a later version of github.com/mattn/go-bindata
func mustAssetString(name string) string {
	return string(file.MustAsset(name))
}

func get(name string) Resource {
	r, ok := resources[name]
	if !ok {
		panic("resource " + name + " not found")
	}
	return r
}

func AWSCPI() Resource {
	return get("cpi")
}

func AWSStemcell() Resource {
	return get("stemcell")
}

func BOSHRelease() Resource {
	return get("bosh")
}

func BPMRelease() Resource {
	return get("bpm")
}

func init() {
	p := file.MustAsset("../../control-tower-ops/createenv-dependencies-and-cli-versions.json")
	err := json.Unmarshal(p, &resources)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(p, &binaries)
	if err != nil {
		panic(err)
	}
}

var binaries map[string]binaryPaths

type binaryPaths struct {
	Mac   string `json:"mac"`
	Linux string `json:"linux"`
}

func (p binaryPaths) path() string {
	switch runtime.GOOS {
	case "darwin":
		return p.Mac
	case "linux":
		return p.Linux
	default:
		panic("OS not supported")
	}
}

// DownloadBOSHCLI returns the path of the downloaded bosh-cli
func DownloadBOSHCLI() (string, error) {
	p := binaries["bosh-cli"].path()
	return bincache.Download(p)
}

// DownloadTerraformCLI returns the path of the downloaded terraform-cli
func DownloadTerraformCLI() (string, error) {
	p := binaries["terraform"].path()
	return bincache.Download(p)
}
