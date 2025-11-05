package main

import (
	"fmt"
	"help-save-a-life/cms/handler"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/yookoala/realpath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	acctgrpc "help-save-a-life/proto/accountType"
	acntgrpc "help-save-a-life/proto/accounts"
	collgrpc "help-save-a-life/proto/collection"
	commgrpc "help-save-a-life/proto/comments"
	currgrpc "help-save-a-life/proto/currency"
	dregrpc "help-save-a-life/proto/dailyReport"
	linkgrpc "help-save-a-life/proto/links"
	docsgrpc "help-save-a-life/proto/medDocs"
	settgrpc "help-save-a-life/proto/settings"
	usergrpc "help-save-a-life/proto/users"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config.yaml")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Connection failed", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
	}
	assetPath, err := realpath.Realpath(filepath.Join(wd, "cms/assets"))
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
	}
	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))

	uc := usergrpc.NewUserServiceClient(conn)
	cc := collgrpc.NewCollectionServiceClient(conn)
	cmc := commgrpc.NewCommentServiceClient(conn)
	drc := dregrpc.NewDailyReportServiceClient(conn)
	cuc := currgrpc.NewCurrencyServiceClient(conn)
	sc := settgrpc.NewSettingsServiceClient(conn)
	atc := acctgrpc.NewAccountTypeServiceClient(conn)
	acc := acntgrpc.NewAccountsServiceClient(conn)
	lnk := linkgrpc.NewLinkServiceClient(conn)
	mds := docsgrpc.NewMedDocsServiceClient(conn)
	r := handler.GetHandler(decoder, store, asst, uc, cc, cmc, drc, cuc, sc, atc, acc, lnk, mds)

	host, port := config.GetString("client.host"), config.GetString("client.port")
	log.Println("Server  starting...")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
