package postgres

import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const insertFirstUser = `
INSERT INTO users (
	user_id,
	name,
	batch,
	email,
	password,
	created_by,
	updated_by
) VALUES (
	'a67dc254-7410-4394-8b53-eb148eab0477',
	$1,
	$2,
	$3,
	$4,
	'a67dc254-7410-4394-8b53-eb148eab0477',
	'a67dc254-7410-4394-8b53-eb148eab0477'
) RETURNING
	user_id;
`

const insertFirstCurrency = `
INSERT INTO currency (
	name,
	exchange_rate,
	created_by,
	updated_by
) VALUES (
	'BDT',
	1,
	'a67dc254-7410-4394-8b53-eb148eab0477',
	'a67dc254-7410-4394-8b53-eb148eab0477'
) RETURNING
	id;
`

func InsertDemoData() {
	configPath := flag.String("config", "env/config.yaml", "config file")

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile(*configPath)
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := Open(config)
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}
	defer func() { _ = db.Close() }()
	passByte, err := bcrypt.GenerateFromPassword([]byte(config.GetString("admin.password")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(insertFirstUser, config.GetString("admin.name"), config.GetString("admin.batch"), config.GetString("admin.email"), passByte)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	_, err = db.Exec(insertFirstCurrency)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer db.Close()
}
