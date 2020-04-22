package file

import (
	"bufio"
	"io/ioutil"
	"os"
	Path "path/filepath"
	"strings"
	"sync"
)

func IsExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsExist(err) {
		return true
	} else {
		return false
	}
}

func IsExistDir(dirPath string) bool {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IsExistFile(filePath string) bool {
	if fileInfo, err := os.Stat(filePath); err == nil {
		return fileInfo.Mode().IsRegular()
	} else {
		return false
	}
}

func IsSameAbsPath(filename1, filename2 string) bool {
	var (
		abs1, abs2 string
		err        error
	)
	if abs1, err = Path.Abs(filename1); err != nil {
		return false
	}
	if abs2, err = Path.Abs(filename2); err != nil {
		return false
	}
	return abs1 == abs2
}

func ListFilesSuffix(dir, suffix string) []string {
	var res = make([]string, 0)
	valid := IsExistDir(dir)
	if !valid {
		return res
	}
	if fileInfos, err := ioutil.ReadDir(dir); err == nil {
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				file := fileInfo.Name()
				if strings.HasSuffix(file, suffix) {
					res = append(res, file)
				}
			}
		}
	}
	return res
}

func ListFilesPrefix(dir, prefix string) []string {
	var res = make([]string, 0)
	valid := IsExistDir(dir)
	if !valid {
		return res
	}
	if fileInfoList, err := ioutil.ReadDir(dir); err == nil {
		for _, fileInfo := range fileInfoList {
			if !fileInfo.IsDir() {
				file := fileInfo.Name()
				if strings.HasPrefix(file, prefix) {
					res = append(res, file)
				}
			}
		}
	}
	return res
}

func ForEachLine(filename string, lineHandler func(line string) error, errHandler func(line string, err error) error) (successfulNum, failedNum int64, err error) {
	var fd *os.File
	fd, err = os.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		err = fd.Close()
	}()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		err = lineHandler(line)
		if err != nil {
			failedNum++
			err = errHandler(line, err)
			if err != nil {
				break
			}
		} else {
			successfulNum++
		}
	}
	return
}

func ForEachLineConcurrent(filename string, lineHandler func(line string) error, concurrent int) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	var waitGroup sync.WaitGroup
	jobs := make(chan bool, concurrent)
	started := make(chan bool)
	for scanner.Scan() {
		waitGroup.Add(1)
		line := scanner.Text()
		jobs <- true
		go func() {
			started <- true
			defer func() {
				waitGroup.Done()
				<-jobs

			}()
			_ = lineHandler(line)
		}()
		<-started
	}
	waitGroup.Wait()
	return err
}
