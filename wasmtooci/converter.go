package wasmtooci

import (
	"os"
	"path/filepath"
)

const (
	mediaTypeIndex = "application/vnd.oci.image.manifest.v1+json"
	mediaTypeConfig = "application/vnd.wasm.config.v0+json"
	mediaTypeLayer = "application/wasm"
	mediaTypeManifest = "application/vnd.oci.image.manifest.v1+json"
	manifestSchemaVersion = 2
)

func NewConverter() (*Converter, error) {
	return &Converter{}, nil
}

// TODO: Store in /tmp then archive and compress
func (*Converter) Convert(wasmFilePath string) (string,int64, error) {
	// Create temp dir for OCI artifact
	// TODO: create rand temp root dir in /tmp
	blobDir := filepath.Join(os.TempDir(), "blobs", "sha256")
	os.MkdirAll(blobDir, os.ModePerm)
	
	// TODO: func CreateOCIWasmComponent
	//wasmName, wasmSize, err := StoreFileAsCAS(wasmFilePath, blobDir)
	wasmFile, err := os.Open(wasmFilePath)	
	if err != nil {
		return "",0, err
	}

	wasmCASFileName, wasmCASFileSize, err := StoreAsCAS(wasmFile, blobDir)
	if err != nil {
		return "",0, err	
	}

	// TODO: func CreateOCIConfig
	config,err := NewConfig(wasmCASFileName)
	if err != nil {
		return "", 0, err
	}
	
	wasmCompBuffer, err := MarshalToBuffer(config)
	if err != nil {
		return "", 0, err
	}

	configCASFileName, configCASFileSize, err := StoreAsCAS(wasmCompBuffer, blobDir)	
	if err != nil {
		return "", 0, err
	}

	// TODO: func CreateOCIManifest
	configDesc := NewDescriptor(mediaTypeConfig,configCASFileName,configCASFileSize,
	)

	layerDesc := NewDescriptor(mediaTypeLayer,wasmCASFileName,wasmCASFileSize)

	manifest := NewManifest(mediaTypeManifest, configDesc, []Descriptor{layerDesc},manifestSchemaVersion)

	manifestBuffer, err := MarshalToBuffer(manifest)
	if err != nil {
		return "", 0, err
	}

	_, _, err = StoreAsCAS(manifestBuffer, blobDir)	
	if err != nil {
		return "", 0, err
	}

	//TODO: create index.json

	return "",0,nil
}
