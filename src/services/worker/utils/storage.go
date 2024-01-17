package utils

import (
	"cc/src/pkg/models/task"
	"fmt"
	"log"
	"os"
)

func CreateTaskDirectory(task_id string) error {

	err := os.Mkdir(task_id, 0755)
	if err != nil {
		fmt.Println("Error al crear el directorio:", err)
		return err
	}

	task_result_dir := task_id + task.RESULT_DIR
	err = os.Mkdir(task_result_dir, 0755)
	if err != nil {
		fmt.Println("Error al crear el directorio:", err)
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

func CreateErrorFile(task_id string, error_msg string) string {
	file_errors := OpenFile(task_id + task.RESULT_DIR + task.ERROR_FILE)
	defer file_errors.Close()
	file_errors.WriteString(error_msg)

	return file_errors.Name()
}

func IsFileEmpty(file *os.File) (bool, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	// Verificar si el tamaño del archivo es 0
	return fileInfo.Size() == 0, nil
}

func OpenFile(file_path string) *os.File {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		// return err
	}

	return file
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
