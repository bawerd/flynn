package zfs

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"

	zfs "github.com/flynn/flynn/Godeps/_workspace/src/github.com/mistifyio/go-zfs"

	. "github.com/flynn/flynn/Godeps/_workspace/src/github.com/flynn/go-check"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

/*
option 1
--------

apt-get install zfs-fuse


option 2
--------

gpg --keyserver pgp.mit.edu --recv-keys "F6B0FC61"
gpg --armor --export "F6B0FC61" | apt-key add -
echo deb http://ppa.launchpad.net/zfs-native/stable/ubuntu trusty main > /etc/apt/sources.list.d/zfs.list
apt-get update
apt-get install -y ubuntu-zfs

this drags in g++, and proceeds to spend a large number of seconds on "Building initial module for 3.13.0-43-generic"...
okay, minutes.  ugh.
*/

//

func (S) TestThings(c *C) {
	dataset, err := prepareYourself("/mount/path", "test")

	fmt.Printf("::: %#v\n", dataset)
	fmt.Printf("::: %#v\n", err)
}

func prepareYourself(mountPath string, poolName string) (*zfs.Dataset, error) {
	if _, err := exec.LookPath("zfs"); err != nil {
		return nil, fmt.Errorf("zfs command is not available")
	}

	backingFile, err := ioutil.TempFile("/tmp/", "zfs-")
	if err != nil {
		return nil, err
	}
	defer backingFile.Close()

	err = backingFile.Truncate(int64(math.Pow(2, float64(30))))
	if err != nil {
		return nil, err
	}
	defer os.Remove(backingFile.Name())
	pool, err := zfs.CreateZpool(poolName, nil, backingFile.Name()) // the default point where this mounts is in "/poolName", so... you're gonna wanna override that
	if err != nil {
		return nil, err
	}
	defer pool.Destroy() // this appears to somehow race itself or otherwise fail, sometimes.  makes cleanup a titch hard, gonna need to get to the bottom of that.  may be because we're holding references to things after the destroy, but it should at least error.

	logger := Logger{}
	zfs.SetLogger(&logger) // package wide global.  ugh.  have to propagate their singleton mistakes out and do this once somewhere.

	dataset, err := zfs.GetDataset(poolName)
	if err != nil {
		return nil, err
	}

	_, err = zfs.CreateFilesystem(path.Join(poolName, "beta"), map[string]string{
		"mountpoint": mountPath,
	})
	if err != nil {
		return nil, err
	}

	mountedDir, err := os.Open(mountPath)
	if err != nil {
		return nil, err
	}
	dirlist, err := mountedDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("::: mountpath dirlist, fresh: %#v\n", dirlist)

	_, err = os.Create(filepath.Join(mountPath, "alpha"))

	mountedDir, err = os.Open(mountPath)
	if err != nil {
		return nil, err
	}
	dirlist, err = mountedDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("::: mountpath dirlist, written: %#v\n", dirlist)

	mountedDir, err = os.Open("/" + poolName)
	if err != nil {
		return nil, err
	}
	dirlist, err = mountedDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("::: poolmount dirlist: %#v\n", dirlist)

	err = cloneFilesystem(poolName+"2", poolName, mountPath+"2")
	if err != nil {
		return nil, err
	}

	mountedDir, err = os.Open(mountPath + "2")
	if err != nil {
		return nil, err
	}
	dirlist, err = mountedDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	fmt.Printf("::: mountpath2 dirlist, forked: %#v\n", dirlist)

	return dataset, nil
}

func cloneFilesystem(newDatasetName string, parentDatasetName string, mountPath string) error {
	parentDataset, err := zfs.GetDataset(parentDatasetName)
	if parentDataset == nil {
		return err
	}
	snapshotName := fmt.Sprintf("%d", time.Now().Nanosecond())
	snapshot, err := parentDataset.Snapshot(snapshotName, false)
	if err != nil {
		return err
	}

	_, err = snapshot.Clone(newDatasetName, map[string]string{
		"mountpoint": mountPath,
	})
	if err != nil {
		snapshot.Destroy(zfs.DestroyDeferDeletion)
		return err
	}
	err = snapshot.Destroy(zfs.DestroyDeferDeletion)
	return err
}
