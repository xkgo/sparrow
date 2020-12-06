package env

import (
	"github.com/xkgo/sparrow/util/FileUtils"
	"github.com/xkgo/sparrow/util/StringUtils"
	"os"
	"regexp"
)

const (
	// sparrow profile dirs key
	SparrowProfileDirsKey = "sparrow-profile-dirs"
)

/*
	activeProfiles 配置搜索文件夹列表，来源规则如下：
	1. 如果是程序自己设置了（相当于是程序自己定义了路径），那么直接使用程序自己设定的
	2. 上一步没有获取到，解析命令行参数，当存在 --sparrow-profile-dirs=...... 的时候，那么直接以 --sparrow-profile-dirs 指定的为准，
	   如果--sparrow-profile-dirs设置了，但是为空字符串，那么默认: ./,./config,./conf
	3. 上一步没有获取到，那么计算系统环境变量中，是否定义了 sparrow-profile-dirs， 如果定义了就以此为准，如果定义的是空字符串，那么就是默认：./,./config,./conf
	4. 上一步没有，那么检查：./,./config,./conf，搜索是否存在（理论上是本框架支持的文件格式） application*.properties|yml|toml
	5. 上一步没有，向上获取到一个目录，然后重复上一步，直到找到符合上一步的为止
	6. 如果始终找不到，那么说明不需要配置文件，系统一样是可以运行的
*/
func resolveProfileDirs(customDirs []string) []string {
	if customDirs != nil && len(customDirs) > 0 {
		// 检查每一个文件夹，是否包含 application*.properties|yml|toml
		return filterAndGetValidProfileDirs(customDirs)
	}
	// 解析命令行参数，当存在 --sparrow-profile-dirs=...... 的时候，那么直接以 --sparrow-profile-dirs 指定的为准
	tempDirs, exists := GetCommandLineProperty(SparrowProfileDirsKey)
	if exists {
		if len(tempDirs) < 1 {
			tempDirs = "./,./config,./conf"
		}
		return filterAndGetValidProfileDirs(StringUtils.SplitByRegex(tempDirs, "[,，；;]+"))
	}

	// 计算系统环境变量中，是否定义了 sparrow-profile-dirs， 如果定义了就以此为准，如果定义的是空字符串，那么就是默认：./,./config,./conf
	tempDirs, exists = os.LookupEnv(SparrowProfileDirsKey)
	if exists {
		if len(tempDirs) < 1 {
			tempDirs = "./,./config,./conf"
		}
		return filterAndGetValidProfileDirs(StringUtils.SplitByRegex(tempDirs, "[,，；;]+"))
	}

	// 使用默认的
	profileDirs := filterAndGetValidProfileDirs([]string{"./", "./config", "./conf"})
	if len(profileDirs) > 0 {
		return profileDirs
	}

	profileDir := ""
	// 循环向上一目录进行查找，直到找到有配置文件的为止
	_ = FileUtils.ScanParent("./", func(parent *FileUtils.FileInfo) (stop bool) {
		path := parent.Path
		if isValidProfileDir(path) {
			profileDir = path
			return true
		}
		return false
	})

	if len(profileDir) > 0 {
		profileDirs = make([]string, 0)
		return append(profileDirs, profileDir)
	}

	return make([]string, 0)
}

func filterAndGetValidProfileDirs(profileDirs []string) []string {
	validProfileDirs := make([]string, 0)
	if profileDirs != nil && len(profileDirs) < 1 {
		return validProfileDirs
	}

	for _, profileDir := range profileDirs {
		if isValidProfileDir(profileDir) {
			validProfileDirs = append(validProfileDirs, profileDir)
		}
	}

	return validProfileDirs
}

var applicationFileRegex, _ = regexp.Compile("(?i)^application.*\\.(properties|yml|toml)$")

/**
检查是否是合法的 profileDir 目录，合法的定义：
1. 文件夹存在
2. 该文件夹下面，包含 application*.properties|toml|yml
*/
func isValidProfileDir(profileDir string) bool {
	// 检查是否包含  application*.properties|toml|yml
	subFiles := ListDirApplicationFiles(profileDir)

	return len(subFiles) > 0
}

func ListDirApplicationFiles(dir string) []*FileUtils.FileInfo {
	return FileUtils.ListDirFiles(dir, func(fileInfo os.FileInfo) bool {
		if fileInfo.IsDir() {
			return false
		}
		return applicationFileRegex.MatchString(fileInfo.Name())
	}, 1)
}
