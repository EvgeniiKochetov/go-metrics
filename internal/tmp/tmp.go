package tmp

import (
	"fmt"

	"os"

	"time"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/logger"

	"github.com/EvgeniiKochetov/go-metrics-tpl/internal/handler"
)

func SaveInFile(filename string, interval string) {

	freq, err := time.ParseDuration(interval)
	if err != nil {
		logger.Log.Error("can't convert interval")
		return
	}

	f, err := OpenFile(filename)
	if err != nil {
		fmt.Println("can't open file")
		logger.Log.Error("can't open file")
		return
	}
	defer f.Close()

	for {
		fmt.Println("saving file")
		err = handler.Memory.SaveStorage(f.Name())
		if err != nil {
			fmt.Println("error save metrics")
			logger.Log.Error("error save metrics")
			return
		}
		time.Sleep(time.Duration(freq))
	}

}

func OpenFile(filename string) (os.File, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return *file, err
	}

	return *file, nil
}
