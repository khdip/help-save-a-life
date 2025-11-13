# Help Save A Life

A simple and ready to deploy web application to collect and manage funds for an invidual who needs donations (i.e.: Medical purposes). It's a Go-based web application built using microservice architecture and gRPC.

## Features
- You can create multiple users to manage the website.
- You can change the title logo, banner title, banner description, banner image, about section from settings.
- You can choose which sections will show in homepage.
- You can add, edit, delete each transaction from the collection page.
- If adding every transaction is not feasible for you, you can add the total amount of each day at the end of the day. You can also edit and delete an entry if needed.
- You can add multiple currencies with an exchange rate if the funding is coming from different countries.
- The total collection amount is calculated from the total of collection or from the total of daily report. You can choose how it will be calculated from settings. You may also use a custom value if you are not keeping track of each transaction or total of each day.
- You can upload medical documents as pdf or jpg or png to show in homepage inside a slider. You can upload multiple medical documents at once.
- You can put important links such as facebook page link or facebook event link or news link from important links menu.
- Input fields are validated from backend.
- People can fillup the contact form and you can access the comments from comments menu in the admin center.
- You can add multiple account types for example: Bank, MFS, Bkash, Nagad etc.
- You may add one or more accounts for each account type. These accounts will show in the donate section in the homepage.
- You may choose to show each transaction in homepage. The last four digits of the account number are hidden.

## Tools Used
- The Go Programming Language (Version 1.25.0)
- PostgreSQL (Version 17.4)
- Protocol Buffers (Version 32.0)
- gRPC

## Setup
1. Download and install the Go programming language. [Download](https://go.dev/)
2. Download and install PostgreSQL. [Download](https://www.postgresql.org/) 
3. Install Protocol Buffer Compiler. [Install](https://protobuf.dev/installation/)
4. Install the Protocol compiler plugin for Go. [Install](https://grpc.io/docs/languages/go/quickstart/)
5. Clone this repository.
```
git clone https://github.com/khdip/help-save-a-life.git
```
6. Create a database. Run below command in psql shell.
```
CREATE DATABASE save_life;
```
> [!NOTE]  
> If you want to create a database with a different name, make sure to change the database name in env file.
7. Download all the go package dependencies.
```
go mod tidy
```
8. Modify the config.yaml files inside the the folders server > env and cms > env according to your preference.
9. Navigate to the server directory and run the migration. Assuming you are in the root folder of the project:
```
cd server
go run migrations/migrate.go up
```
9. Run the seeder to seed initial data into the database.
```
go run seeder/seeder.go
```
10. Navigate back to the root folder of the project and run the grpc server.
```
cd ..
go run server/main.go
```
11. Open another terminal window and run the grpc client server.
```
go run cms/main.go
```
12. Login to the admin dashboard using the admin credentials in the env file you mentioned by going to the route /login. Then go to the Settings option from menu and change accordingly.

## Routes
- Homepage:                 <mark>/</mark>
- Login:                    <mark>/login</mark>
- Logout:                   <mark>/logout</mark>
- Dashboard:                <mark>/dashboard</mark>
- Users list:               <mark>/users</mark>
- Collection list:          <mark>/collection</mark>
- Comments list:            <mark>/comments</mark>
- Daily report list:        <mark>/daily_report</mark>
- Currency List:            <mark>/currencies</mark>
- Account type list:        <mark>/account_types</mark>
- Accounts list:            <mark>/accounts</mark>
- Link list:                <mark>/links</mark>
- Medical documents list:   <mark>/med_docs</mark>
- Settings:                 <mark>/settings</mark>

## Advance Customizations
- To create a new db migration file navigate to the server folder and run the below command.
```
go run migrations/migrate.go create create_<your_table_name>_table sql
```
- To generate the pb.go and grpc.pb.go files from a .proto file:
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./proto/<your_folder_name>/<your_proto_file_name>.proto 
```

## License
This project is licensed under the MIT License—see the LICENSE file for details.