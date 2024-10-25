package render

func DoOn(queue chan func(), work func()) {
	done := make(chan bool, 1)

	queue <- func() {
		work()
		done <- true
	}

	<-done
}
