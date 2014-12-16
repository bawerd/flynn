package volume

import (
	"github.com/flynn/flynn/host/types"
)

/*
	A Volume is a persistent and sharable filesystem.  Unlike most of the filesystem in a job's
	container, which is ephemeral and is discarded after job termination, Volumes can be used to
	store data and may be reconnected to a later job (or to multiple jobs).

	Volumes may also support additional features for their section of the filesystem, such
	storage quotas, read-only mounts, snapshotting operation, etc.

	The Flynn host service maintains a locally persistent knowledge
	of mounts, and supplies this passively to the orchestration API.
	The host service does *not* perform services such as garbage collection of unmounted
	volumes (how is it to know if you still want that data preserved for a future job?)
	or transport and persistence of volumes between hosts (that should be orchestrated via
	the API from a higher level service).
*/
type Volume interface {
	ID() string // guid (v4, random.  not globally sync'd, entropy should be high enough to be unique)
	Mounts() map[VolumeMount]struct{}
	// Note: NOT provided: a method that gets the host's path to a mount.  Not all backends have such a useable raw path on the host.
	// REVIEW: but maybe they could?  Designing such that we always have such a staging area outside of the job container could be done...?

	Mount(job host.ActiveJob, path string) (VolumeMount, error)

	TakeSnapshot() (Volume, error)
}

/*
	VolumeMount names the location in which a shared+persistent filesystem is mounted into a job's container.

	A Volume has a one-to-many relationship with `VolumeMount`s -- the same volume
	may be mounted to many containers (or even multiple places within a single container).
*/
type VolumeMount struct {
	JobID    string // job which the volume is mounted to // REVIEW: may be appropriate to use the reified ActiveJob since it doesn't make sense to consider a mount on a container that's dead... but it's also nice to have this object make sense as a map key
	Location string // path within the container where the mount shall appear
}

func NewVolume() (Volume, error) {
	// TODO: probably needs to sprout some sort of factory or attach to existing host state, since the boltdb will be a near-singleton
	return nil, nil
}

func NewVolumeFromSnapshot( /* host reference?? */ v Volume) (Volume, error) {
	// TODO
	return nil, nil
}

func NewVolumeUsingBackend(b Backend) (Volume, error) {
	// TODO
	// not actually sure this is going to end up being at all well defined.  construction might need different inputs per backend.
	return nil, nil
}

type Backend string

const (
	BindmountBackend = "bind"
	ZFSBackend       = "zfs"
)
