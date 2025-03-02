package stores

import (
	"os"
	"path/filepath"

	"github.com/tuxounet/k2/types"

	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

func getRunDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return cwd
}

func unmarshallFile[TItem any](filePath string) (*TItem, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result TItem
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

type FileStore struct {
	// Path to the directory where the files are stored
	Dir string
	K2  []*types.IK2[any]
}

// NewFileStore creates a new FileStore
func NewFileStore(dir string) *FileStore {

	if !filepath.IsAbs(dir) {
		dir = filepath.Join(getRunDir(), dir)
	}

	return &FileStore{
		Dir: dir,
		K2:  make([]*types.IK2[any], 0),
	}
}

func (fs *FileStore) GetKey(key string) (*types.IK2[any], error) {
	localPath := filepath.Join(fs.Dir, key)
	ret, err := unmarshallFile[types.IK2[any]](filepath.Join(fs.Dir, key))
	if err != nil {
		return nil, err
	}
	ret.K2.Metadata.Path = localPath
	ret.K2.Metadata.Folder = filepath.Dir(localPath)
	return ret, nil
}

func (fs *FileStore) GetAsInventory(key string) (*types.IK2Inventory, error) {
	localPath := filepath.Join(fs.Dir, key)
	ret, err := unmarshallFile[types.IK2Inventory](localPath)
	if err != nil {
		return nil, err
	}
	ret.K2.Metadata.Path = localPath
	ret.K2.Metadata.Folder = filepath.Dir(localPath)
	return ret, nil

}

func (fs *FileStore) GetAsTemplateApply(key string) (*types.IK2TemplateApply, error) {
	localPath := filepath.Join(fs.Dir, key)
	ret, err := unmarshallFile[types.IK2TemplateApply](localPath)
	if err != nil {
		return nil, err
	}
	ret.K2.Metadata.Path = localPath
	ret.K2.Metadata.Folder = filepath.Dir(localPath)
	return ret, nil

}
func (fs *FileStore) GetAsTemplate(key string) (*types.IK2Template, error) {
	localPath := filepath.Join(fs.Dir, key)
	ret, err := unmarshallFile[types.IK2Template](localPath)
	if err != nil {
		return nil, err
	}
	ret.K2.Metadata.Path = localPath
	ret.K2.Metadata.Folder = filepath.Dir(localPath)
	return ret, nil

}

func (fs *FileStore) Scan(patterns []string) ([]string, error) {

	var files []string
	err := filepath.Walk(fs.Dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, file := range files {
		relative, err := filepath.Rel(fs.Dir, file)
		if err != nil {
			return nil, err
		}
		for _, pattern := range patterns {

			matched, err := glob.Compile(pattern)
			if err != nil {
				return nil, err
			}
			if matched.Match(relative) {
				result = append(result, relative)
				break
			}

		}
	}

	return result, nil

}
