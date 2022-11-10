package vaku

// Mount is a high level representation of selected fields of a
// vault mount that are relevant to vaku.
type Mount struct {
	Path    string
	Type    string
	Version string
}

// mountProvider is used to get a list of all mounts that the user has access to.
type mountProvider interface {
	ListMounts() ([]Mount, error)
}

// defaultMountProvider is used if no other mountProvider is supplied.
type defaultMountProvider struct {
	client *Client
}

// ListMounts lists mounts using the sys/mounts endpoint.
func (p defaultMountProvider) ListMounts() ([]Mount, error) {
	mounts, err := p.client.vc.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	result := make([]Mount, 0)
	for mountPath, data := range mounts {
		mount := Mount{
			Path:    mountPath,
			Type:    data.Type,
			Version: data.Options["version"],
		}
		result = append(result, mount)
	}
	return result, nil
}
