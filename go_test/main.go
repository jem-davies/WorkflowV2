package main

import (
	"fmt"
	"time"
)

// TODO - use a lock rather than the sleep on line 93
// TODO - alter from using array of size 6 to slices

func main() {
	// workflow_matrix from yaml config
	var workflow_matrix [6][6]int = [6][6]int{
		{0, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0},
	}
	// branch_list from yaml config
	var branch_list [6]string = [6]string{"A", "B", "C", "D", "E", "F"}

	var branch_seconds [6]int = [6]int{3, 40, 3, 3, 3, 3}

	// NS - Not Started
	// F - Finished
	// R - Running
	var branch_status [6]string = [6]string{"NS", "NS", "NS", "NS", "NS", "NS"}

	type branchStatusMessage struct {
		index  int
		status string
	}
	// transpose the matrix for convenience
	transposed_matrix := transposeMatrix(workflow_matrix)

	print("\n")
	for i, idx := range transposed_matrix {
		fmt.Print(i, idx, "\n")
	}

	statusChannel := make(chan branchStatusMessage)
	updateChannel := make(chan int)

	go func() {
		for {
			index := <-updateChannel
			transposed_matrix = zeroOutColumn(transposed_matrix, index)
		}
	}()

	go func() {
		for {
			statusUpdate := <-statusChannel
			branch_status[statusUpdate.index] = statusUpdate.status
			//fmt.Print(branch_status, '\n')
		}
	}()

	for {
		for i, row := range transposed_matrix {
			if allZeros(row) && branch_status[i] == "NS" {
				branch_name := branch_list[i]
				branch_second := branch_seconds[i]
				message := branchStatusMessage{index: i, status: "R"}
				statusChannel <- message

				go func(branch_name string, index int, branch_second int) {

					fmt.Printf("Node: id=%s, started\n", branch_name)

					time.Sleep(time.Second * time.Duration(branch_second)) // simulate the running of a go routine

					fmt.Printf("Node: id=%s, finished\n", branch_name)
					message = branchStatusMessage{index: index, status: "F"}
					statusChannel <- message
					updateChannel <- index
				}(branch_name, i, branch_second)
			}
			// print("\n")
			// for i, idx := range transposed_matrix {
			// 	fmt.Print(i, idx, "\n")
			// }
		}
		//time.Sleep(time.Second * 10)
	}
}

func transposeMatrix(matrix [6][6]int) [6][6]int {
	var transposed [6][6]int

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			transposed[i][j] = matrix[j][i]
		}
	}

	return transposed
}

func allZeros(arr [6]int) bool {
	for _, v := range arr {
		if v != 0 {
			return false
		}
	}
	return true
}

func zeroOutColumn(matrix [6][6]int, columnIndex int) [6][6]int {
	for i := 0; i < 6; i++ {
		matrix[i][columnIndex] = 0
	}
	return matrix
}

/*

Node: id=A, started
Node: id=A, finished
Node: id=C, started
Node: id=B, started
Node: id=C, finished
Node: id=E, started
Node: id=E, finished
Node: id=F, started
Node: id=F, finished
Node: id=B, finished
Node: id=D, started
Node: id=D, finished

*/
