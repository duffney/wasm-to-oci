package wasmtooci

type Converter struct {
}

type Config struct {
	Author string `json:"author,omitempty"`
	Created string `json:"time"`
	Architecture string `json:"architecture"`
	Os string `json:"os"`
	LayerDigests []Descriptor `json:"layerDigest,omitempty"`
	Component Component `json:"component,omitempty"`
}

type Component struct {
	Exports []string `json:"exports"`
	Imports []string `json:"imports"`
	Target string `json:"target"`
}

type Descriptor struct {
	Digest string `json:"digest"`
	MediaType string `json:"mediaType,omitempty"`
	Size int64 `json:"size,omitempty"`
	Urls []string `json:"urls,omitempty"`
	Annotations []map[string]string `json:"annotations,omitempty"`
	Data string `json:"data,omitempty"`
	ArtifactType string `json:"artifactType,omitempty"`
}

type Manifest struct {
	Config Descriptor `json:"config"` 
	Layers []Descriptor `json:"layers"`
	MediaType string `json:"mediaType"`
	SchemaVersion int `json:"schemaVersion"`
}

type Index struct {
	SchemaVersion int `json:"schemaVersion"`
	Manifests []Descriptor `json:"manifests"`
	Annotations []map[string]string `json:"annotations,omitempty"`
}
