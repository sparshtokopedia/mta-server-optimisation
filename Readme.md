mta-server-optimization

To set up the server, follow these steps:

Open the app.go file in your project.

Initialize the Go module by running the following command in your terminal or command prompt:
go mod init mta-server-optimiser
This command creates a go.mod file that tracks the project's dependencies.

Run go mod tidy to ensure that the go.mod file reflects the correct and updated dependencies based on your code.

Use go mod vendor to create a vendor directory that contains all the dependencies required by your project. This step ensures that you have a local copy of the dependencies.

Build the project by running:
go build
This command compiles the Go code and generates an executable binary file.

Finally, run the server by executing:
go run app.go
This command starts the server and runs your application.

Make sure you are in the correct directory when running these commands. Adjust the commands accordingly if your project structure or file names are different.

with endpoint localhost:8080/mta-hosting-optimize
And all the unit and integration test cases are written in app_test.go 
And for changing the value of "X" , you can change it in getInstanceName function or you can set the value in .env file

Output will be this 
[ "mta-prod-1", "mta-prod-3" ]


