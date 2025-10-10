package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	//fmt.Printf("Value: %#v, type: %T\n", err, err) // Value: (*fs.PathError)(nil), type: *fs.PathError
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
