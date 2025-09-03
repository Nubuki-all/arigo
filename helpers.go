package arigo

import (
	"os"
	"path/filepath"
    "strings"
)


func DeleteControlFile(status Status) error {
	name, err := GetDownloadName(status)
	if err != nil {
		return err
	}
	ctrlFile := filepath.Join(status.Dir, name + ".aria2")
	return os.Remove(ctrlFile) // error can be ignored
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
		if strings.HasPrefix(name, status.Dir) {
			name = filepath.Base(name)
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


func RemoveUnselectedFiles(status Status) {
	for _, file := range status.Files {
		if !file.Selected{
			_ = os.Remove(file.Path)
		}
	}
	name, _ := GetDownloadName(status)
	_, _ = RemoveEmptyDirs(filepath.Join(status.Dir, name))
}


// RemoveEmptyDirs walks root recursively and removes any directories that are empty.
// It returns the number of directories removed and the first non-nil error encountered (if any).
// Symlinked directories are skipped (not followed).
func RemoveEmptyDirs(root string) (int, error) {
	var removed int

	// inner recursive function that returns (isEmpty bool, err error)
	var walk func(path string) (bool, error)
	walk = func(path string) (bool, error) {
		entries, err := os.ReadDir(path)
		if err != nil {
			return false, err
		}

		// iterate children first
		for _, e := range entries {
			if e.IsDir() {
				// skip symlinked directories
				info, err := e.Info()
				if err != nil {
					return false, err
				}
				if info.Mode()&os.ModeSymlink != 0 {
					// don't follow symlinks
					continue
				}
				childPath := filepath.Join(path, e.Name())
				empty, err := walk(childPath)
				if err != nil {
					return false, err
				}
				if empty {
					// try to remove it
					if err := os.Remove(childPath); err != nil {
						// if remove fails, we treat dir as non-empty for parent
						return false, err
					}
					removed++
				}
			}
		}

		// After children processed, re-read entries to see if any remain.
		entries, err = os.ReadDir(path)
		if err != nil {
			return false, err
		}
		return len(entries) == 0, nil
	}

	isEmpty, err := walk(root)
	if err != nil {
		return removed, err
	}
	// optionally remove the root itself if you want:
	if isEmpty {
		// don't remove root by default — uncomment if desired:
		// if err := os.Remove(root); err == nil {
		//     removed++
		// } else {
		//     return removed, err
		// }
	}
	return removed, nil
}