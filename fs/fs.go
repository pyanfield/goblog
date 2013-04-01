// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package fs contains functions that access or modify the filesystem
// that are useful to the goblog application.
package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// MakeDirIfNotExists creates the given directory if it does not exist.
// 创建一个文件夹如果这个文件夹不存在
func MakeDirIfNotExists(dir string) error {
	// 得到所给路径 dir的 FileInfo
	st, err := os.Stat(dir)
	if err != nil {
		// It may just not exist. Make it or error out trying.
		// 如果文件夹不存在就创建一个文件夹
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0750)
		}
	}

	// It exists, make sure it's a directory.
	// 不是文件夹提示错误
	if !st.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	return nil
}

// CopyFilesRecursively copies the contents of the directory src
// into dest.
// 通过递归调用，将 src 下的文件子文件全都复制到 dest下面
func CopyFilesRecursively(dest, src string) error {

	// Read the list of entries for src.
	// 返回src下面的所有子文件夹,子文件信息
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		s := path.Join(src, file.Name())
		d := path.Join(dest, file.Name())

		// If the file is a directory, make it and then recursively call
		// this function.
		// 如果 file 是文件，那么建立一个同名的文件夹，然后递归这个文件夹下
		// 如果是文件，那么就复制这个文件
		if file.IsDir() {
			err = MakeDirIfNotExists(d)
			if err != nil {
				return err
			}

			err = CopyFilesRecursively(d, s)
			if err != nil {
				return err
			}
		} else {
			// If the file is a file, then copy the file to dest.
			CopyFile(d, s)
		}

		// Set the create/mod times to be the same as the src.
		// 修改目标文件文件的创立时间和修改时间与源文件的相同
		err = os.Chtimes(d, file.ModTime(), file.ModTime())
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy file makes an exact copy fo the file at src and saves it to
// dest. The contents of dest are overwritten if it exists.
// 复制源文件到目标文件夹中，如果这个文件已经存在，那么覆盖这个文件
func CopyFile(dest, src string) error {
	// 读取 src 文件，返回 *File
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 建立文件，且mode为 0666, 返回 *File
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// GetTimes attempts to get the create and last modify times of the
// given path. It first tries using git for the first and last
// commits. If that fails, it uses the system time.
func GetTimes(p string) (time.Time, time.Time, error) {

	// First try to get the git times.
	first, err1 := getGitTime(p, true)
	second, err2 := getGitTime(p, false)

	if err1 == nil && err2 == nil {
		// We got the times, so let's return those!
		return time.Unix(first, 0), time.Unix(second, 0), nil
	}

	// One of them must have failed, so let's get the system times.
	fi, err := os.Stat(p)
	if err != nil {
		return time.Now(), time.Now(), nil
	}

	return fi.ModTime(), fi.ModTime(), nil
}

// getGitTime calls various git commands to get the UNIX Epoch time
// the file was made. if first is true, the original commit time is
// returned, otherwise, the most recent commit time is returned.
func getGitTime(p string, first bool) (int64, error) {
	dir, file := path.Split(p)

	// Get our current directory and defer jumping back to it.
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	defer os.Chdir(cwd)

	// Change directories.
	err = os.Chdir(dir)
	if err != nil {
		return 0, err
	}

	// Call the revlist
	cmd := exec.Command("git", "rev-list", "--max-parents=1", "HEAD", file)
	revs, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Get the first or last item.
	revslist := strings.Split(string(revs), "\n")
	var rev string
	if !first {
		rev = revslist[0]
	} else {
		// The last one will actually be an empty line, so we want the
		// second to last.
		if len(revslist) < 2 {
			rev = revslist[0]
		} else {
			rev = revslist[len(revslist)-2]
		}
	}

	if rev == "" {
		return 0, fmt.Errorf("not found")
	}

	// Get the commit time.
	cmd = exec.Command("git", "show", "-s", "--format=%at", rev)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Convert to int.
	i, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil

}
