Installation Steps

    Clone the repository

    Open your terminal and use the git command to clone the repository:


git clone https://github.com/kokal33/csvChallenge

cd csvChallenge

Install all dependencies

Use the go mod tidy command to download and install all required dependencies:

    go mod tidy

    This will install all the dependencies mentioned in the go.mod file, including the govaluate package.

Running the program

To run the program, you can use the go run command followed by the name of the .go file that contains the main function. For instance, if the main function is in a file named main.go, you would run:

bash

go run main.go