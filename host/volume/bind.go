package volume

import (
	"fmt"

	"github.com/flynn/flynn/host/types"
	"github.com/flynn/flynn/pkg/random"
)

type bindmountVolume struct {
	id         string
	storageDir string
	mounts     map[VolumeMount]struct{}
}

func NewBindVolume(hostPath string) Volume {
	return &bindmountVolume{
		id:         random.UUID(),
		storageDir: hostPath,
		mounts:     make(map[VolumeMount]struct{}),
	}
	// not much for potential error codes; nothing's mounted yet, and for direct bind mounts, we don't allocate our own internal storage.
}

func (v *bindmountVolume) ID() string {
	return v.id
}

func (v *bindmountVolume) Mounts() map[VolumeMount]struct{} {
	return v.mounts
}

func (v *bindmountVolume) Mount(job host.ActiveJob, path string) (VolumeMount, error) {
	mount := VolumeMount{
		JobID:    job.Job.ID,
		Location: path,
	}
	if _, exists := v.mounts[mount]; exists {
		return VolumeMount{}, fmt.Errorf("volume: cannot make same mount twice!")
	}
	// TODO: fire syscalls
	v.mounts[mount] = struct{}{}
	return mount, nil
}

func (v *bindmountVolume) TakeSnapshot() (Volume, error) {
	return nil, fmt.Errorf("snapshots not supported on bind mount volumes")
}
