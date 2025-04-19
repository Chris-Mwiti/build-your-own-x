package main

import (
	"fmt"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/manager"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/node"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/scheduler"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/google/uuid"
)


func main(){
	fmt.Println("Upcoming orchestra project")

	task := task.Task{
		ID: uuid.New(),
		Name: "Task-one",
		State: task.Pending,
		Image: "Image One",
		Memory: 1024,
		Disk: 1,
	}
}