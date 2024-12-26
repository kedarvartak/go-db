

# GoDatabase  

A lightweight, performant, and educational database implemented in Go, showcasing core database concepts and functionalities. This project is designed for developers and enthusiasts who want to understand the inner workings of a database system.

## Features  

### Core Functionality  
- **Write-Ahead Logging (WAL)**: Ensures data integrity by persisting changes to a log before applying them to the main database.  
- **Page Management**: Efficiently handles data storage and retrieval with a structured page layout system.  
- **Serialization/Deserialization**: Transforms in-memory data structures into a format suitable for storage or transmission and vice versa.  

### Advanced Indexing  
- **B-Tree Indexing**: Implements a balanced tree structure for efficient query lookups and data retrieval.  

### Query Processing  
- **Query Parsing**: Interprets and validates user queries, transforming them into executable operations.  

### Storage Structure  
- **Page Layout**: Organizes data into pages for structured and efficient storage management.  

## Getting Started  

### Prerequisites  
- Go 1.20 or later  
- Basic knowledge of database principles is recommended  

### Installation  
1. Clone the repository:  
   ```bash  
   git clone https://github.com/your-username/GoDatabase.git   
   ```  

## Usage  
After building the project, you can interact with the database using the provided command-line interface (CLI).  

### Example  - 
I have only implemented basic parsing and insertions yet 


2. Insert data:  
   ```sql  
   INSERT INTO users (id, name) VALUES (1, 'Alice');  
   ```  
3. Retrieve data:  
   ```sql  
   SELECT * FROM users;  
   ```
   

## Contributing  

We welcome contributions from the community!  
1. Fork the repository.  
2. Create a new branch for your feature or bug fix.  
3. Submit a pull request describing your changes.  

## License  

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.  

## Contact  

For questions or feedback, please open an issue or contact the repository owner at `kedar.vartak22@vit.edu`.  

--- 

Let me know if you'd like any further adjustments!
