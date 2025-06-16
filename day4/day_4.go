package main

import (
	"fmt"
)

type task struct {
	Id     int
	Desc   string
	Status bool
}

func IdGen() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}
func AddTask(tasks *[]task, getNextId func() int, Description string) {
	newTask := task{
		Id:     getNextId(),
		Desc:   Description,
		Status: false,
	}
	*tasks = append(*tasks, newTask)
	fmt.Printf("Added task %v - %s\n", newTask.Id, newTask.Desc)
}
func TaskList(tasks *[]task) {
	fmt.Println("\nPending Tasks:")
	for _, task := range *tasks {
		if task.Status == false {
			fmt.Printf("%d :%s \n", task.Id, task.Desc)
		}
	}
}

func TaskComplete(tasks *[]task, id int) {
	for i, task := range *tasks {
		if task.Id == id {
			(*tasks)[i].Status = true
			fmt.Printf("n Task %d marked as complete\n  ", id)
			return
		}
	}
	fmt.Printf("Task %d not found\n  ", id)
}
func main() {
	//this is main
	var tasks []task
	nxtId := IdGen()
	AddTask(&tasks, nxtId, "Buy Medicines")
	AddTask(&tasks, nxtId, "GO to Gym")
	TaskList(&tasks)
	TaskComplete(&tasks, 1)
	TaskList(&tasks)
	//end  ...
}
