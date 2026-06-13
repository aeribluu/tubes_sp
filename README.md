# SecurePass

SecurePass is an application for storing and managing account login information locally and securely (not really). The main data used includes service account data, usernames, and passwords. The application users are individuals who want to secure their digital access.

**Group Member:**
- **103012540009**
    - Anything else that is not `deleteService()`, `selectionSort()`, and `binarySearch()` function. 
- **103012540020**
    - `deleteService()`, `selectionSort()`, and `binarySearch()` function.

---

## Features

| #  | Feature              | Description                                                      |
| -- | -------------------- | ---------------------------------------------------------------- |
| 1  | User Register        | Create a new user account                                        |
| 2  | User Login           | Login using username and password                                |
| 3  | Add Service          | Add saved service data like service name, login ID, and password |
| 4  | View Services        | Show saved services in a table with pagination                   |
| 5  | Detail Service       | View full information of a service                               |
| 6  | Modify Service       | Edit service name, login ID, or password                         |
| 7  | Delete Service       | Remove a service from the saved list                             |
| 8  | Search Service       | Search service using normal search or exact mode                 |
| 9  | Sort Service         | Sort service by name, login ID, created date, or modified date   |
| 10 | Password Strength    | Check saved password strength                                    |
| 11 | Show / Hide Password | User must verify password before viewing saved password          |

---

## Main Menu

When the program starts, it shows:

```text
1. Login
2. Register
3. Exit
```

After login, it shows:

```text
1. Services Saved
2. Add Service
3. Log Out
4. Log Out & Exit
```

---

## Services Saved Menu

Inside the service table, the user can use these options:

| Option | Meaning                      |
| ------ | ---------------------------- |
| `n`    | Next page                    |
| `p`    | Previous page                |
| `s`    | Search service               |
| `o`    | Open sort menu               |
| `r`    | Change rows per page         |
| `i`    | Open detail / modify service |
| `d`    | Delete service               |
| `c`    | Clear search                 |
| `v`    | Show password                |
| `h`    | Hide password                |
| `b`    | Back to main menu            |

---

## Sort Menu

The sort menu has these options:

```text
1. Default Order
2. Service Name Asc
3. Service Name Desc
4. Login ID Asc
5. Login ID Desc
6. Date Created Asc
7. Date Created Desc
8. Date Modified Asc
9. Date Modified Desc
10. Exact Service Name Search Mode
11. Back
```

The exact service name search mode sorts the data by service name ascending, then search will use binary search to find exact service names.

---

## Data Structures

### Service

```go
type Service struct {
    ID                          int
    Service_Name, Login_ID      string
    Password, Password_Strength string
    Date_Created, Date_Modified time.Time
}
```

This struct stores one saved account/service.

| Field               | Description                                   |
| ------------------- | --------------------------------------------- |
| `ID`                | Service ID / tracking number                  |
| `Service_Name`      | Name of the service, like Google or GitHub    |
| `Login_ID`          | Email, username, or login ID for the service  |
| `Password`          | Saved password                                |
| `Password_Strength` | Password strength result, like Weak or Strong |
| `Date_Created`      | Time when service was created                 |
| `Date_Modified`     | Time when service was last changed            |

---

### User

```go
type User struct {
    Username      string
    Service_Count int
    Services      [SMAX]Service
}
```

This struct stores user vault data.

| Field           | Description                      |
| --------------- | -------------------------------- |
| `Username`      | User's username                  |
| `Service_Count` | How many services the user saved |
| `Services`      | Array of service data            |

---

### UserAccount

```go
type UserAccount struct {
    Secret string
    Data   User
}
```

This stores the login password and the user's vault data.

| Field    | Description            |
| -------- | ---------------------- |
| `Secret` | User login password    |
| `Data`   | User data and services |

---

### Session

```go
type Session struct {
    A_Index int
    Vault   *User
}
```

This is used to remember who is currently logged in.

| Field     | Description                      |
| --------- | -------------------------------- |
| `A_Index` | Current user account index       |
| `Vault`   | Pointer to logged in user's data |

---

### ServiceTableWidths

```go
type ServiceTableWidths struct {
    No_Width,
    Service_Name_Width,
    Login_ID_Width,
    Password_Width,
    Date_Created_Width,
    Date_Modified_Width int
}
```

This is only for table display. It helps the table column width adjust based on the text length.

---

## Constants

| Constant            | Value | Meaning                       |
| ------------------- | ----- | ----------------------------- |
| `StateExit`         | `-2`  | Program exit state            |
| `StateUnauthorized` | `-1`  | No user logged in             |
| `SMAX`              | `255` | Maximum service data per user |
| `UMAX`              | `8`   | Maximum user accounts         |

---

## Algorithms

### Searching

| Search Type       | Function           | Used When                      |
| ----------------- | ------------------ | ------------------------------ |
| Sequential Search | `sequentialSearch` | Normal search mode             |
| Binary Search     | `binarySearch`     | Exact Service Name Search Mode |

Normal search checks several fields:

* Service name
* Login ID
* Date created
* Date modified

Normal search uses `containsFoldASCII`, so it can search partial text and ignores uppercase/lowercase.

Example:

```text
Search: google
```

Can match:

```text
Google
Google2
user@gmail.com
```

Exact service name search mode is different. It uses binary search and only finds service names that exactly match the keyword.

Example:

```text
Search: Google
```

Can match:

```text
Google
```

But it will not match:

```text
Google2
```

This is intentional so binary search has a clear purpose and does not confuse normal search.

---

### Sorting

| Sort Type         | Function                | Used For                           |
| ----------------- | ----------------------- | ---------------------------------- |
| Insertion Sort    | `insertionSortServices` | Ascending sort and default sort    |
| Selection Sort    | `selectionSortServices` | Descending sort                    |
| Apply Sort Helper | `applySort`             | Chooses which sort function to use |

Sort fields:

* ID / Default
* Service Name
* Login ID
* Date Created
* Date Modified

The exact service name search mode also sorts by service name ascending first, because binary search needs the data to be sorted.

---

## Application Flow

```text
main
 └── Welcome Menu
      ├── Login
      ├── Register
      └── Exit

After login:
Main Menu
 ├── Services Saved
 │    ├── View table
 │    ├── Search service
 │    ├── Sort service
 │    ├── Change rows per page
 │    ├── View detail
 │    ├── Modify service
 │    ├── Delete service
 │    └── Show / hide password
 ├── Add Service
 ├── Log Out
 └── Log Out & Exit
```

---

## Important Function List

### Main flow

| Function      | Purpose                                      |
| ------------- | -------------------------------------------- |
| `main`        | Starts program loop and checks session state |
| `welcomeMenu` | Shows login/register/exit menu               |
| `mainMenu`    | Shows menu after user login                  |
| `exit`        | Clears session and exits program             |

---

### User functions

| Function                | Purpose                                                            |
| ----------------------- | ------------------------------------------------------------------ |
| `registerUser`          | Creates new user                                                   |
| `authUser`              | Checks username and password                                       |
| `addUsn`                | Input username                                                     |
| `addPwd`                | Input password and confirmation                                    |
| `verifySessionPassword` | Checks current user password again before showing/editing password |

---

### Service CRUD functions

| Function            | Purpose                          |
| ------------------- | -------------------------------- |
| `addService`        | Adds new service to current user |
| `addServiceName`    | Input service name               |
| `addServiceLoginID` | Input login ID                   |
| `getNextServiceID`  | Gets next service ID             |
| `viewServiceDetail` | Shows one service detail         |
| `modifyService`     | Edits service data               |
| `deleteService`     | Deletes one service              |

---

### Display functions

| Function                       | Purpose                                                        |
| ------------------------------ | -------------------------------------------------------------- |
| `displayPagedServices`         | Shows service table with pagination, search, sort, and options |
| `displayServiceDetail`         | Prints one service detail                                      |
| `displayServiceTableStatus`    | Prints page and sort status                                    |
| `displayServiceTableHeader`    | Prints table header                                            |
| `displayServiceRow`            | Prints service row                                             |
| `displayServiceEmptyRow`       | Prints empty row when there is no result                       |
| `displayServiceTableBorder`    | Prints table border                                            |
| `getDynamicServiceTableWidths` | Adjusts table width from current data                          |
| `fitText`                      | Cuts text if too long                                          |

---

### Search and sort functions

| Function                | Purpose                                               |
| ----------------------- | ----------------------------------------------------- |
| `sequentialSearch`      | Searches services one by one                          |
| `binarySearch`          | Searches exact service name using binary search       |
| `chooseSortMenu`        | Shows sort menu                                       |
| `applySort`             | Applies selected sort                                 |
| `insertionSortServices` | Sorts services with insertion sort                    |
| `selectionSortServices` | Sorts services with selection sort                    |
| `compareService`        | Compares service data by selected field               |
| `compareText`           | Compares text using lowercase                         |
| `compareTime`           | Compares date/time                                    |
| `toLower`               | Converts text to lowercase                            |
| `containsFoldASCII`     | Checks text contains keyword without case sensitivity |

---

### Password strength

| Function            | Purpose                                    |
| ------------------- | ------------------------------------------ |
| `calculateStrength` | Checks password and returns strength level |

Password strength result:

| Return | Text      |
| ------ | --------- |
| `'0'`  | Very Weak |
| `'1'`  | Weak      |
| `'2'`  | Good      |
| `'3'`  | Strong    |

The password is checked based on length, uppercase letter, lowercase letter, number, and symbol.

---

## Data Storage

The program stores all data in memory using global arrays:

```go
var user_accounts [UMAX]UserAccount
var user_count int
```

Because of that:

* Data exists only while the program is running
* If the program closes, all users and services are gone
* There is no file save yet
* There is no database yet

---
:::
