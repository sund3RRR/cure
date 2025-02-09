package types

import (
	"os"
	"path/filepath"
	"strings"
)

type FileType int8

const (
	RegularFile FileType = iota
	DirectoryFile
	SymlinkFile
)

type ErrorHandler func(ftype FileType, file File, err error) error

type WriteParams struct {
	Prefix Path
	Owner  string
	Mode   os.FileMode
}

type File interface {
	Write(p WriteParams, errHandler ErrorHandler) (Path, error)
	GetPath() Path
	GetType() FileType
}

type Regular struct {
	Type    FileType // File type [Regular, Directory, Symlink]
	Name    string
	Owner   string
	Path    Path
	Mode    os.FileMode // File mode (e.g. 0755)
	Content []byte      // for regular file
}

func NewRegularFile(path Path, content []byte) File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &Regular{
		Type:    RegularFile,
		Path:    path,
		Name:    name,
		Content: content,
	}
}

func (r *Regular) GetPath() Path {
	return r.Path
}

func (r *Regular) GetType() FileType {
	return r.Type
}

func (r *Regular) Write(p WriteParams, errHandler ErrorHandler) (Path, error) {
	err := os.WriteFile(filepath.Join(p.Prefix.String(), r.Path.String()), r.Content, r.Mode)
	return r.GetPath(), errHandler(r.Type, r, err)
}

type Directory struct {
	Type  FileType // File type [Regular, Directory, Symlink]
	Name  string
	Owner string
	Path  Path
	Mode  os.FileMode // File mode (e.g. 0755)
	Files []File      // for directory
}

func NewDirectory(path Path, files ...File) File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &Directory{
		Type:  DirectoryFile,
		Path:  path,
		Name:  name,
		Files: files,
	}
}

func (d *Directory) GetPath() Path {
	return d.Path
}

func (d *Directory) GetType() FileType {
	return d.Type
}

func (d *Directory) Write(p WriteParams, errHandler ErrorHandler) (Path, error) {
	writeDirFunc := func(f File) error {
		err := os.MkdirAll(filepath.Join(p.Prefix.String(), d.Path.String()), p.Mode)
		if err := errHandler(d.Type, f, err); err != nil {
			return err
		}
		for _, file := range d.Files {
			_, err := file.Write(p, errHandler)
			if err := errHandler(d.Type, file, err); err != nil {
				return err
			}
		}
		return nil
	}

	return d.Path, writeDirFunc(d)
}

type Symlink struct {
	Type    FileType // File type [Regular, Directory, Symlink]
	Name    string
	Owner   string
	Path    Path
	PointTo Path // for symlink
}

func NewSymlink(path Path, pointTo Path) File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &Symlink{
		Type:    SymlinkFile,
		Path:    path,
		Name:    name,
		PointTo: pointTo,
	}
}

func (s *Symlink) GetPath() Path {
	return s.Path
}

func (s *Symlink) GetType() FileType {
	return s.Type
}

func (s *Symlink) Write(p WriteParams, errHandler ErrorHandler) (Path, error) {
	err := os.Symlink(s.PointTo.String(), filepath.Join(p.Prefix.String(), s.Path.String()))
	return s.GetPath(), errHandler(s.Type, s, err)
}
