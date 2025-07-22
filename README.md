## 🛠️ Prerequisites

Before getting started, ensure the following tools are installed on your system:

* [Go](https://golang.org/dl/)
* [MySQL](https://www.mysql.com/)
* [Git](https://git-scm.com/)
* [go-migrate (CLI)](https://github.com/golang-migrate/migrate)

---

## 🚀 Getting Started

Follow these steps to set up and run the project locally:

### 1. Clone the Repository

```bash
git clone git@github.com:boardwallfloor/shv_be.git
cd shv_be
```

---

### 2. Configure the Database

Update your MySQL connection details in the `setDBConn` function inside `/internal/main.go`:

```go
// Example
dsn := "root:your_password@/your_database?parseTime=true"
```
Ensure parseTime is added

> Replace `your_password` and `your_database` with your actual MySQL credentials.

---

### 3. Install Dependencies

Make sure you're in the project root and then run:

```bash
go mod tidy
```

---

### 4. Run Database Migrations

Ensure your MySQL server is running, then apply the database migrations:

```bash
migrate -database "mysql://root:your_password@tcp(localhost:3306)/your_database" -path internal/migrations up
```

> Again, replace `your_password` and `your_database` accordingly.

---

### 5. Run the Application

Start the Go server:

```bash
go run internal/main.go
```

The server will start on **port `8080`** by default.
