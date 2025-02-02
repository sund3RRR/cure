package types

import (
	"os"
	"path/filepath"
	"strings"
)

type FileType int8

const (
	Regular FileType = iota
	Directory
	Symlink
)

type File struct {
	Type    FileType // File type [Regular, Directory, Symlink]
	Name    string
	Owner   string
	Path    Path
	Mode    os.FileMode // File mode (e.g. 0755)
	PointTo Path        // for symlink
	Content []byte      // for regular file
	Files   []*File     // for directory
}

type ErrorHandler func(ftype FileType, file *File, err error) error
type WriteParams struct {
	Prefix Path
	Owner  string
	Mode   os.FileMode
}

func NewFile(path Path, content []byte) *File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &File{
		Type:    Regular,
		Path:    path,
		Name:    name,
		Content: content,
	}
}

func NewDirectory(path Path, files ...*File) *File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &File{
		Type:  Directory,
		Path:  path,
		Name:  name,
		Files: files,
	}
}

func NewSymlink(path Path, pointTo Path) *File {
	splitted := strings.Split(path.String(), "/")
	name := ""
	if len(splitted) > 0 {
		name = splitted[len(splitted)-1]
	}
	return &File{
		Type:    Symlink,
		Path:    path,
		Name:    name,
		PointTo: pointTo,
	}
}

func (f *File) Write(p WriteParams, errHandler ErrorHandler) (Path, error) {
	switch f.Type {
	case Regular:
		err := os.WriteFile(filepath.Join(p.Prefix.String(), f.Path.String()), f.Content, f.Mode)
		return f.GetPath(), errHandler(f.Type, f, err)
	case Directory:
		writeDirFunc := func(f *File) error {
			err := os.MkdirAll(filepath.Join(p.Prefix.String(), f.Path.String()), p.Mode)
			if err := errHandler(f.Type, f, err); err != nil {
				return err
			}
			for _, file := range f.Files {
				_, err := file.Write(p, errHandler)
				if err := errHandler(f.Type, f, err); err != nil {
					return err
				}
			}
			return nil
		}
		return f.Path, writeDirFunc(f)
	case Symlink:
		err := os.Symlink(f.PointTo.String(), filepath.Join(p.Prefix.String(), f.Path.String()))
		return f.GetPath(), errHandler(f.Type, f, err)
	default:
		return "", ErrInvalidFileType
	}
}

func (f *File) GetPath() Path {
	return f.Path
}

func (f *File) GetType() FileType {
	return f.Type
}
