package utils

import (
	"cc/src/pkg/models/task"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func CreateTaskDirectory(task_id string) error {

	err := os.Mkdir(task_id, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

	task_result_dir := task_id + task.RESULT_DIR
	err = os.Mkdir(task_result_dir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
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

func CreateErrorFile(task_id string, error_msg string) (string, error) {
	file_errors, err := OpenFile(task_id + task.RESULT_DIR + task.ERROR_FILE)
	if err != nil {
		fmt.Println("Unable to open file in '", task_id+task.RESULT_DIR+task.ERROR_FILE, "':", err)
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
	// Verificar si el directorio existe
	if _, err := os.Stat(dir + task.REPO_DIR); err == nil {
		// Si el directorio existe, lo eliminamos
		err := os.RemoveAll(dir + task.REPO_DIR)
		if err != nil {
			log.Println("Error removing directory:", err)
			return err
			// return err
		}
	}

	return nil
}

func CloneRepo(repo_url string, task_dir string) error {

	//git.Plain cant clone if dir already exists, so delete it if so
	CheckDirectoryExists(task_dir + task.REPO_DIR)

	_, err := git.PlainClone(task_dir+task.REPO_DIR, false, &git.CloneOptions{
		URL:      repo_url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("Error doing git clone:", err)
		return err
	}

	return nil
}
