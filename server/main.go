package main

import (
	"fmt"
	"log"
	"net"

	usergrpc "github.com/khdip/help-save-a-life/proto/users"
	usercore "github.com/khdip/help-save-a-life/server/core/users"
	usersvc "github.com/khdip/help-save-a-life/server/services/users"

	collgrpc "github.com/khdip/help-save-a-life/proto/collection"
	collcore "github.com/khdip/help-save-a-life/server/core/collection"
	collsvc "github.com/khdip/help-save-a-life/server/services/collection"

	commgrpc "github.com/khdip/help-save-a-life/proto/comments"
	commcore "github.com/khdip/help-save-a-life/server/core/comments"
	commsvc "github.com/khdip/help-save-a-life/server/services/comments"

	dregrpc "github.com/khdip/help-save-a-life/proto/dailyReport"
	drecore "github.com/khdip/help-save-a-life/server/core/dailyReport"
	dresvc "github.com/khdip/help-save-a-life/server/services/dailyReport"

	currgrpc "github.com/khdip/help-save-a-life/proto/currency"
	currcore "github.com/khdip/help-save-a-life/server/core/currency"
	currsvc "github.com/khdip/help-save-a-life/server/services/currency"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	acctcore "github.com/khdip/help-save-a-life/server/core/accountType"
	acctsvc "github.com/khdip/help-save-a-life/server/services/accountType"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	acntcore "github.com/khdip/help-save-a-life/server/core/accounts"
	acntsvc "github.com/khdip/help-save-a-life/server/services/accounts"

	linkgrpc "github.com/khdip/help-save-a-life/proto/links"
	linkcore "github.com/khdip/help-save-a-life/server/core/links"
	linksvc "github.com/khdip/help-save-a-life/server/services/links"

	docsgrpc "github.com/khdip/help-save-a-life/proto/medDocs"
	docscore "github.com/khdip/help-save-a-life/server/core/medDocs"
	docssvc "github.com/khdip/help-save-a-life/server/services/medDocs"

	settgrpc "github.com/khdip/help-save-a-life/proto/settings"
	settcore "github.com/khdip/help-save-a-life/server/core/settings"
	settsvc "github.com/khdip/help-save-a-life/server/services/settings"

	"strconv"
	"strings"

	"github.com/khdip/help-save-a-life/server/storage/postgres"

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

	linkC := linkcore.New(store)
	linkS := linksvc.New(linkC)
	linkgrpc.RegisterLinkServiceServer(grpcServer, linkS)

	docsC := docscore.New(store)
	docsS := docssvc.New(docsC)
	docsgrpc.RegisterMedDocsServiceServer(grpcServer, docsS)

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
