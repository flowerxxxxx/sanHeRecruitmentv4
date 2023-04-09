package backupModule

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func fileCopy(from, to string) error {
	var err error

	f, err := os.Stat(from)
	if err != nil {
		return err
	}

	fn := func(fromFile string) error {
		//复制文件的路径
		rel, err := filepath.Rel(from, fromFile)
		if err != nil {
			return err
		}
		toFile := filepath.Join(to, rel)

		//创建复制文件目录
		if err = os.MkdirAll(filepath.Dir(toFile), 0777); err != nil {
			return err
		}

		//读取源文件
		file, err := os.Open(fromFile)
		if err != nil {
			return err
		}

		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建复制文件用于保存
		out, err := os.Create(toFile)
		if err != nil {
			return err
		}

		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, err = io.Copy(out, bufReader)
		return err
	}

	//转绝对路径
	pwd, _ := os.Getwd()
	if !filepath.IsAbs(from) {
		from = filepath.Join(pwd, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}

	//复制
	if f.IsDir() {
		return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return fn(path)
			} else {
				if err = os.MkdirAll(path, 0777); err != nil {
					return err
				}
			}
			return err
		})
	} else {
		return fn(from)
	}
}

func linuxFileCp(srcFolder, destFolder string) (err error) {
	//srcFolder := "copyUtil/from/path"
	//destFolder := "copyUtil/to/path"
	cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
	err = cpCmd.Run()
	return
}
