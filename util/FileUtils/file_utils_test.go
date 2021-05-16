package FileUtils

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestListDirFiles(t *testing.T) {

	showFileInfos(ListDirFiles("../", nil, 1))
	showFileInfos(ListDirFiles("./", nil, 1))
	showFileInfos(ListDirFiles("../../", nil, 1))
}

func showFileInfos(fileInfos []*FileInfo) {
	fmt.Println("========================================================================================================================")
	if len(fileInfos) < 1 {
		fmt.Println("FileInfos is empty")
	} else {
		for idx, fi := range fileInfos {
			fmt.Println(idx, ":", fi.Info.Name(), ", Dir:", fi.Info.IsDir(), ", Path:", fi.Path)
		}
	}
}

func TestScanParent(t *testing.T) {

	_ = ScanParent(".//", func(parent *FileInfo) (stop bool) {

		fmt.Println(parent.Path, parent.Info.Name())
		return false
	})
}

func TestGetParentPath(t *testing.T) {
	assert.Equal(t, GetParentPath("/data/log/user"), "/data/log")
	assert.Equal(t, GetParentPath("/data/log///user"), "/data/log")
	assert.Equal(t, GetParentPath("/data/log/user/"), "/data/log")
	assert.Equal(t, GetParentPath("/data"), "")

	assert.Equal(t, GetParentPath("E:/data/log/user"), "E:/data/log")
	assert.Equal(t, GetParentPath("E:\\data\\log\\user/"), "E:\\data\\log")
}
