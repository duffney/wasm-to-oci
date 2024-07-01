package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/duffney/wasm-to-oci/wasmtooci"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

//TODO: add const for media types
const (
	ociDir string = "blobs/sha256"
	mediaTypeIndex = "application/vnd.oci.image.manifest.v1+json"
)

//TODO **define manifest type
type Manifest struct {
	Config ManifestConfig 
	Layers []ManifestLayer
	MediaType string `json:"mediaType"`
	SchemaVersion int `json:"schemaVersion"`
}

//TODO refactor redudant with ManifestLayer
type ManifestConfig struct {
	Digest string `json:"digest"`
	MediaType string `json:"mediaType"`
	Size int64 `json:"size"`
}

type ManifestLayer struct {
	Digest string `json:"digest"`
	MediaType string `json:"mediaType"`
	Size int64 `json:"size"`
}

type Descriptor struct {
	Digest string `json:"digest"`
	MediaType string `json:"mediaType"`
	Size int64 `json:"size"`
}

func (d *Descriptor) New(p string) (Descriptor, error){
	// Open the file
	file, err := os.Open(p)
	if err != nil {
		return Descriptor{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return Descriptor{}, fmt.Errorf("failed to get file info: %w", err)
	}
	size := fileInfo.Size()

	// Compute SHA256 hash of the file content
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return Descriptor{}, fmt.Errorf("failed to hash file content: %w", err)
	}
	digest := hex.EncodeToString(hash.Sum(nil))

	return Descriptor{
		Digest:    digest,
		MediaType: mediaTypeIndex,
		Size:      size,
	}, nil
}

func NewDescriptor(p string) (Descriptor, error) {

	// Open the file
	file, err := os.Open(p)
	if err != nil {
		return Descriptor{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return Descriptor{}, fmt.Errorf("failed to get file info: %w", err)
	}
	size := fileInfo.Size()

	// Compute SHA256 hash of the file content
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return Descriptor{}, fmt.Errorf("failed to hash file content: %w", err)
	}
	digest := hex.EncodeToString(hash.Sum(nil))

	return Descriptor{
		Digest:    digest,
		MediaType: mediaTypeIndex,
		Size:      size,
	}, nil
}

//TODO: rename to ociconfig
type ImageConfig struct {
	Author string `json:"author,omitempty"`
	Created time.Time `json:"time"`
	Architecture string `json:"architecture"`
	Os string `json:"os"`
	LayerDigests []string `json:"layerDigest,omitempty"`
	Component Component `json:"component,omitempty"`
}

type Component struct {
	Exports []string `json:"exports"`
	Imports []string `json:"imports"`
	Target string `json:"target"`
}

type Index struct {
	SchemaVersion int
	Manifests []Descriptor
}

func main() {
	// createWasmLayer > genConfig > genImgManifest > genImageIndex
	// TODO: How do I create a layer from a wasm component?	
	// take a wasm component file, calc sha256 hash, rename file to hash, mv file to blobs/sha256, return hash
	//wasmFile := "fake.wasm"
	//comp,compSize, err := createWasmLayer(wasmFile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(comp)
	//fmt.Printf("Component:%s Size:%d", comp, compSize)

	//config,configSize := createOciWasmConfig([]string{comp})
	//fmt.Println(config)
	//fmt.Println(configSize)

	//manConfig := createManifestConfig(config,configSize) 
	//manLayer := createManifestLayer(comp,compSize)

	//_ = createManifest(manConfig, manLayer)	

	//TODO create Image Index 
	//indexDesc, err := NewDescriptor("blob/sha256/"+manifestName)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//createIndex(indexDesc)

	/*
		v2
		1. define path to wasmfile and dest
		2. init convert type
		3. convert wasm to oci
	*/
	c, err := wasmtooci.NewConverter()
	if err != nil {
		log.Fatal(err)
	}

	wasmFilePath := "fake.wasm"
	name, size, err := c.Convert(wasmFilePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Component hash:%s size: %d\n", name, size)
}

//TODO: add filepath as arg 
func createWasmLayer(p string) (string,int64,error) {
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	//TODO: repalce with const
	layerDir := "blobs/sha256"

	if _, err := os.Stat(layerDir); os.IsNotExist(err) {
		err := os.MkdirAll(layerDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	//TODO: move file to blob with hash name
    hashFile, err := os.Create(filepath.Join(layerDir,hashString))
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = io.Copy(hashFile, f)	
	if err != nil {
		log.Fatal(err)
	}

	binFile, err := os.Stat(filepath.Join(layerDir,hashString))
	if err != nil {
		log.Fatal(err)
	}


	return hashString,binFile.Size(), nil
}

//TODO: Add MarshalJSON implements the custom JSON marshaling for Component
// handles target nul when empty
func createOciWasmConfig(digests []string) (string,int64) {
    component := Component{
        Exports: []string{"wasi:fake/todo", "wasi:fake/later"},
    }
	
	var shaDigests []string
	for _, d := range digests {
		shaDigests = append(shaDigests,"sha256:"+d)
	}

    config := ImageConfig{
        Author:       "Josh Duffney",
        Created:      time.Now().Round(0), // Round to seconds for consistency
        Architecture: "wasm",
        Os:           "wasip2",
        LayerDigests: shaDigests,
        Component:    component,
    }

    // Marshal the config struct to JSON with consistent formatting
    byteData, err := json.Marshal(config)
    if err != nil {
        log.Fatal(err)
    }

    // Calculate SHA-256 hash
    hash := sha256.Sum256(byteData)
    configName := hex.EncodeToString(hash[:])

    // Debugging: Print the hash and JSON for verification
    fmt.Println("Computed hash:", configName)
    fmt.Println("JSON data:", string(byteData))

    // Write JSON data to a file based on the hash
    configPath := filepath.Join(ociDir, configName)
    err = os.WriteFile(configPath, byteData, 0644)
    if err != nil {
        log.Fatal(err)
    }

	configFile, err := os.Stat(configPath)
	if err != nil {
		log.Fatal(err)
	}

    return configName, configFile.Size() 
}

// TODO: convert to method
func createManifestConfig (d string, s int64) ManifestConfig {
	digest := "sha256:"+d
	config := ManifestConfig{
		Digest: digest,
		MediaType: "application/vnd.wasm.config.v0+json",
		Size: s,
	}

	return config 
}

//TODO: convert to method
func createManifestLayer (d string, s int64) ManifestLayer {
	digest := "sha256:"+d

	layer := ManifestLayer{
		Digest: digest,
		MediaType: "application/wasm",
		Size: s,
	}

	return layer
}

func createManifest (c ManifestConfig, l ManifestLayer) string {
	manifest := Manifest{
		Config: c,
		Layers: []ManifestLayer{l},
		MediaType: "application/vnd.oci.image.manifest.v1+json",
		SchemaVersion: 2,
	}

    // Marshal the config struct to JSON with consistent formatting
    byteData, err := json.Marshal(manifest)
    if err != nil {
        log.Fatal(err)
    }

    // Calculate SHA-256 hash
    hash := sha256.Sum256(byteData)
    manifestName := hex.EncodeToString(hash[:])

    // Debugging: Print the hash and JSON for verification
    fmt.Println("Computed hash:", manifest)
    fmt.Println("JSON data:", string(byteData))

    // Write JSON data to a file based on the hash
    manifestPath := filepath.Join(ociDir, manifestName)
    err = os.WriteFile(manifestPath, byteData, 0644)
    if err != nil {
        log.Fatal(err)
    }

	return manifestName
}


//TODO: old 
func genImageIndex() {
    index := v1.Index{
        MediaType: v1.MediaTypeImageIndex,
        Manifests: []v1.Descriptor{
            {
                MediaType:   v1.MediaTypeImageManifest,
                Digest:      "sha256:exampledigest1",
                Size:        7023,
                Annotations: map[string]string{"org.opencontainers.image.ref.name": "example1"},
            },
            {
                MediaType:   v1.MediaTypeImageManifest,
                Digest:      "sha256:exampledigest2",
                Size:        7033,
                Annotations: map[string]string{"org.opencontainers.image.ref.name": "example2"},
            },
        },
        Annotations: map[string]string{
            "example.key": "example.value",
        },
    }

    // Convert the index struct to JSON
    jsonData, err := json.MarshalIndent(index, "", "    ")
    if err != nil {
        fmt.Println("Error marshalling index:", err)
        return
    }

    // Write the JSON data to a file
    file, err := os.Create("index.json")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    _, err = file.Write(jsonData)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    fmt.Println("index.json created successfully")
}
/*
flow genLayers, genManifest, genConfig, genImageIndex
*/
//TODO: funcs genIndex, genManifest, getConfig, genLayers
