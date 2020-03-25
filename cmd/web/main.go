package main

import (
	"andreyladmj/analytics/pkg/grpcapi"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	//"github.com/golang-migrate/migrate/v4"
	//mysql_migrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	grpcClient *grpcapi.GRPCClient
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4001", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LUTC|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ltime|log.Lshortfile)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Could not connect to the server: %v", err)
	}
	defer cc.Close()

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		session:    makeSession(),
		grpcClient: grpcapi.New(cc),
	}

	//tslconfig := &tls.Config{
	//	PreferServerCipherSuites: true,
	//	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256}, // these two have assembly implementations
	//}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		//TLSConfig:    tslconfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}


	infoLog.Printf("Starting server on %s", *addr)
	//err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func makeSession() *sessions.Session {
	secret := StringWithCharset(32, "")
	session := sessions.New([]byte(secret))
	session.Lifetime = 10 * time.Minute
	//session.Secure = true
	session.SameSite = http.SameSiteStrictMode
	return session
}


