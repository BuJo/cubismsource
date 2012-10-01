package main

import "fmt"
import "log"
import "time"
import "net/http"

//import "code.google.com/p/gosqlite/sqlite"

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func handleHttpRoot(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Printf("got req: %#v\n", r)

	return nil
}
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		http.Error(w, e.Message, e.Code)
	}
}

func main() {
	/*
		filename := "jvm-ram.db"
		version := sqlite.Version()
		fmt.Printf("hello, world %s\n", version)

		dbconn, err := sqlite.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer dbconn.Close()

		stmt, stmterr := dbconn.Prepare("select date, free from ram where site = ?1 and date between ?2 and ?3")
		if stmterr != nil {
			fmt.Print(err)
		}
		defer stmt.Finalize()

		err = stmt.Exec("BSL", "2012-08-22 1:00", "2012-08-24 23:00")
		if err != nil {
			fmt.Print(err)
		}

		var date string
		var free int

		for stmt.Next() {
			err = stmt.Scan(&date, &free)
			fmt.Printf("D: %s f: %s\n", date, free)
		}
	*/
	s := &http.Server{
		Addr:           ":8080",
		Handler:        appHandler(handleHttpRoot),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
