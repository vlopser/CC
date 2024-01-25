package utils

import (
	"archive/zip"
	"cc/src/pkg/models/task"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func CreateTaskDirectory(task_id string) error {
	log.Println("Creating task directory")

	err := os.Mkdir(task_id, 0755)
	if err != nil {
<<<<<<< HEAD:src/services/worker/utils/storage.go
		fmt.Println("Error creating directory:", err)
=======
>>>>>>> main:src/pkg/utils/storage.go
		return err
	}

	task_result_dir := task_id + task.RESULT_DIR
	err = os.Mkdir(task_result_dir, 0755)
	if err != nil {
<<<<<<< HEAD:src/services/worker/utils/storage.go
		fmt.Println("Error creating directory:", err)
=======
>>>>>>> main:src/pkg/utils/storage.go
		return err
	}

	return nil
}

func CleanDirectory(root_dir string) {
	err := os.RemoveAll(root_dir)
	if err != nil {
		log.Println(err)
	}
}

func RemoveAllFiles(files []string) error {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Println("Error removing file '", file, "':", err.Error())
			return err
		}
	}

	return nil
}

func CreateErrorFile(task_id string, error_msg string) (string, error) {
	log.Printf("Creating error file.\n")

	file_errors, err := OpenFile(task_id + task.RESULT_DIR + task.ERROR_FILE)
	if err != nil {
		return "", err
	}
	defer file_errors.Close()
	file_errors.WriteString(error_msg)

	return file_errors.Name(), nil
}

func IsFileEmpty(file *os.File) (bool, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	// Verificar si el tamaño del archivo es 0
	return fileInfo.Size() == 0, nil
}

func OpenFile(file_path string) (*os.File, error) {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CheckDirectoryExists(dir string) error {
	if _, err := os.Stat(dir + task.REPO_DIR); err == nil {
		err := os.RemoveAll(dir + task.REPO_DIR)
		if err != nil {
			log.Println("Error removing directory:", err)
			return err
		}
	}

	return nil
}

func CloneRepo(repo_url string, task_dir string) error {
	log.Printf("Cloning repo '%s' to dir '%s'.\n", repo_url, task_dir)

	//git.Plain cant clone if dir already exists, delete it if so
	CheckDirectoryExists(task_dir + task.REPO_DIR)

	_, err := git.PlainClone(task_dir+task.REPO_DIR, false, &git.CloneOptions{
		URL: repo_url,
	})
	if err != nil {
		return err
	}

	changePermissionsDir(task_dir, os.FileMode(0755))
	return nil
}

func changePermissionsDir(dir string, perm os.FileMode) error {
	err := os.Chmod(dir, perm)
	if err != nil {
		return err
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, perm)
	})

	if err != nil {
		return err
	}

	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// Get file information
	fileInfo, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	// Create file header zip
	fileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	fileHeader.Name = filename
	// Creat new file in zip
	fileInZip, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return err
	}
	// Copy the contents of the file to the new file in the zip
	_, err = io.Copy(fileInZip, fileToZip)
	if err != nil {
		return err
	}
	return nil
}
