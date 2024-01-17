package main

import (
	. "cc/src/pkg/lib/TaskManager"
	"cc/src/services/worker/utils"
	"context"

	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/nats-io/nats.go"
)

func waitForTasks(nats_server *nats.Conn, wg *sync.WaitGroup) {
	GetTasks(nats_server, handleRequest)

	wg.Add(1)

	log.Println("Waiting for request. (Presiona Ctrl+C para salir)")

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func execCommand(root_dir string, stdout_file *os.File, stderr_file *os.File, command string) error {
	// Add command executed to files
	stdout_file.WriteString(">> " + command + "\n")
	stderr_file.WriteString(">> " + command + "\n")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	GidMappings: "grupo_sin_permisos",
	// }

	cmd.Stdout = stdout_file
	cmd.Stderr = stderr_file

	cmd.Dir = root_dir

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when executing `", command, "':", err)
		return err
	}

	return nil
}

func executeTask(task_dir string) error {

	stdout_file := utils.OpenFile(task_dir + task.RESULT_DIR + task.STDOUT_FILE)
	defer stdout_file.Close()

	stderr_file := utils.OpenFile(task_dir + task.RESULT_DIR + task.STDERR_FILE)
	defer stderr_file.Close()

	err := execCommand(task_dir+task.REPO_DIR, stdout_file, stderr_file, "go mod download")
	// if err != nil {  // Quizás no puede hacer go mod download porqeu no hay .mod, dejemos que go run decida si puede ejecutarse
	// 	return err
	// }

	err = execCommand(task_dir+task.REPO_DIR, stdout_file, stderr_file, "go run main.go")
	if err != nil {
		return err
	}

	return nil
}

func cloneRepo(repo_url string, task_dir string) error {

	//git.Plain cant clone if dir already exists, so deletes it if so
	err := utils.CheckDirectoryExists(task_dir + task.REPO_DIR)

	_, err = git.PlainClone(task_dir+task.REPO_DIR, false, &git.CloneOptions{
		URL:      repo_url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("Error doing git clone:", err)
		return err
	}

	return nil
}

func handleRequest(t task.Task, nats_server *nats.Conn) {
	log.Println("Request", t.TaskId.String(), "received!")

	utils.CreateTaskDirectory(t.TaskId.String())

	err := cloneRepo(t.RepoUrl, t.TaskId.String())
	if err != nil {
		error_file := utils.CreateErrorFile(t.TaskId.String(), "Error cloning repo: "+err.Error())
		PostResult(
			nats_server,
			result.Result{
				TaskId: t.TaskId,
				Files:  []string{error_file},
			},
		)
		SetTaskStatusToFinishedWithErrors(nats_server, t.TaskId.String())
		utils.CleanDirectory(t.TaskId.String())
		return
	}

	SetTaskStatusToExecuting(nats_server, t.TaskId.String())

	init := time.Now()
	err = executeTask(t.TaskId.String())
	end := time.Now()

	if err != nil {
		error_file := utils.CreateErrorFile(t.TaskId.String(), "Error executing task: "+err.Error())
		PostResult(
			nats_server,
			result.Result{
				TaskId: t.TaskId,
				Files: []string{
					t.TaskId.String() + task.RESULT_DIR + task.STDERR_FILE, //Even if an error happend during execution, there may be still some stdout and/or stderr
					t.TaskId.String() + task.RESULT_DIR + task.STDOUT_FILE,
					error_file,
				},
			},
		)
		SetTaskStatusToFinishedWithErrors(nats_server, t.TaskId.String())
		utils.CleanDirectory(t.TaskId.String())
		return
	}

	PostResult(
		nats_server,
		result.Result{
			TaskId: t.TaskId,
			Files: []string{
				t.TaskId.String() + task.RESULT_DIR + task.STDERR_FILE,
				t.TaskId.String() + task.RESULT_DIR + task.STDOUT_FILE,
			},
			TimeElapsed: end.Sub(init),
		},
	)

	SetTaskStatusToFinished(nats_server, t.TaskId.String())
	utils.CleanDirectory(t.TaskId.String()) //Clean all tmp directories created for the task

	//CODIGO DEL FRONTEND
	// result, err := GetResult(nats_server, t.TaskId.String())
	// if err != nil {
	// 	log.Println("Error when storing the result of", t.TaskId, ":", err)
	// 	return
	// }
	// println("Solucion en ", result.Files)
}

func main() {

	time.Sleep(10 * time.Second)

	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS")) //nats.DefaultURL
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	wg := sync.WaitGroup{}
	go utils.WaitForSigkill(&wg)

	GetTasks(nats_server, handleRequest)

	// Codigo frontend
	// task := task.Task{
	// 	TaskId: uuid.New(),
	// 	RepoUrl:  "https://github.com/go-training/helloworld.git",
	// 	Status: task.PENDING,
	// }

	// EnqueueTask(task, nats_server)

	waitForTasks(nats_server, &wg)
}
