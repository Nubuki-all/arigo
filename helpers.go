package arigo

import (
	"fmt"
	"path/filepath"
    "strings"
    "os"
)


func DeleteControlFile(status Status) error {
	name, err := GetDownloadName(status)
	if err != nil {
		return err
	}
	ctrlFile := filepath.Join(status.Dir, name + ".aria2")
	return os.Remove(ctrlFile)
}

func GetDownloadName(status Status) (string, error) {
	var name string
	name = status.BitTorrent.Info.Name
	if name != "" {
		return name, nil
	}
	files := status.Files
	if len(files) > 0 {
		file := files[0]
		name = file.Path
		if strings.HasPrefix(name, "[METADATA]") {
			return name, nil
		}
		filePath, err := filepath.Abs(name)
		if err != nil {
			return name, err
		}
		dirPath, err := filepath.Abs(status.Dir)
		if err != nil {
			return name, err
		}
		if strings.HasPrefix(filePath, dirPath) {
			name = filepath.Base(filePath)
		} else {
			if uris := file.URIs; len(uris) > 0 {
				tempStr := strings.Split(uris[0].URI, "/")
				if len(tempStr) > 0 {
					name = tempStr[len(tempStr)-1]
				}
			}
		}
	}
	return name, nil
}

func RemoveFiles(files []File) {
	for _, file := range files {
		_ = os.Remove(file.Path)
	}
}


type Error struct {
	Code    ExitStatus `json:"code,string"`   // The code of the last error for this item, if any.
	Message string     `json:"message"`       // The human readable error message associated to ErrorCode

}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}