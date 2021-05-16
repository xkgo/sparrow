package FileUtils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
判断给定的 文件路径是否是一个存在的目录
*/
func IsDirExists(path string) bool {
	if len(path) < 1 {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	if !fileInfo.IsDir() {
		return false
	}
	return true
}

type FileInfo struct {
	Info   os.FileInfo
	Path   string    // 文件全路径
	Parent *FileInfo // 上级目录
}

/**
获取没指定文件夹下的所有文件列表
@param dir 要操作的目录
@param acceptor 过滤器，该接口符合条件的情况下才会返回，如果为 nil 则返回全部
@param depth 扫描深度，有效值>0, 如果传入一个无效值，那么默认就是1，即返回该目录第一层子文件列表， 否则的话，返回第 depth 深度下面的文件
*/
func ListDirFiles(dir string, acceptor func(fileInfo os.FileInfo) bool, depth int) []*FileInfo {
	return listDirFilesWithDepth(dir, acceptor, 1, depth, nil)
}

func listDirFilesWithDepth(dir string, acceptor func(fileInfo os.FileInfo) bool, curDepth, maxDepth int, files []*FileInfo) []*FileInfo {
	if curDepth < 1 {
		curDepth = 1
	}
	if maxDepth < 1 {
		maxDepth = 1
	}

	if curDepth > maxDepth {
		return files
	}

	if !IsDirExists(dir) {
		return make([]*FileInfo, 0)
	}
	dir, _ = filepath.Abs(dir)
	if files == nil {
		files = make([]*FileInfo, 0)
	}
	subFiles, err := ioutil.ReadDir(dir)
	if err != nil || len(subFiles) < 1 {
		return files
	}

	dirInfo, _ := os.Stat(dir)
	parent := &FileInfo{
		Info:   dirInfo,
		Path:   dir,
		Parent: nil,
	}
	acceptFiles := make([]*FileInfo, 0)
	for _, fi := range subFiles {
		fii := &FileInfo{
			Info:   fi,
			Path:   dir + string(filepath.Separator) + fi.Name(),
			Parent: parent,
		}
		if acceptor == nil || acceptor(fi) {
			acceptFiles = append(acceptFiles, fii)
		}
		if fi.IsDir() {
			// 如果是目录，继续遍历
			acceptFiles = listDirFilesWithDepth(dir+string(filepath.Separator)+fi.Name(), acceptor, curDepth+1, maxDepth, acceptFiles)
		}
	}

	if len(acceptFiles) > 0 {
		files = append(files, acceptFiles...)
	}

	return files
}

/**
扫描父级文件，直到没有父级为止
@param path 源文件
@param consumer 当扫描到父级目录的时候，执行consumer处理，返回true则不会继续往上扫描
@return error 如果扫描异常则返回 error
*/
func ScanParent(path string, consumer func(parent *FileInfo) (stop bool)) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	parentPath := path[0 : len(path)-len(fi.Name())-1]
	for parentPath != path {
		pfi, err := os.Stat(parentPath)
		if err != nil {
			return err
		}
		fileInfo := &FileInfo{
			Info: pfi,
			Path: parentPath,
		}
		if consumer(fileInfo) {
			return nil
		}
		path = parentPath
		parentPath = parentPath[0 : len(parentPath)-len(pfi.Name())-1]
	}
	return nil
}

/**
获取给定路径的上级目录路径，返回的上级目录中，不包含 / 结尾
实例：
/data/web/log 		--> /data/web
E:\data\web\log 	--> E:\data\web
/					--> ""
/data			    --> ""
*/
func GetParentPath(path string) string {
	if len(path) < 1 {
		return ""
	}
	var sep1 = '\\'
	var sep2 = '/'

	strRunes := []rune(path)
	lastIndex := len(strRunes) - 1

	firstNotSepIndex := lastIndex
	// 获取第一个不是分隔符的下标
	for i := lastIndex; i >= 0; i-- {
		r := strRunes[i]
		if r != sep1 && r != sep2 {
			firstNotSepIndex = i
			break
		}
	}

	// 从第一个不是分隔符下标开始对比
	foundSep := false
	for i := firstNotSepIndex; i >= 0; i-- {
		r := strRunes[i]
		if r != sep1 && r != sep2 {
			if foundSep {
				return string(strRunes[0 : i+1])
			}
		} else {
			foundSep = true
		}
	}
	return ""
}
