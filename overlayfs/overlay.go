package overlayfs

import (
	"fmt"
	"os"
	"os/exec"
)

const MountPoint = "/mnt/tiny-docker/"
const MountName = "tiny-docker"

// /mnt/tiny-docker/<容器名>/{merge,lower,upper,lower}
func CreateOverlayFsDirs(containerName string) error {
	// err := os.MkdirAll(MountPoint+containerName+"/lower", 0755)
	// if err != nil {
	// 	return err
	// }
	err := os.MkdirAll(MountPoint+containerName+"/upper", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(MountPoint+containerName+"/work", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(MountPoint+containerName+"/merge", 0755)
	if err != nil {
		return err
	}
	return nil
}

func MountOverlay(lowerPath string, containerName string) error {
	CreateOverlayFsDirs(containerName)
	upperPath := MountPoint + containerName + "/upper"
	workPath := MountPoint + containerName + "/work"
	mergePath := MountPoint + containerName + "/merge"
	mountinfo := "sudo mount -t overlay " + MountName + " -o " + "lowerdir=" + lowerPath + "," + "upperdir=" + upperPath + "," + "workdir=" + workPath + " " + mergePath
	cmd := exec.Command("/bin/bash", "-c", mountinfo)
	if err := cmd.Run(); err != nil {
		fmt.Println("mount err = ", err)
		return err
	}
	return nil
}
