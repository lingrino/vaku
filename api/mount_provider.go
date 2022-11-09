package vaku

type Mount struct {
	Path    string
	Type    string
	Version string
}

type mountProvider interface {
	ListMounts() ([]Mount, error)
}

type defaultMountProvider struct {
	client *Client
}

func (p defaultMountProvider) ListMounts() ([]Mount, error) {
	mounts, err := p.client.vc.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	result := make([]Mount, 0)
	for mountPath, data := range mounts {
		// Ensure '/' so that no match on foo/bar/ when actual path is foo/barbar/
		mount := Mount{
			Path:    mountPath,
			Type:    data.Type,
			Version: data.Options["version"],
		}
		result = append(result, mount)
	}
	return result, nil
}
