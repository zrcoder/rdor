package debug

import (
	"fmt"
	"log"
	"os"
)

func Log(format string, v ...any) {
	if os.Getenv("RDD") == "1" {
		msg := fmt.Sprintf(format+"\n", v...)
		f, err := os.OpenFile("_rdor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()
		_, err = f.WriteString(msg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
