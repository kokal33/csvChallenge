Journey and Learnings

I set out to create a mapping of the cells to get the formulas working. It took a couple of iterations to get it right. The first thing I learned was how to work with struct and mapping in Go. Here is a brief outline of my journey:

    My plan:
        Make a mapping of the cells, so the formulas work.
        Make a print function to print the mapped contents as a table.
        Execute the functions.

    First Iteration:

    I started by creating a struct to hold the cell information (row, col, and val). I then mapped these cells to their respective values. This approach was a bit complicated, so I decided to simplify it.

    Second Iteration:

    In this step, I decided to ditch the struct and only use mapping. I split the data into rows and columns and created a map to store each cell's data. This approach was more in line with what I needed.

    Improving the Print Function:

    To create a table representation of the data, I added a print function. This function prints the column headers dynamically and the contents of each cell. I also added a feature to print 3 lines of space when a header is encountered for better readability. However, I needed to make some adjustments:
        Added a row number for empty lines.
        Prevented the function from printing 3 empty spaces before the first headers.

    Mapping the Cells:

    I realized I had not mapped the three empty spaces in between headers, so I updated my function to increment the row in the mapping three times whenever a header is encountered.

    Reading from a File:

    Until this point, I had been using a hardcoded string as input. I then wrote a readFile() function to read CSV data from a file.

    Processing Formulas:

    I then moved onto processing the formulas in the cells. I started by processing the "^^" operation, which copies the formula from the cell above in the same column, using the ProcessDoubleCaret() function.

    Next, I needed to weed out the standalone formulas and simplify them. Once I had the clean standalone formulas, I could solve them. I attempted to use a function mapper to handle the various operations dynamically. However, this proved to be quite challenging, so I decided to use a switch-case structure, ensuring the parameters were converted into the right type before calling the respective function.

    After the standalone formulas were solved, I mapped them to the global Mapping so that the evaluated ones could take the values from them. I also had to map the "E^" and "E^v" values to their respective cells before solving.

    All I had left in the cell formulas after this were the basic expressions (*,+,-,/). Although I originally planned to implement the evaluations of these simple expressions myself, I decided to use the govaluate package for this purpose, to make the process more straightforward.

    Some snips from my previos iteration failures :) 

    var funcMap = map[string]interface{}{
    	"incFrom": incFrom,
    	"text":    text,
    	"concat":  concat,
    	"split":   split,
    	"spread":  spread,
    	"sum":     sum,
    	"bte":     bte,
    }

    func GetFunction(name string) (interface{}, bool) {
    	function, found := funcMap[name]
    	return function, found
    }
    

    Final Steps:

    Once all the operations were completed, I mapped everything to the global mapping and called the print function to display the final output.

Files and Functions

The main Go files and their respective functions:

    hello.go:
    This file contains the main() function, which reads the data, maps the cells, and calls the print function. It also includes the readFile() function for reading the CSV data from a file.

    helper.go:
    This file contains all the helper functions, cleaning, formatting, mapping and all operations

    functions.go:
    This file has all the main functions like solving, expression understanding and all the excel functions written in Go.

    Conclusion

In conclusion, this project helped me significantly in understanding Go language's nuances. This README contains a brief overview of the work I've done. Please refer to the comments in the source code for a more detailed explanation.