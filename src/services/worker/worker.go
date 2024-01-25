package main

import (
	store "cc/src/pkg/lib/StoreManager"
	. "cc/src/pkg/lib/TaskManager"
	"cc/src/services/worker/utils"
	"context"
	"strings"

	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	MAX_TIME_EXECUTION = 10 * time.Second
)

func waitForTasks(nats_server *nats.Conn, wg *sync.WaitGroup) {
	ReceiveTasks(nats_server, handleRequest)

	wg.Add(1)

	log.Println("Waiting for requests. (Presiona Ctrl+C para salir)")

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func execCommand(root_dir string, stdout_file *os.File, stderr_file *os.File, command string, args ...string) error {
	command_with_args := command + " " + strings.Join(args, " ")

	stdout_file.WriteString(">> " + command_with_args + "\n")
	stderr_file.WriteString(">> " + command_with_args + "\n")

	ctx, cancel := context.WithTimeout(context.Background(), MAX_TIME_EXECUTION)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command_with_args)

	cmd.Stdout = stdout_file
	cmd.Stderr = stderr_file

	cmd.Dir = root_dir

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when executing '", command, "':", err)
		return err
	}

	return nil
}

func executeTask(t task.Task) error {

	task_dir := t.TaskId.String()

	stdout_file, _ := utils.OpenFile(task_dir + task.RESULT_DIR + task.STDOUT_FILE)
	defer stdout_file.Close()

	stderr_file, _ := utils.OpenFile(task_dir + task.RESULT_DIR + task.STDERR_FILE)
	defer stderr_file.Close()

	err := execCommand(task_dir+task.REPO_DIR, stdout_file, stderr_file, "go mod download")
	err = execCommand(task_dir+task.REPO_DIR, stdout_file, stderr_file, "go run main.go", t.Parameters...)
	if err != nil {
		return err
	}

	return nil
}

func handleRequest(t task.Task, nats_server *nats.Conn) {
	log.Println("Request", t.TaskId.String(), "received!")

	store.SetInOutMsgs(nats_server, store.OUT_MSGS)

	err := utils.CreateTaskDirectory(t.TaskId.String())
	if err != nil {
		log.Println("Error when creating repo directory at '", t.TaskId.String(), "':", err.Error())
		return
	}

	err = utils.CloneRepo(t.RepoUrl, t.TaskId.String())
	if err != nil {
		error_file, err := utils.CreateErrorFile(t.TaskId.String(), "Error cloning repo: "+err.Error())
		if err != nil {
			log.Println("Unable to create error file for task '", t.TaskId.String(), "':")
			SetTaskStatusToUnexpectedError(nats_server, t)
			return
		}

		err = CreateTaskResult(nats_server, result.Result{TaskId: t.TaskId, Files: []string{error_file}})
		if err != nil {
			log.Println("Unable to post result for task '", t.TaskId.String(), "':")
			SetTaskStatusToUnexpectedError(nats_server, t)
			return
		}

		SetTaskStatusToFinishedWithErrors(nats_server, t)
		return
	}
	defer utils.CleanDirectory(t.TaskId.String()) // Clean all tmp directories created for the task

	SetTaskStatusToExecuting(nats_server, t)

	init := time.Now()
	err = executeTask(t)
	end := time.Now()

	if err != nil {
		error_file, err := utils.CreateErrorFile(t.TaskId.String(), "Error executing task: "+err.Error())
		if err != nil {
			log.Println("Unable to create error file for task '", t.TaskId.String(), "':")
			SetTaskStatusToUnexpectedError(nats_server, t)
			return
		}

		err = CreateTaskResult(
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
		if err != nil {
			log.Println("Unable to post result for task '", t.TaskId.String(), "':")
			SetTaskStatusToUnexpectedError(nats_server, t)
			return
		}

		SetTaskStatusToFinishedWithErrors(nats_server, t)
		return
	}

	CreateTaskResult(
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
	if err != nil {
		log.Println("Unable to post result for task '", t.TaskId.String(), "':")
		SetTaskStatusToUnexpectedError(nats_server, t)
		return
	}

	SetTaskStatusToFinished(nats_server, t)
}

func main() {

	// time.Sleep(10 * time.Second)

	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	wg := sync.WaitGroup{}
	go utils.WaitForSigkill(&wg)

	waitForTasks(nats_server, &wg)
}
