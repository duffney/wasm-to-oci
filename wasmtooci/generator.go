package wasmtooci

import (
)

//TODO: add components as input
//Should be method??
func NewConfig(hash string) (Config, error) {
	// create descriptor for config
	wasmBinDigest := "sha256:"+hash
	desc := Descriptor{
		Digest: wasmBinDigest,
	}

	//component for config 
	//TODO: implmement custom marshal interface for import[] & target null
	//TODO: extract exports and imports from .wasm
	component := Component{
		Exports: []string{"wasi:fake/todo", "wasi:fake/later"},
	}

	//create config
	config := Config{
		Author: "Josh Duffney",
		Created: "2024-06-26T17:57:13.394956Z",
		//Created: time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"),
		Architecture: "wasm",
		Os: "wasip2",
		LayerDigests: []Descriptor{desc},
		Component: component,
	}
	return config, nil
}


/*
NewDescriptor creates a new Descriptor with mandatory fields and optional parameters.
desc := NewDescriptor("application/wasm", "sha256:abcdef123456", 1024,
    WithUrls([]string{"https://example.com/artifact"}),
    WithAnnotations([]map[string]string{{"key1": "value1"}, {"key2": "value2"}}),
    WithData("additional data"),
    WithArtifactType("wasm module"),
)
*/
func NewDescriptor(mediaType, digest string, size int64, options ...func(*Descriptor)) Descriptor {
	desc := Descriptor{
		MediaType: mediaType,
		Digest:    digest,
		Size:      size,
	}

	for _, option := range options {
		option(&desc)
	}

	return desc
}

// WithUrls sets the Urls field of Descriptor.
func WithUrls(urls []string) func(*Descriptor) {
	return func(d *Descriptor) {
		d.Urls = urls
	}
}

// WithAnnotations sets the Annotations field of Descriptor.
func WithAnnotations(annotations []map[string]string) func(*Descriptor) {
	return func(d *Descriptor) {
		d.Annotations = annotations
	}
}

// WithData sets the Data field of Descriptor.
func WithData(data string) func(*Descriptor) {
	return func(d *Descriptor) {
		d.Data = data
	}
}

// WithArtifactType sets the ArtifactType field of Descriptor.
func WithArtifactType(artifactType string) func(*Descriptor) {
	return func(d *Descriptor) {
		d.ArtifactType = artifactType
	}
}

func NewManifest(mediaType string, config Descriptor, layers []Descriptor, schemaVersion int) (Manifest)  {
	return Manifest{
		Config: config,
		Layers: layers,
		MediaType: mediaType,
		SchemaVersion: schemaVersion,
		}
}
