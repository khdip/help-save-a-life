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
	'05cf-71sl-cas7-oqtf-fsds',
	'Dip',
	15,
	'khdip.ku@gmail.com',
	$1,
	'05cf-71sl-cas7-oqtf-fsds',
	'05cf-71sl-cas7-oqtf-fsds'
) RETURNING
	user_id;
`

func CreateFirstUser() {
	passByte, err := bcrypt.GenerateFromPassword([]byte("arthaBainchod"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
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
	_, err = db.Exec(insertFirstUser, passByte)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer db.Close()
}
