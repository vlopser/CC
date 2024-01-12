package main

import (
	// . "cc/pkg/lib/QueueManager"
	// . "cc/pkg/lib/StoreManager"
	. "cc/pkg/lib/TaskManager"

	"cc/worker/utils"

	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/nats-io/nats.go"
)

func goModTidy(root_dir string, file_errors *os.File) error {
	cmd := exec.Command("go", "mod", "tidy")

	cmd.Dir = root_dir + task.REPO_DIR
	cmd.Stderr = file_errors

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when executing 'go mod tidy':", err)
		return err
	}

	if isEmpty, err := utils.IsFileEmpty(file_errors); !isEmpty {
		return err
	}

	return nil
}

func waitForTasks(nats_server *nats.Conn, wg *sync.WaitGroup) {
	GetTasks(nats_server, handleRequest)

	wg.Add(1)

	log.Println("Waiting for request. (Presiona Ctrl+C para salir)")

	wg.Wait()                   // Esperar a que se reciba una señal de interrupción
	time.Sleep(1 * time.Second) // Permitir tiempo para desuscribirse antes de salir
}

func goRun(go_file string, output_file *os.File, error_file *os.File) error {
	cmd := exec.Command("go", "run", go_file)

	cmd.Stdout = output_file
	cmd.Stderr = error_file

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when executing `go run", go_file, "':", err)
		return err
	}

	return nil
}

func executeTask(root_dir string) error {

	file := utils.OpenFile(root_dir + task.RESULT_DIR + task.OUTPUT_FILE)
	defer file.Close()

	file_errors := utils.OpenFile(root_dir + task.RESULT_DIR + task.ERRORS_FILE)
	defer file_errors.Close()

	err := goModTidy(root_dir, file_errors)
	if err != nil {
		return err
	}

	err = goRun(root_dir+task.REPO_DIR+"/main.go", file, file_errors)

	return nil
}

func cloneRepo(repo_url string, root_dir string) error {

	//git.Plain cant clone if dir already exists, so deletes it if so
	err := utils.CheckDirectoryExists(root_dir + task.REPO_DIR)

	_, err = git.PlainClone(root_dir+task.REPO_DIR, false, &git.CloneOptions{
		URL:      repo_url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("Error doing git clone:", err)
		return err
	}

	return nil
}

func manageTask(t task.Task) error {

	task_dir := t.TaskId.String()
	utils.CreateDirectories(task_dir)

	err := cloneRepo(t.Input, task_dir)
	if err != nil {
		file_errors, _ := os.OpenFile(t.TaskId.String()+task.RESULT_DIR+task.ERRORS_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		defer file_errors.Close()
		file_errors.WriteString("Error cloning repo: " + err.Error())

		return err
	}

	// If an error occured, error file is automatically created in this method
	err = executeTask(task_dir)
	if err != nil {
		return err
	}

	return nil
}

func handleRequest(task task.Task, nats_server *nats.Conn) {
	log.Println("Request", task.TaskId.String(), "received!")

	SetTaskStatusToExecuting(nats_server, task.TaskId.String())

	init := time.Now()
	err := manageTask(task)
	end := time.Now()

	res := result.Result{
		TaskId:      task.TaskId,
		TimeElapsed: end.Sub(init),
	}

	if err != nil {
		res.Files = []string{task.TaskId.String() + "/result/errors.txt"}
		PostResult(nats_server, res)
		SetTaskStatusToFinishedWithErrors(nats_server, task.TaskId.String())
	} else {
		res.Files = []string{task.TaskId.String() + "/result/output.txt"}
		PostResult(nats_server, res)
		SetTaskStatusToFinished(nats_server, task.TaskId.String())
	}

	utils.CleanDirectory(task.TaskId.String()) //Clean all tmp directories created for the task

	//CODIGO DEL FRONTEND
	// result, err := GetResult(nats_server, task.TaskId.String())
	// if err != nil {
	// 	log.Println("Error when storing the result of", task.TaskId, ":", err)
	// 	return
	// }
	// println("Solucion en ", result.Files)
}

func main() {

	nats_server, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	wg := sync.WaitGroup{}
	go utils.WaitForSigkill(&wg)

	GetTasks(nats_server, handleRequest)

	// task := task.Task{
	// 	TaskId: uuid.New(),
	// 	Input:  "https://github.com/go-training/helloworld.git",
	// 	Status: task.PENDING,
	// }

	// EnqueueTask(task, nats_server)

	waitForTasks(nats_server, &wg)
}
