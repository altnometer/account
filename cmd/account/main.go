package main

import (
	"flag"

	"github.com/altnometer/account/service"
)

func main() {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	service.StartWebServer(*port)
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
