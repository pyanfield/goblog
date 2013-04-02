// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"github.com/pyanfield/goblog/archives"
	"github.com/pyanfield/goblog/blogs"
	"github.com/pyanfield/goblog/fs"
	"github.com/pyanfield/goblog/rss"
	"github.com/pyanfield/goblog/tags"
	"github.com/pyanfield/goblog/templates"
	"go/build"
	"log"
	"os"
	"path"
)

var (
	LOG   = log.New(os.Stderr, "---", log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = LOG
)

const (
	// 本框架路径
	GOBLOG_IMPORT_PATH = "github.com/pyanfield/goblog"
	//默认情况下会在 GOPATH下面的 goblog_workspace下面去创建文件夹来存放我们所有的文件
	GOBLOG_CONTENT = "goblog_workspace"
)

func main() {
	// Parse the flags.
	flag.Parse()
	setupDirectories()

	// First load the templates.
	// 返回的是 tmplts 是map[string]*template.Template，一个以模版文件名字为key值的Template的map
	tmplts, err := templates.LoadTemplates(TemplateDir)
	if err != nil {
		fmt.Println("loading templates:", err)
		os.Exit(1)
	}

	// Now, move the static files over.
	// 复制 static 文件夹下所有的子文件夹和子文件到 public 文件夹下
	err = fs.CopyFilesRecursively(OutputDir, StaticDir)
	if err != nil {
		fmt.Println("making output dir:", err)
		os.Exit(1)
	}

	// Get a list of files from the BlogDir.
	// 得到Blog文件夹下的所有md文件列表，如果在Blog下有子文件夹，那么这个文件夹的名字作为前缀，以"-"为连接符，形成新的文件名
	entries, err := blogs.GetBlogFiles(BlogDir)
	if err != nil {
		fmt.Println("getting blog file list:", err)
		os.Exit(1)
	}

	// Iteratively Parse each blog for it's useful data and generate a
	// page for each blog.
	for _, blog := range entries {
		contents, err := blog.Parse()
		if err != nil {
			fmt.Println("parsing blog", blog, ":", err)
			os.Exit(1)
		}

		err = tmplts.MakeBlogEntry(OutputDir, blog, contents)
		if err != nil {
			fmt.Println("generating blog html", blog, ":", err)
			os.Exit(1)
		}
	}

	// Generate the about page.
	err = tmplts.MakeAbout(OutputDir)
	if err != nil {
		fmt.Println("generating about.html:", err)
		os.Exit(1)
	}

	// Get a sorted list of tags.
	tagentries := tags.ParseBlogs(entries)
	t := tagentries.Slice()

	// Generate the tags page.
	err = tmplts.MakeTags(OutputDir, t)
	if err != nil {
		fmt.Println("generating tags.html:", err)
		os.Exit(1)
	}

	// Get a sort list of archives.
	dateentries := archives.ParseBlogs(entries)
	a := dateentries.Slice()

	// Generate the archive page.
	err = tmplts.MakeArchive(OutputDir, a)
	if err != nil {
		fmt.Println("generating archive.html:", err)
		os.Exit(1)
	}

	// Generate the index page.
	mostRecent := archives.GetMostRecent(a, MaxIndexEntries)
	err = tmplts.MakeIndex(OutputDir, mostRecent)
	if err != nil {
		fmt.Println("generating index.html:", err)
		os.Exit(1)
	}

	// Generate the RSS feed.
	err = rss.MakeRss(archives.GetMostRecent(a, 10), URL, TemplateDir, OutputDir)
	if err != nil {
		fmt.Println("generating feed.rss:", err)
		fmt.Println("no rss will be available")
	}

}

// SetupDirectories is a helper function that prepends the working
// directory to the other directories.
func setupDirectories() {
	// 检查是否通过命令输入了 WorkingDir ，如果没有，那就按照 GOBLOG_CONTENT建立博客的文件结构
	if WorkingDir == "./" {
		src := findSrcPath()
		if err := fs.MakeDirIfNotExists(GOBLOG_CONTENT); err != nil {
			ERROR.Fatalln("make diretory blog content:", err)
		}
		WorkingDir = path.Join(src, GOBLOG_CONTENT)
	}
	OutputDir = path.Join(WorkingDir, OutputDir)
	// Next, let's clear out the OutputDir if requested.
	if EmptyOutputDir {
		// 删除当前路径及其路径下的所有子文件。如果这个路径不存在将返回nil
		err := os.RemoveAll(OutputDir)
		if err != nil {
			fmt.Println("cleaning output dir:", err)
			os.Exit(1)
		}
	}
	if err := fs.MakeDirIfNotExists(OutputDir); err != nil {
		ERROR.Fatalln(err)
	}
	TemplateDir = path.Join(WorkingDir, TemplateDir)
	if err := fs.MakeDirIfNotExists(TemplateDir); err != nil {
		ERROR.Fatalln(err)
	}
	StaticDir = path.Join(WorkingDir, StaticDir)
	if err := fs.MakeDirIfNotExists(StaticDir); err != nil {
		ERROR.Fatalln(err)
	}
	BlogDir = path.Join(WorkingDir, BlogDir)
	if err := fs.MakeDirIfNotExists(BlogDir); err != nil {
		ERROR.Fatalln(err)
	}
}

//查找本项目的源地址
func findSrcPath() string {
	// 检查是否定义了 GOPATH 环境变量
	if gopath := os.Getenv("GOPATH"); gopath == "" {
		ERROR.Fatalln("Please set GOPATH first")
	}
	goblogPkg, err := build.Import(GOBLOG_IMPORT_PATH, "", build.FindOnly)
	// LOG.Println("appPkg >> ", goblogPkg.SrcRoot)
	if err != nil {
		ERROR.Fatalln("Failed to import", "", "with error:", err)
	}
	return goblogPkg.SrcRoot
}
