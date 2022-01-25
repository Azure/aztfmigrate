package coverage

type Coverage struct {
	ApiVersion string   `json:"api_version"`
	IdPattern  string   `json:"api_path"`
	Operation  string   `json:"operation"`
	Properties []string `json:"properties"`
}
