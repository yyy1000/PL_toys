package WorkScheduler

// version 1
/*
func Schedule(servers []string, numTask int, call func(srv string, task int)) {
	idle := make(chan string, len(servers))
	// use a buffered channel as a concurrent blocking queue
	for _, srv := range servers {
		idle <- srv
	}

	for task := 0; task < numTask; task++ {
		srv := <-idle
		go func(taskNum int) {
			call(srv, taskNum)
			idle <- srv
		}(task)
	}

	for i := 0; i < len(servers); i++ {
		<-idle
	}
}
*/
func Schedule(servers []string, numTask int, call func(srv string, task int) bool) {
	work := make(chan int, numTask)
	done := make(chan bool)

	runTasks := func(srv string) {
		for task := range work {
			if call(srv, task) {
				done <- true
			} else {
				work <- task
			}

		}
	}

	go func() {
		for _, srv := range servers {
			go runTasks(srv)
		}
	}()

	for task := 0; task < numTask; task++ {
		work <- task
	}

	for i := 0; i < numTask; i++ {
		<-done
	}
	close(work)
}
