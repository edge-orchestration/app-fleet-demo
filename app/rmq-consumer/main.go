package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	RabbitMq_Host string = "localhost"
	RabbitMq_Port string = "5672"
	RabbitMq_User string = "guest"
	RabbitMq_Pass string = "guest"
)

type templateData struct {
	Messages *[]string
}

func main() {
	flag.StringVar(&RabbitMq_Host, "rabbitmq-host", LookupEnvOrString("RABBITMQ_HOST", RabbitMq_Host), "RabbitMQ Host")
	flag.StringVar(&RabbitMq_Port, "rabbitmq-port", LookupEnvOrString("RABBITMQ_PORT", RabbitMq_Port), "RabbitMQ port")
	flag.StringVar(&RabbitMq_User, "rabbitmq-user", LookupEnvOrString("RABBITMQ_USER", RabbitMq_User), "RabbitMQ User")
	flag.StringVar(&RabbitMq_Pass, "rabbitmq-pass", LookupEnvOrString("RABBITMQ_PASS", RabbitMq_Pass), "RabbitMQ Password")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/consumer", home)
	fileServer := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/consumer/assets/", http.StripPrefix("/consumer/assets", fileServer))

	// start web server
	log.Println("Starting RabbitMQ Demo app - Consumer on: 4002")
	log.Printf("RabbitMQ Instance: %s:%s\n", RabbitMq_Host, RabbitMq_Port)
	log.Printf("RabbitMQ User: %s\n", RabbitMq_User)

	err := http.ListenAndServe("0.0.0.0:4002", mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error: "+err.Error(), 500)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), 500)
	}
	err = ts.Execute(w, &templateData{Messages: nil})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error: "+err.Error(), 500)
	}
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func RabbitMQInstanceConnectionPath() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s?heartbeat=30&connection_timeout=120", RabbitMq_User, RabbitMq_Pass, RabbitMq_Host, RabbitMq_Port)
}
