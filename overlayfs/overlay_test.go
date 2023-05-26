package overlayfs

import "testing"

func TestOverlay(t *testing.T) {
	MountOverlay("/root/busybox", "first")
}
