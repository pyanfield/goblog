// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package blog contains structures, methods and functions for
// manipulating blog entries.
package blogs

import (
	"bytes"
	"fmt"
	"github.com/pyanfield/goblog/fs"
	md "github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"
)

// BlogEntry is a representation of a blog entry. It contains the
// information necessary to generate the blogs contents and store
// information about the blog.
type BlogEntry struct {
	// Name is the name of the entry gleaned from the filename.
	Name string

	// Path is the path to the markdown file from the cwd.
	Path string

	// Aurhor is the name of the person who wrote the page.
	Author string

	// Title is the title of the Entry.
	Title string

	// Description is the description of the Entry.
	Description string

	// Url is the HTML file name of this entry (Name + ".html").
	Url string

	// Tags is a list of tags this blog entry contains. It is generated
	// when when the Parse method is called.
	Tags []string

	// Languages is a list of languages this blog entry contains. It is
	// generated when when the Parse method is called.
	Languages []string

	// Created is the date the blog entry was created. It is generated
	// when the Parse metod is called.
	Created time.Time

	// Updated is the date the blog entry was last updated. It is
	// generated when the Parse metod is called.
	Updated time.Time
}

// Parse reads the contents of the path for this BlogEntry. It gleans
// information from the file and saves it to this BlogEntry. It then
// formats the markdown to HTML and returns that.
func (be *BlogEntry) Parse() (string, error) {
	// Get the files contents.
	orgContents, err := ioutil.ReadFile(be.Path)
	if err != nil {
		return "", err
	}

	// Save some of the meta data.
	err = be.gleanInfo(string(orgContents))
	if err != nil {
		return "", err
	}

	// Return the markdown content.
	return string(md.MarkdownCommon(orgContents)), nil
}

// CDate is a helper function for the templating system that returns
// the Created date as a string or "" if there is no value.
func (be *BlogEntry) CDate() string {
	if be.Created.IsZero() {
		return ""
	}

	return be.Created.Format("2006-01-02")
}

// PubDate is a helper function for the templating system that returns
// the Created date as an RFC822 string or "" if there is no value.
func (be *BlogEntry) PubDate() string {
	if be.Created.IsZero() {
		return ""
	}

	return be.Created.Format(time.RFC822)
}

// UDate is a helper function for the templating system that returns
// the Updated date as a string or "" if there is no value or if it's
// identical to the Created date.
func (be *BlogEntry) UDate() string {
	if be.Updated.IsZero() {
		return ""
	}

	if be.Created.Equal(be.Updated) {
		return ""
	}

	return be.Updated.Format("2006-01-02")
}

// GetBlogFiles looks in the given directory for blog entries and
// returns a list of them. Blog entries must have the '.md'
// extension. Entries are searched in the directory recursively. If a
// files is in a directory, the directory name is used as a prefix to
// the blog entries name concatenated with a '-'. The blog is not
// parsed or read. You should do that yourself elsewhere.
// 返回 dir 文件夹下的 BlogEntry 的list,这些 Blog的文件必须是以 .md结尾. 如果Blog文件在一个子文件夹内，
// 那么这个子文件夹的名字会作为Blog的前缀，并且以 "-" 来作为文件夹和BLOG文件的连接. 
// BLOG 不会被解析和读取
func GetBlogFiles(dir string) ([]*BlogEntry, error) {
	entries := []*BlogEntry{}

	// Read the list of entries for dir.
	// 读取dir下的文件夹和文件信息
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 完整的本地路径
		p := path.Join(dir, file.Name())

		// If it's a directory, then recursively call this function and
		// merge the two slices.
		if file.IsDir() {
			blogs, err := GetBlogFiles(p)
			if err != nil {
				return nil, err
			}

			// Add the blog entries to this list.
			for _, blog := range blogs {
				newName, err := MakeBlogName(file.Name(), blog.Name)
				if err != nil {
					return nil, err
				}

				entries = append(entries, &BlogEntry{
					Name: newName,
					Url:  newName + ".html",
					Path: blog.Path,
				})

			}
		} else {
			// We only deal with .md (markdown) files.
			// 如果是文件，那么检测是否以.md 为文件扩展名的markdown文件
			if path.Ext(p) != ".md" {
				continue
			}

			// Get the name of the entry.
			// 将文件和其扩展名分割，只读取文件名
			pieces := strings.Split(file.Name(), ".md")
			newName := pieces[0]

			// Just create the new entry.
			// 生成新的BlogEntry，保存其文件名，带有新的扩展名html的URL和文件路径
			// 将其加入到 entries里面
			entries = append(entries, &BlogEntry{
				Name: newName,
				Url:  newName + ".html",
				Path: p,
			})
		}
	}

	return entries, nil
}

// MakeBlogName concatenates all of the given names with a dash and
// replaces illegal blog name characters with a dash.
// 将不合法的名字替换成 "-"，并且如果这个使用"-"作为文件夹与子文件的连接符
func MakeBlogName(names ...string) (string, error) {
	re := regexp.MustCompile("[^-a-zA-Z0-9_]")
	buf := new(bytes.Buffer)

	end := len(names) - 1
	for i, name := range names {
		// fmt.Println(name)
		_, err := buf.Write(re.ReplaceAll([]byte(name), []byte("-")))
		if err != nil {
			return "", err
		}
		// 如果是文件夹的话，那么在其后面添加 "-" 以作为和子文件的连接
		if i != end {
			_, err := buf.Write([]byte("-"))
			if err != nil {
				return "", err
			}

		}
	}
	return buf.String(), nil
}

// gleanInfo is a helper function that searches for various comments
// that contain useful information about the blog. It also uses fetchs
// the update and create dates.
func (be *BlogEntry) gleanInfo(contents string) error {
	/* These are the patterns we are searching for */
	var err error

	be.Title, err = regexSingle("Title", contents)
	if err != nil {
		return err
	}

	be.Author, err = regexSingle("Author", contents)
	if err != nil {
		return err
	}

	be.Description, err = regexSingle("Description", contents)
	if err != nil {
		return err
	}

	be.Languages, err = regexList("Languages", contents)
	if err != nil {
		return err
	}

	be.Tags, err = regexList("Tags", contents)
	if err != nil {
		return err
	}

	be.Created, be.Updated, err = fs.GetTimes(be.Path)
	if err != nil {
		return err
	}

	return nil
}

// regexList is a helper function that performs a regex search for an
// HTML comment with the given title. It returns the list (comma
// separated) of values.
func regexList(key, contents string) ([]string, error) {
	val, err := regexSingle(key, contents)
	if err != nil {
		return nil, err
	}

	if val == "" {
		return []string{}, nil
	}

	return strings.Split(val, ","), nil
}

// regexSingle is a helper function that performs a regex search for
// an HTML comment with the given title. It returns the value if it
// was found, or "" if it wasn't.
func regexSingle(key, contents string) (string, error) {
	pattern := "<!--[ ]*" + key + ":(.*)-->"

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(contents)
	if len(matches) < 2 {
		return "", nil
	}

	return strings.TrimSpace(matches[1]), nil
}
