package types

import (
	"os"
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

func (f *File) Write(mode os.FileMode, owner string) (Path, error) {
	switch f.Type {
	case Regular:
		return f.Path, os.WriteFile(f.Path.String(), f.Content, f.Mode)
	case Directory:
		writeDirFunc := func(f *File) error {
			for _, file := range f.Files {
				_, err := file.Write(mode, owner)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return f.Path, writeDirFunc(f)
	case Symlink:
		return f.Path, os.Symlink(f.PointTo.String(), f.Path.String())
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
