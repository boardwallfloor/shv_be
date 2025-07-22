## 🛠️ Prerequisites

Make sure you have the following installed on your system:

* [Go](https://golang.org/dl/)
* [MySQL](https://www.mysql.com/)
* [Git](https://git-scm.com/)
* [go-migrate](https://github.com/golang-migrate/migrate) (CLI tool)

---

## 🚀 Getting Started

Follow these steps to run the project locally:

### 1. Clone the Repository

```bash
git clone git@github.com:boardwallfloor/shv_be.git
cd shv_be
```

### 2. Configure the Database

Open `/internal/main.go` and update the `setDBConn` function with your MySQL credentials.

```go
// Example
dsn := "root:your_password@tcp(localhost:3306)/your_database"
```

### 3. Run Migrations

Ensure your MySQL server is running, then run the following command from the project root:

```bash
migrate -database "mysql://root:your_password@tcp(localhost:3306)/your_database" -path internal/migrations up
```

> Replace `your_password` and `your_database` with your actual credentials.

### 4. Run the Application

From the project root, start the application:

```bash
go run internal/main.go
```
Server will run at port `8080` by default
