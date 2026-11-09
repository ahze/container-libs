package zfs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	graphdriver "go.podman.io/storage/drivers"
	"golang.org/x/sys/unix"
)

func checkRootdirFs(rootdir string) error {
	var buf unix.Statfs_t
	if err := unix.Statfs(rootdir, &buf); err != nil {
		return fmt.Errorf("failed to access '%s': %s", rootdir, err)
	}

	// on FreeBSD buf.Fstypename contains ['z', 'f', 's', 0 ... ]
	if (buf.Fstypename[0] != 122) || (buf.Fstypename[1] != 102) || (buf.Fstypename[2] != 115) || (buf.Fstypename[3] != 0) {
		logrus.WithField("storage-driver", "zfs").Debugf("no zfs dataset found for rootdir '%s'", rootdir)
		return fmt.Errorf("no zfs dataset found for rootdir '%s': %w", rootdir, graphdriver.ErrPrerequisites)
	}

	return nil
}

func getMountpoint(id string) string {
	return id
}

func detachUnmount(mountpoint string) error {
	// FreeBSD has no lazy unmount. A plain unmount leaves a busy mount in
	// place for its other users; MNT_FORCE would revoke their open files.
	return unix.Unmount(mountpoint, 0)
}
