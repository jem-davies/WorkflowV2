package main

import (
	"fmt"
	"time"
)

func main() {
	// workflow_matrix from yaml config
	var workflow_matrix [][]int = [][]int{
		{0, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0},
	}

	// branch_list from yaml config
	var branch_list []string = []string{"A", "B", "C", "D", "E", "F"}

	// for simulating the running of a go routine with different lengths of time
	var branch_seconds []int = []int{3, 15, 3, 3, 3, 3}

	// NS - Not Started
	// F - Finished
	// R - Running
	var branch_status []string = []string{"NS", "NS", "NS", "NS", "NS", "NS"}

	for !isAllFinished(branch_status) {
		for i, status := range branch_status {
			if isColumnAllZeros(workflow_matrix, i) && status == "NS" {
				branch_status[i] = "R"

				go func(branch_name string, i int, branch_second int) {
					fmt.Printf("Node: id=%s, started\n", branch_name)
					time.Sleep(time.Second * time.Duration(branch_second))
					fmt.Printf("Node: id=%s, finished\n", branch_name)
					workflow_matrix = zeroOutRow(workflow_matrix, i)
					branch_status[i] = "F"
				}(branch_list[i], i, branch_seconds[i])
			}
		}
	}
	fmt.Println("DAG Workflow execution finished!")
}

func isColumnAllZeros(matrix [][]int, columnIdx int) bool {
	for i := 0; i < len(matrix); i++ {
		if matrix[i][columnIdx] != 0 {
			return false
		}
	}
	return true
}

func zeroOutRow(matrix [][]int, rowIdx int) [][]int {
	for i := 0; i < len(matrix); i++ {
		matrix[rowIdx][i] = 0
	}
	return matrix
}

func isAllFinished(branch_status []string) bool {
	for i := 0; i < len(branch_status); i++ {
		if branch_status[i] != "F" {
			return false
		}
	}
	return true
}
