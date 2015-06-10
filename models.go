package resource

type Source struct {
	Token      string `json:"token"`
	ProjectID  string `json:"project_id"`
	TrackerURL string `json:"tracker_url"`
	Repo       string `json:"repo"`
	PrivateKey string `json:"private_key"`
}

type Version struct {
	StoryID   int    `json:"story_id"`
	Ref       string `json:"ref"`
	Timestamp int64  `json:"timestamp"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
