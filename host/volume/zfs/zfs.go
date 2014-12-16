package zfs

import (
	"fmt"

	"github.com/flynn/flynn/host/types"
	. "github.com/flynn/flynn/host/volume"
	"github.com/flynn/flynn/pkg/random"
)

type zfsVolume struct {
	id     string
	mounts map[VolumeMount]struct{}
}

func NewZFSVolume() (Volume, error) {
	v := &zfsVolume{
		id:     random.UUID(),
		mounts: make(map[VolumeMount]struct{}),
	}
	// TODO: claim area, make new filesystem
	return v, nil
}

func (v *zfsVolume) ID() string {
	return v.id
}

func (v *zfsVolume) Mounts() map[VolumeMount]struct{} {
	return v.mounts
}

func (v *zfsVolume) Mount(job host.ActiveJob, path string) (VolumeMount, error) {
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

func (v *zfsVolume) TakeSnapshot() (Volume, error) {
	// TODO: lots
	return nil, nil
}
