package storage

import (
	"github.com/google/logger"
	"os"
	"path"
	"path/filepath"
)

type StorageType struct {
	RootFolder   string
	ExportFolder string
	debug        bool
}

func NewStorage(root string, debug bool) *StorageType {
	storage := StorageType{
		RootFolder: root,
		debug:      debug,
	}
	return &storage
}

func (storage *StorageType) CreateOutputFolder() bool {
	storage.ExportFolder = storage.RootFolder

	logger.Info("The output folder will be ", storage.ExportFolder)

	if _, error := os.Stat(storage.ExportFolder); os.IsNotExist(error) {
		if err := os.MkdirAll(storage.ExportFolder, os.ModePerm); err != nil {
			return false
		}
	}

	return true
}

func (storage *StorageType) DownloadPath(meta map[string]string, formatID string) string {
	return filepath.Join(storage.ExportFolder, meta["id"]+"-"+formatID, meta["filename"])
}

func (storage *StorageType) ConvertedDownloadPath(filePath string) string {
	folder := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	// extension := filepath.Ext(filename)

	return filepath.Join(folder, "aac-" + filename)
}

func (storage *StorageType) EnsureFolder(filepath string) bool {
	pathname := path.Dir(filepath)

	if _, error := os.Stat(pathname); os.IsNotExist(error) {
		if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
			return false
		}
	}
	return true
}
