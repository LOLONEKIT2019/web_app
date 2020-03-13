package web_app

import (
	. "./handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/status", http.HandlerFunc(RootHandler))

	r.Handle("/token", http.HandlerFunc(CreateToken))
	r.Handle("/users", JwtAuthentication(http.HandlerFunc(UsersRouter)))
	r.Handle("/users/{id}", JwtAuthentication(http.HandlerFunc(UsersRouter)))
	r.Handle("/contacts/{id}", JwtAuthentication(http.HandlerFunc(ContactRouter)))

	log.Println("start")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

/*//endpoints
http.HandleFunc("/users", UsersRouter)
http.HandleFunc("/users/", UsersRouter)

//API status
http.HandleFunc("/status", RootHandler)


s := &http.Server{
	Addr:           ":8000",
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

log.Println("web-app start seccesfuly ! on port ", s.Addr)
err := s.ListenAndServe()
if err != nil {
	log.Println(err)
	os.Exit(1)
}

*/
