package main

import (
	store "cc/src/pkg/lib/StoreManager"
	. "cc/src/pkg/lib/TaskManager"
	"cc/src/pkg/utils"
	"context"
	"strconv"
	"strings"
	"syscall"

	"cc/src/pkg/models/result"
	"cc/src/pkg/models/task"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	MAX_TIME_EXECUTION = 20 * time.Second
)

func waitForTasks(nats_server *nats.Conn, wg *sync.WaitGroup) {
	ReceiveTasks(nats_server, handleRequest)

	wg.Add(1)

	log.Println("Waiting for requests. (Press Ctrl+C to exit)")

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func execCommand(root_dir string, stdout_file *os.File, stderr_file *os.File, command string, args ...string) error {
	command = command + " " + strings.Join(args, " ")
	final_command := "env GOMODCACHE=/home/client/go/pkg/mod GOCACHE='/home/client/.cache/go-build'" + " " + command //ENV needed for client user executions

	stdout_file.WriteString(">> " + command + "\n")
	stderr_file.WriteString(">> " + command + "\n")

	max_seconds_executing, err := strconv.Atoi(os.Getenv("MAX_SECONDS_EXECUTION_PER_REQUEST"))
	if err != nil {
		log.Fatal("Error using environment value MAX_SECONDS_EXECUTION_PER_REQUEST:", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(max_seconds_executing)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", final_command)

	cmd.Stdout = stdout_file
	cmd.Stderr = stderr_file

	// Set client user with less permissions
	client_user, err := user.Lookup("client")
	if err != nil {
		log.Fatal("Error looking up client user:", err.Error())
	}
	client_user_uid, _ := strconv.Atoi(client_user.Uid)
	client_user_gid, _ := strconv.Atoi(client_user.Uid)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(client_user_uid),
			Gid: uint32(client_user_gid),
		},
	}

	cmd.Dir = root_dir

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func executeTask(t task.Task) error {
	log.Println("Executing task.")

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
		log.Println("Unable to clone repo:", err.Error())
		error_file, err := utils.CreateErrorFile(t.TaskId.String(), "Error cloning repo: "+err.Error())
		if err != nil {
			log.Println("Unable to create error file for task '", t.TaskId.String(), "':", err.Error())
			SetTaskStatusToUnexpectedError(nats_server, t)
			return
		}

		err = CreateTaskResult(nats_server, result.Result{TaskId: t.TaskId, Files: []string{error_file}})
		if err != nil {
			log.Println("Unable to create result for task '", t.TaskId.String(), "':", err.Error())
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
		log.Println("Error when executing :", err.Error())
		error_file, err := utils.CreateErrorFile(t.TaskId.String(), "Error executing task: "+err.Error())
		if err != nil {
			log.Println("Unable to create error file for task '", t.TaskId.String(), "':", err.Error())
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
			log.Println("Unable to post result for task '", t.TaskId.String(), "':", err.Error())
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
		log.Println("Unable to post result for task '", t.TaskId.String(), "':", err.Error())
		SetTaskStatusToUnexpectedError(nats_server, t)
		return
	}

	SetTaskStatusToFinished(nats_server, t)
}

func main() {

	nats_server, err := nats.Connect(os.Getenv("NATS_SERVER_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	wg := sync.WaitGroup{}
	go utils.WaitForSigkill(&wg)

	waitForTasks(nats_server, &wg)
}
