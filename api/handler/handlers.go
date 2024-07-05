package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webpaygo/api/controller"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Init() {

	router := mux.NewRouter()

	//view

	router.HandleFunc("/", controller.InitTransaction).Methods("GET")
	router.HandleFunc("/commit", controller.VerifTransaction).Methods("POST")
	// se le pasan los datos de la transacción: order_id, session_id, monto, url_retorno
	router.HandleFunc("/save-transaction", controller.SaveTransaction).Methods("POST")

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	handler := cors.AllowAll().Handler(router)

	fmt.Println("server en escucha: " + PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, handler))

}
