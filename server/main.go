package main

import (
	"fmt"
	"log"
	"net"

	usergrpc "help-save-a-life/proto/users"
	usercore "help-save-a-life/server/core/users"
	usersvc "help-save-a-life/server/services/users"

	collgrpc "help-save-a-life/proto/collection"
	collcore "help-save-a-life/server/core/collection"
	collsvc "help-save-a-life/server/services/collection"

	commgrpc "help-save-a-life/proto/comments"
	commcore "help-save-a-life/server/core/comments"
	commsvc "help-save-a-life/server/services/comments"

	dregrpc "help-save-a-life/proto/dailyReport"
	drecore "help-save-a-life/server/core/dailyReport"
	dresvc "help-save-a-life/server/services/dailyReport"

	currgrpc "help-save-a-life/proto/currency"
	currcore "help-save-a-life/server/core/currency"
	currsvc "help-save-a-life/server/services/currency"

	acctgrpc "help-save-a-life/proto/accountType"
	acctcore "help-save-a-life/server/core/accountType"
	acctsvc "help-save-a-life/server/services/accountType"

	acntgrpc "help-save-a-life/proto/accounts"
	acntcore "help-save-a-life/server/core/accounts"
	acntsvc "help-save-a-life/server/services/accounts"

	settgrpc "help-save-a-life/proto/settings"
	settcore "help-save-a-life/server/core/settings"
	settsvc "help-save-a-life/server/services/settings"

	"help-save-a-life/server/storage/postgres"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("server/env/config.yaml")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()
	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	userC := usercore.New(store)
	userS := usersvc.New(userC)
	usergrpc.RegisterUserServiceServer(grpcServer, userS)

	collC := collcore.New(store)
	collS := collsvc.New(collC)
	collgrpc.RegisterCollectionServiceServer(grpcServer, collS)

	commC := commcore.New(store)
	commS := commsvc.New(commC)
	commgrpc.RegisterCommentServiceServer(grpcServer, commS)

	dreC := drecore.New(store)
	dreS := dresvc.New(dreC)
	dregrpc.RegisterDailyReportServiceServer(grpcServer, dreS)

	currC := currcore.New(store)
	currS := currsvc.New(currC)
	currgrpc.RegisterCurrencyServiceServer(grpcServer, currS)

	acctC := acctcore.New(store)
	acctS := acctsvc.New(acctC)
	acctgrpc.RegisterAccountTypeServiceServer(grpcServer, acctS)

	acntC := acntcore.New(store)
	acntS := acntsvc.New(acntC)
	acntgrpc.RegisterAccountsServiceServer(grpcServer, acntS)

	settC := settcore.New(store)
	settS := settsvc.New(settC)
	settgrpc.RegisterSettingsServiceServer(grpcServer, settS)

	host, port := config.GetString("server.host"), config.GetString("server.port")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	log.Printf("Server is starting on: %s:%s", host, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslmode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransactionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
