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
// 根据提供的路径获取创建和最后一次修改的时间。
// 首先是去尝试去获得第一次和最后一次git commit的时间。如果失败了，就再去使用系统的时间。
func GetTimes(p string) (time.Time, time.Time, error) {

	// First try to get the git times.
	// 去获取 git 里的第一次提交和最后一次提交的时间，第一次提交作为创建时间，最后一次为修改时间。
	first, err1 := getGitTime(p, true)
	second, err2 := getGitTime(p, false)

	if err1 == nil && err2 == nil {
		// We got the times, so let's return those!
		// 返回 Unix 时间
		return time.Unix(first, 0), time.Unix(second, 0), nil
	}

	// One of them must have failed, so let's get the system times.
	// 如果获取 git 提交时间失败的话，那么去获取系统时间
	// 先去抓取文件的 FileInfo信息
	fi, err := os.Stat(p)
	if err != nil {
		return time.Now(), time.Now(), nil
	}

	return fi.ModTime(), fi.ModTime(), nil
}

// getGitTime calls various git commands to get the UNIX Epoch time
// the file was made. if first is true, the original commit time is
// returned, otherwise, the most recent commit time is returned.
// 返回 git 第一个提交和最后一个提交的时间戳。
// 如果 first 为true，则返回的是第一个提交的时间，否则则为最后一次提交的时间
func getGitTime(p string, first bool) (int64, error) {
	//将路径分离成文件路径和文件名字两部分
	dir, file := path.Split(p)

	// Get our current directory and defer jumping back to it.
	// 返回当前文件夹的路径，及 ./goblog 的路径
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	// 无论如何，最后都要回复当前的工作地址
	defer os.Chdir(cwd)

	// Change directories.
	// 修改当前的 working diretory 到 dir 路径下
	err = os.Chdir(dir)
	if err != nil {
		return 0, err
	}

	// Call the revlist
	// 执行 git rev-list --max-parents=1 HEAD file的命令
	// git rev-list 有点类似git log,但是他会输出从一个commit 的SHA-1到另一个commit的 SHA-1
	// 的所有commit的 SHA值
	cmd := exec.Command("git", "rev-list", "--max-parents=1", "HEAD", file)
	// 返回执行命令之后的结果
	revs, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Get the first or last item.
	// 将执行命令返回的结果放入到 revslist里面，一边查找第一个和最后一个值
	revslist := strings.Split(string(revs), "\n")
	var rev string
	if !first {
		rev = revslist[0]
	} else {
		// The last one will actually be an empty line, so we want the
		// second to last.
		// 因为在revslist里面的最后一个值是空白行，所以我们想要的真正的结果是前面的非空白行的值
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
	// %at: author date, UNIX timestamp
	// 通过 git show -s --format=%at rev命令来得到想要的提交时间戳
	cmd = exec.Command("git", "show", "-s", "--format=%at", rev)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Convert to int.
	// 转化成十进制的 int64格式
	i, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil

}
