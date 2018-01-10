package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/altnometer/account/service"
	"github.com/altnometer/kafkalog"
)

func main() {
	port := flag.String("port", "8080", "server port")
	env := flag.String("env", "dev", "environment, accepted values: dev, prod")
	flag.Parse()
	service.StartWebServer(*port)
	if *env == "prod" {
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			panic("Error running os.Pipe()")
		}
		os.Stdout = w
		lw := kafkalog.NewAsyncProducer("loggerID")
		defer func() {
			w.Close()
			os.Stdout = old
			lw.Close()
		}()
		go func() {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				l := sc.Text()
				lw.Send(l)
			}
		}()

	}
	// reader := bufio.NewReader(os.Stdin)
	// kp := kafka.SyncProducer{}
	// err := kp.InitMySyncProducer()
	// if err != nil {
	// 	panic(err)
	// }
	// kp.SendAccMsg([]byte("hello"))
	// fmt.Print("-> ")
	// text, _ := reader.ReadString('\n')
	// text = strings.Replace(text, "\n", "", -1)
	// args := strings.Split(text, "###")
	// cmd := args[0]

	// switch cmd {
	// case "write":
	// 	if len(args) == 2 {
	// 		msg := args[1]
	// 		// event := NewCreateAccountEvent(accName)
	// 		kp.SendAccMsg([]byte(msg))
	// 	} else {
	// 		fmt.Println("Only specify write###yourmessage")
	// 	}
	// default:
	// 	fmt.Printf("Unknown command %s, only: write, is implementd.\n", cmd)
	// }

	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	err = nil
	// }
}
