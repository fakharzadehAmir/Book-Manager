# Book Manager Mini Project

The Book Manager Mini Project is a web application designed to facilitate book management, allowing users to add, retrieve, update, and delete books. It includes user authentication, user profiles, and supports various book-related operations.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- **Go**: Make sure you have Go (Golang) installed on your system. You can download it from the official [Go website](https://golang.org/dl/).

- **PostgreSQL Database**: This project uses PostgreSQL as the database backend. Make sure you have a PostgreSQL server running and accessible. You'll need the connection details (host, port, username, password, and database name) to set up the database connection in the configuration.

## Project Structure

The project is structured into several packages and files, each serving a specific purpose. Here's an overview of the project structure:

### `db` Package

The `db` package handles database operations using the Gorm library. It defines the structure of tables in the database, including users, books, authors, and table of contents. The package provides functions for CRUD operations on these entities.

### `authenticate` Package

The `authenticate` package handles user authentication and token management. It provides functions for user login and token generation/validation.

### `handlers` Package

The `handlers` package contains HTTP request handler functions responsible for handling various endpoints of the application. Each file in this package focuses on a specific aspect of the application:

- `server.go`: Defines the main `BookManagerServer` struct, which holds instances of database, logger, and authentication components.

- `profile.go`: Handles user profile information retrieval, including authorization token validation and fetching user details.

- `book.go`: Manages book-related operations, such as adding new books, retrieving all books, and handling operations on individual books (get, delete, update).

- `auth.go`: Handles user authentication and registration, including login and signup requests, interacting with authentication package and the database.

### `main.go`

The `main.go` file serves as the entry point of the application. It sets up the router, initializes database connection, authentication, and logger components, and maps URLs to the appropriate handler functions.

## Usage

1. Start the application by running `go run main.go` in your terminal.
2. The application will start, and you can access its functionality through various API endpoints.
3. Use API endpoints for user authentication, book management, and profile retrieval.

## Important Information

- The application uses the Gorilla Mux router for routing and URL mapping.
- It employs the Gorm library to interact with the PostgreSQL database.
- User passwords are encrypted using bcrypt before being stored in the database.
- API endpoints are designed following RESTful principles, utilizing appropriate HTTP methods for different actions.

## Future Enhancements

The Book Manager Mini Project can be extended with additional features such as book searching, filtering, and more advanced user management.

## Packages Used

The project utilizes the following external packages:

- Gorilla Mux (github.com/gorilla/mux): A powerful URL router and dispatcher for matching incoming HTTP requests to their respective handler functions.

- Gorm (gorm.io/gorm): A popular ORM library for Go that simplifies database operations by providing a higher-level abstraction over SQL databases.

- Golang Crypto (golang.org/x/crypto): A collection of packages that provide cryptographic functions and tools for Go.

These packages are imported and utilized within the project to achieve various functionalities.