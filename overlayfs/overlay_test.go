package overlayfs

import "testing"

func TestOverlay(t *testing.T) {
	MountOverlay("/root/busybox", "first")
}

func TestDelMnt(t *testing.T) {
	DeleteOverlayMnt("db0ee06c")
}
