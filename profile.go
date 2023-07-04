package gorundut

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/pprof"
)

type RuntimeDumpResult struct {
	HeapDumpPath   string
	GoroutinesPath string
}

type RuntimeProfileName string

const (
	GoroutineProfileName    RuntimeProfileName = "goroutine"
	ThreadcreateProfileName RuntimeProfileName = "threadcreate"
	HeapProfileName         RuntimeProfileName = "heap"
	AllocsProfileName       RuntimeProfileName = "allocs"
	BlockProfileName        RuntimeProfileName = "block"
	MutexProfileName        RuntimeProfileName = "mutex"
)

func MakeProfileDump(directory string, profile RuntimeProfileName) (string, error) {
	profilePath := path.Join(directory, fmt.Sprintf("%v.pprof", profile))
	directoryPath := path.Dir(profilePath)
	if _, err := os.Stat(directoryPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("unable to create directory %v: %w", directoryPath, err)
		}
	}
	file, err := os.OpenFile(profilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return "", fmt.Errorf("unable to open file %v for profile %v: %w", profilePath, profile, err)
	}
	defer file.Close()
	err = pprof.Lookup(string(profile)).WriteTo(file, 0)
	if err != nil {
		return "", fmt.Errorf("unable to write profile %v: %w", profile, err)
	}
	return profilePath, nil
}
