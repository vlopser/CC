package main

import (
	. "cc/pkg/lib/QueueManager"
	. "cc/pkg/lib/StoreManager"
	. "cc/pkg/lib/TaskManager"
	"cc/pkg/models/result"
	"cc/pkg/models/task"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var wg sync.WaitGroup

func cloneRepo(repo_url string, root_dir string) {

	// Verificar si el directorio existe
	if _, err := os.Stat(root_dir + task.REPO_DIR); err == nil {
		// Si el directorio existe, lo eliminamos
		err := os.RemoveAll(root_dir + task.REPO_DIR)
		if err != nil {
			fmt.Println("Error al eliminar el directorio existente:", err)
			return
		}
	}

	// Clonar el repositorio
	_, err := git.PlainClone(root_dir+task.REPO_DIR, false, &git.CloneOptions{
		URL:      repo_url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("Error al clonar el repositorio:", err)
		return
	}

}

func executeTask(root_dir string) {

	// Abrir el archivo en modo escritura, crear si no existe, agregar al final si existe
	file, err := os.OpenFile(root_dir+task.RESULT_DIR+task.OUTPUT_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	file_errors, err := os.OpenFile(root_dir+task.RESULT_DIR+task.ERRORS_FILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file_errors.Close()

	cmd := exec.Command("go", "mod", "tidy")

	cmd.Dir = root_dir + task.REPO_DIR
	cmd.Stderr = file_errors

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error when executing 'go mod tidy':", err)
		return
	}

	cmd = exec.Command("go", "run", root_dir+task.REPO_DIR+"/main.go")

	cmd.Stdout = file
	cmd.Stderr = file_errors

	// Ejecutar el comando
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error when executing `go run "+root_dir+task.REPO_DIR+"/main.go':", err)
		return
	}
}

func createDirectories(root_dir string) {

	err := os.Mkdir(root_dir, 0755)
	if err != nil {
		fmt.Println("Error al crear el directorio:", err)
		return
	}

	task_result_dir := root_dir + task.RESULT_DIR
	err = os.Mkdir(task_result_dir, 0755)
	if err != nil {
		fmt.Println("Error al crear el directorio:", err)
		return
	}

}

func manageTask(task task.Task) {

	task_dir := task.TaskId.String()
	createDirectories(task_dir)

	cloneRepo(task.Input, task_dir)

	executeTask(task_dir)

}

func cleanDirectory(root_dir string) {

	err := os.RemoveAll(root_dir)
	if err != nil {
		log.Println(err)
	}

}

func handleRequest(task task.Task, nats_server *nats.Conn) {
	log.Println("Request received!")

	SetTaskStatusToExecuting(nats_server, task.TaskId.String())

	init := time.Now()
	manageTask(task)
	end := time.Now()

	PostResult(
		nats_server,
		result.Result{
			TaskId:      task.TaskId,
			TimeElapsed: end.Sub(init),
		},
	)

	SetTaskStatusToFinished(nats_server, task.TaskId.String())

	cleanDirectory(task.TaskId.String()) //Clean all tmp directories created for the task

	//CODIGO DEL FRONTEND
	res, err := GetResult(nats_server, task.TaskId.String())
	if err != nil {
		log.Println("Error when storing the result of", task.TaskId, ":", err)
		return
	}
	println("Solucion en ", res.Output)
}

func waitForSigkill() {

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	<-c

	wg.Done()

	os.Exit(0)
}

func waitForTasks(nats_server *nats.Conn) {
	// GetTasks(nats_server, handleFunction)

	wg = sync.WaitGroup{}
	wg.Add(1)

	log.Println("Waiting for request. (Presiona Ctrl+C para salir)")

	wg.Wait()                   // Esperar a que se reciba una señal de interrupción
	time.Sleep(1 * time.Second) // Permitir tiempo para desuscribirse antes de salir
}

func main() {

	nats_server, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nats_server.Close()

	go waitForSigkill()

	GetTasks(nats_server, handleRequest)

	task := task.Task{
		TaskId: uuid.New(),
		Input:  "https://github.com/go-training/helloworld.git",
		Status: task.PENDING,
	}

	EnqueueTask(task, nats_server)

	waitForTasks(nats_server)
}
