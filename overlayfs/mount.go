// 用于挂载普通文件
package overlayfs

import (
	"os"
	"path"
	"syscall"
)

func MountFS(vol_map map[string]string, container_id string) error {
	for vol1, vol2 := range vol_map {
		if _, err := os.Stat(vol1); os.IsNotExist(err) {
			//不存在则创建
			if err = os.MkdirAll(vol1, 0755); err != nil {
				return err
			}
		}
		//vol2转换成容器内部的路径
		vol2 = path.Join(MountPoint, container_id, "merge", vol2)
		if _, err := os.Stat(vol2); os.IsNotExist(err) {
			//不存在则创建
			if err = os.MkdirAll(vol2, 0755); err != nil {
				return err
			}
		}
		//挂载 vol1 ——> vol2
		err := syscall.Mount(vol1, vol2, "", syscall.MS_BIND, "")
		if err != nil {
			return err
		}
	}
	return nil
}

// 卸载mount
func RemoveMountFs(vol_map map[string]string, container_id string) error {
	for _, vol2 := range vol_map {
		//vol2转换成容器内部的路径
		vol2 = path.Join(MountPoint, container_id, "merge", vol2)
		//卸载mount
		if err := syscall.Unmount(vol2, 0); err != nil {
			return err
		}
	}
	return nil
}
