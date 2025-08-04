package rclone

import (
	"context"
	"fmt"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	_ "github.com/rclone/rclone/backend/all"
	"github.com/rclone/rclone/cmd"
	"github.com/rclone/rclone/fs/operations"
)

func Sync(src string, dst string, dstFileName string) error {
	fsrc, srcFileName := cmd.NewFsFile(src)
	fdst := cmd.NewFsDir([]string{dst})
	log.Debug(fmt.Sprintf("Syncing %s %s %s %s", fsrc.Root(), fdst.Root(), srcFileName, dstFileName))
	return operations.CopyFile(context.TODO(), fdst, fsrc, dstFileName, srcFileName)
}
