package store

import "os"

type FileStore struct {
	staged []stagedFile
	committed []string
}

type stagedFile struct {
	path string
	data []byte
}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (fs *FileStore) Stage(path string, data []byte) {
	fs.staged = append(fs.staged, stagedFile{path: path, data: data})
}

func (fs *FileStore) Flush() error {
	for _, file := range fs.staged {
		if fs.fileExists(file.path) {
			continue
		}
		if err := os.WriteFile(file.path, file.data, 0o644); err != nil {
			return err
		}
		fs.committed = append(fs.committed, file.path)
	}
	fs.staged = nil
	return nil
}

func (fs *FileStore) Rollback() error {
	for _, path := range fs.committed {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	fs.committed = nil
	return nil
}

func (fs *FileStore) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
