# Accounting-System

## How to Run

Follow the steps below to set up and run the tests for the project:

### 1. Set Up Environment Variables

- Copy the `.env.test.example` file to `.env.test` and replace your own environment variables

### 2. Set Up the Database

Apply the necessary SQL migration files to your database:

```bash
psql -U your_user -d your_database -f db/sql/001_create_dl_table.sql
psql -U your_user -d your_database -f db/sql/002_create_sl_table.sql
psql -U your_user -d your_database -f db/sql/003_create_voucher_table.sql
psql -U your_user -d your_database -f db/sql/004_create_voucher_item_table.sql
```
### 3. Run Tests

Navigate to the `internal/services` directory and run tests:

```bash
cd internal/services
go test ./...
```


## Description

This project is a final project for the Golang Bootcamp, focusing on building a simple accounting system. The objective is to create a system that allows users to manage detailed and defined entities, which are then used to generate accounting records (vouchers). The system supports CRUD operations for the following entities:

### Entities:
1. **SL (Subsidiary Ledger)**
   - **Fields:**
     - `code` (string)
     - `title` (string)
     - `hasDL` (boolean)

2. **DL (Detail Ledger)**
   - **Fields:**
     - `code` (string)
     - `title` (string)

3. **Voucher**
   - **Fields:**
     - `number` (string)

4. **VoucherItem**
   - **Fields:**
     - `voucher_id` (refrence to voucher)
     - `sl_id` (refrence to sl)
     - `dl_id` (refrence to dl)
     - `debit_amount` (integer)
     - `credit_amount` (integer)

If you want to see how the request structures should be implemented and detailed further, please refer to the **Description.pdf** file for additional information and examples.

## Additional Information
- The project is implemented service-based.

- It is implemented avoiding design patterns and interfaces to take advantage of Go efficiency.

- There is a high emphasis on clean coding.

- The tests are written in BDD (Behavior-Driven Development).

- It involves interaction with a PostgreSQL (psql) database using GORM.

- It is implemented using a database-first approach.

- It uses a connection pool for database connections.

- The tests are implemented so cleanly that they serve as documentation.

- The tests are not unit tests but system tests, ensuring the entire system, including database interactions, is tested.