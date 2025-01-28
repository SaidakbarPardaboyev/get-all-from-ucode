

Here is a detailed README file for the `get-all` package:


**README.md**
===============

**get-all** Package
-------------------

The **get-all** package is a utility package designed to simplify data retrieval from various data sources. It provides a unified interface for fetching data from different databases, making it easier to switch between data sources or use multiple sources simultaneously.

**Features**
------------

*   Supports multiple database types, including MongoDB and PostgreSQL
*   Provides a simple and consistent API for data retrieval
*   Allows for filtering, sorting, and limiting of data
*   Supports counting of records matching a filter

**Usage**
---------

### Importing the Package

To use the `get-all` package, simply import it into your Go program:
```go
import "get_all"
```
### Creating a New Object Instance

To create a new instance of the `object` struct, pass in your database configuration:
```go
obj := get_all.object{
    config: &get_all.Config{
        DB_HOST:     "your_host",
        DB_PORT:     "your_port",
        DB_USER:     "your_user",
        DB_PASSWORD: "your_password",
        DB_NAME:     "your_database",
        DB_TYPE:     "your_database_type",
    },
}
```
### Retrieving Data

To retrieve data from a database, use the `Items` method:
```go
items, err := obj.Items("your_collection").GetAll().Count()
if err != nil {
    // Handle error
}
```
You can also apply filters to the data:
```go
filter := map[string]interface{}{
    "field": "value",
}
items, err := obj.Items("your_collection").GetAll().Filter(filter).Count()
if err != nil {
    // Handle error
}
```
Sort the data:
```go
sort := map[string]interface{}{
    "field": 1,
}
items, err := obj.Items("your_collection").GetAll().Sort(sort).Count()
if err != nil {
    // Handle error
}
```
Limit the number of records returned:
```go
limit := int64(10)
items, err := obj.Items("your_collection").GetAll().Limit(limit).Count()
if err != nil {
    // Handle error
}
```
### Counting Records

To count the number of records matching a filter, use the `Count` method:
```go
count, err := obj.Items("your_collection").GetAll().Filter(filter).Count()
if err != nil {
    // Handle error
}
```
### Executing Queries

To execute queries on MongoDB and PostgreSQL databases, use the `execMongo` and `execPostgres` methods:
```go
mongoResults, err := obj.Items("your_collection").GetAll().execMongo()
if err != nil {
    // Handle error
}

postgresResults, err := obj.Items("your_collection").GetAll().execPostgres()
if err != nil {
    // Handle error
}
```
**API Documentation**
--------------------

The **get-all** package provides the following API:

### object

*   `Items(collection string) ItemsI`: Retrieves data from the specified collection.
*   `Config() *Config`: Returns the database configuration.

### ItemsI

*   `GetAll() *GetAllI`: Retrieves all data from the collection.
*   `Filter(filter map[string]interface{}) *GetAllI`: Applies a filter to the data.
*   `Sort(sort map[string]interface{}) *GetAllI`: Sorts the data.
*   `Limit(limit int64) *GetAllI`: Limits the number of records returned.

### GetAllI

*   `Count() (int64, error)`: Returns the number of records matching the filter.
*   `execMongo() ([]map[string]interface{}, error)`: Executes the query on a MongoDB database.
*   `execPostgres() ([]map[string]interface{}, error)`: Executes the query on a PostgreSQL database.

**Contributing**
---------------

Contributions to the **get-all** package are welcome. If you find a bug or would like to add a new feature, please submit a pull request.

**License**
----------

The **get-all** package is licensed under the MIT License.