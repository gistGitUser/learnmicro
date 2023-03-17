package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"p02/handlers"
)

/*
в целом нормальная практика для каждой структуры
использовать отдельный логер, т.к. эти структуры
не всегда могут быть связаны
 */


func main() {

	//стоит почитать документацию по стандартным пакетам, т.к.
	//это явно улучшит понимание языка
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	// create the handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	// create a new server
	s := http.Server{
		Addr:         ":9090",      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		//эти таймауты очень важны, т.к.
		//по сути определяют то для чего
		//предназначен сервер, т.е.
		//большие таймауты говорят о том, что мы
		//хотим загружать отдавать большие файлы
		//и наоборот
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		//если тут установлен 0, то соединение
		//бесконечное
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	//для grac shut нужно запустить сервер в горутине
	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	//ctrl+c or kill -9
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	/*
	эта штука важна, т.к. если сервер выполняет какую-либо работу
	то без этого работа прервется, а так будет возможность
	корректно завершить работу
	 */
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//мягко выключаев сервер с попыткой дождаться окончания работы
	//в течении 30 секунд
	s.Shutdown(ctx)

}
