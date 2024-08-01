package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"logger"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	Mutex      sync.Mutex
	sum256Map  map[string]string
	dirty      bool
	sum256List []string
)

func UtilsInit() {
	os.MkdirAll(viper.GetString("originalPrefix"), 0777)
	os.MkdirAll(viper.GetString("thumbnailPrefix"), 0777)

	RefreshSum256Map()
	RefreshSum256List()
}

func GetTotalPage(pageSize int) int {
	lenSum256List := len(sum256List)
	if lenSum256List%pageSize == 0 {
		return lenSum256List / pageSize
	}
	return len(sum256List)/pageSize + 1
}

func RefreshSum256List() {
	Mutex.Lock()
	defer Mutex.Unlock()

	sum256List = make([]string, 0)
	for k := range sum256Map {
		sum256List = append(sum256List, k+".jpg")
	}
	sort.Strings(sum256List)
}

func GetSum256ListPagingQuery(pageNum, pageSize int) []string {
	if dirty {
		RefreshSum256List()
		dirty = false
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	startIndex := (pageNum - 1) * pageSize
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := startIndex + pageSize
	if endIndex > len(sum256List) {
		endIndex = len(sum256List)
	}
	return sum256List[startIndex:endIndex]
}

func GetSHA256(imageBytes []byte) string {
	hasher := sha256.New()
	if _, err := hasher.Write(imageBytes); err != nil {
		logger.Logger.Println(err)
		return ""
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func IsDirty() bool {
	return dirty
}

func CheckHexStr(checkStr string) bool {
	if checkStr == "" {
		return false
	}
	if len(checkStr)&1 != 0 {
		checkStr = "0" + checkStr
	}

	_, err := hex.DecodeString(checkStr)
	return err == nil
}

func AddSum256Map(newSum256 string) bool {
	// 检查格式
	if ok := CheckHexStr(newSum256); !ok {
		return false
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	sum256Map[newSum256] = newSum256 + ".jpg"
	return true
}

func CheckSum256InSum256Map(checkSum256 string) bool {
	_, ok := sum256Map[checkSum256]
	return ok
}

func RemoveSum256Map(rmSum256 string) {
	if ok := CheckHexStr(rmSum256); !ok {
		return
	}
	if ok := CheckSum256InSum256Map(rmSum256); !ok {
		return
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	delete(sum256Map, rmSum256)
	dirty = true
}

func RefreshSum256Map() {
	tmpSum256Map := make(map[string]string, 0)

	dir, err := os.Open(viper.GetString("thumbnailPrefix"))
	if err != nil {
		logger.Logger.Fatal(err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		logger.Logger.Fatal(err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			tmpSum256Map[strings.Split(name, ".")[0]] = name
		}
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	sum256Map = tmpSum256Map
	dirty = true
}
