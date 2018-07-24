package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
)

const (
	defaultPort              = "80"
	defaultApiUser           = "http://172.16.100.243:30000/user/users"
	defaultPublicKey         = "./auth/test_jwt_keys/public.pem"
	defaultApiDatachange     = "http://172.16.100.243:30000/queue/datachange"
	defaultFirebaseServerKey = "AAAAebG27Tw:APA91bE0lgckfj3yQrwtewYym8uZWGSRzvN5BOA0OZ1oUQ0QZbz52Q3wxX27dNizJ7xDDSMPWEYYRKA2hahSlc8ooWXz76yoM7Ei5L_9UIftqxM4v3WkUy_c0-91diA3HirhpzCQJB2s"
	defaultApiMail           = "http://172.16.100.243:30000/mailer/notification/send"
	defaultApiOrder          = "http://172.16.100.243:30000/order/"
	defaultUserName          = "admin"
	defaultPassWord          = "p@ssword"
	defaultUrlGraphql        = "http://172.16.100.243:30000/graphql"
)

func main() {
	log.Print("Starting")
	// Get config values
	var (
		port     = flag.String("http.port", defaultPort, "HTTP listen port") //helper.GetEnvString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+*port, "HTTP listen address")
	)
	// Parse predefined config
	flag.Parse()
	// Setup logger
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestamp)
	// Handle request
	http.Handle("/", initHandler(logger))
	// Keep main process running
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	// Terminate when receive terminate signal
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func initHandler(logger kitlog.Logger) http.Handler {
	// Init Database handler
	// mongoHandler := initDatabaseHandler(logger)
	// Init service
	// messageService := message.NewNotifyMessageService(mongoHandler)
	// tempalteService := notifytemplate.NewNotifyTemplateService(mongoHandler)
	// // Init routing
	mux := http.NewServeMux()
	// // api datachange
	// Handle messages
	// mux.Handle("/messages/", message.MakeNotifyMessageHandler(messageService, logger, apiUser, publicKey, apiDataChange, firebaseServerKey, apiMail, apiOrder, username, password, graphql))
	// mux.Handle("/messages/notify_template/", notifytemplate.MakeNotifyMessageHandler(tempalteService, logger, publicKey))

	return mux
	// Handle cors
	// httpHandler := security.HandleCorsAccess(mux)
	// // Handle access log
	// return accesslog.NewApacheLoggingHandler(httpHandler, os.Stderr)
}

// func initDatabaseHandler(logger kitlog.Logger) db.DatabaseHandler {
// 	// Get config values
// 	var (
// 		mongoURL      = helper.GetEnvString("MONGO_HOST", db.DbHost)
// 		mongoPort     = helper.GetEnvString("MONGO_PORT", strconv.Itoa(db.DbPort))
// 		mongoUser     = helper.GetEnvString("MONGO_USER", db.DbUser)
// 		mongoPass     = helper.GetEnvString("MONGO_PASS", db.DbPass)
// 		mongoDataBase = helper.GetEnvString("MONGO_DB", db.CollectionName)
// 		authDatabase  = helper.GetEnvString("MONGO_AUTH_DB", db.AuthDb)
// 	)

// 	// Init db handler
// 	mongoPortNumber, ok := strconv.Atoi(mongoPort)
// 	if ok != nil {
// 		logger.Log("[App.error]: Cannot get db port")
// 		os.Exit(1)
// 	}
// 	return db.NewMongoHandler(mongoURL, mongoPortNumber, mongoDataBase, authDatabase, mongoUser, mongoPass)
// }
