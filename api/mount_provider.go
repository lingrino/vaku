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
		return nil, newWrapErr("", ErrMountInfo, newWrapErr(err.Error(), ErrListMounts, nil))
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

// StaticMountProvider is a mount provider that returns a single static mount.
// This is useful when the user doesn't have permission to list mounts but knows
// the mount path and version.
type StaticMountProvider struct {
	mount Mount
}

// NewStaticMountProvider creates a new static mount provider with the given mount.
func NewStaticMountProvider(path, version string) *StaticMountProvider {
	return &StaticMountProvider{
		mount: Mount{
			Path:    path,
			Type:    "kv",
			Version: version,
		},
	}
}

// ListMounts returns the single static mount.
func (p *StaticMountProvider) ListMounts() ([]Mount, error) {
	return []Mount{p.mount}, nil
}
