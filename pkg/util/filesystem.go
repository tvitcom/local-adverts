package util

import (
	"os"
	"fmt"
	// "log"
	"time"
	"errors"
	"path/filepath"
    "context"
)

// Make new dir with name as the eid in current prevPath directory
// For example prevPath="upload/", uid="123"
// will make new directory "upload/123"
func MakeUploadDirByUserId(prevPath, uid string) (string, error) {
fmt.Println("TO MakeUploadDirByUserId:", prevPath, uid)
	newDir := prevPath + uid
	dir, err := os.Open(prevPath)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	_, err = os.Stat(newDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(newDir, 0766)
		if errDir != nil {
			return "", err
		}
	}
	return newDir + "/", nil
}

func FileDeletionByMask(ctx context.Context, fpathes string) error {
    cmdTimeout := 10 * time.Second

    if fpathes == "" {
        return errors.New("File path is empty")
    }
	
	ctx, cancel := context.WithTimeout(ctx, cmdTimeout)
	defer cancel()

	files, _ := filepath.Glob(fpathes)
	for _, file:= range files {
    	err := os.Remove(file)
    	if err != nil {
    		return err
    	}
	}
    return nil
   }