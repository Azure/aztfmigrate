package types

type Mapping struct {
	ResourceType         string `json:"resourceType"`
	ExampleConfiguration string `json:"exampleConfiguration,omitempty"`
	IdPattern            string `json:"idPattern"`
}
