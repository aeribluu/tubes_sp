package main

import (
	"fmt"
	"time"
)

const StateExit = -2
const StateUnauthorized = -1
const SMAX = 255
const UMAX = 8

type Service struct {
	Name, Email, Password       string
	Date_Created, Last_Modified time.Time
}

type User struct {
	Username      string
	Service_Count int
	Services      [SMAX]Service
}

type UserAccount struct {
	Secret string
	Data   User
}

type Session struct {
	A_Index int
	Vault   *User
}

var user_accounts [UMAX]UserAccount
var user_count int

func main() {
	var current_session Session
	current_session.A_Index = -1
	current_session.Vault = nil

	for current_session.A_Index != StateExit {
		if current_session.A_Index == StateUnauthorized {
			welcomeMenu(&current_session)
		} else {
			mainMenu(&current_session)
		}
	}
}

func exit(current_session *Session) {
	clearScreen()
	printHeader("Exitting")
	displayPause("See you.", "exit")
	current_session.A_Index = -2
	current_session.Vault = nil
	clearScreen()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printHeader(section string) {
	clearScreen()
	fmt.Println("SecurePass")
	fmt.Printf(">> %v\n\n", section)
}

func displayPause(msg, to_what string) {
	var pause string
	fmt.Println(msg)
	fmt.Printf("Press Enter to %v...", to_what)
	fmt.Scanln(&pause)
}

// TODOTODO
func displayPagedServices(current_session *Session, mode, sorted_by string) string {
	var target, option string
	var page, page_max, rows_per_page, page_lo, page_hi int
	var isSearching bool

	option = "?"
	rows_per_page = 5
	for option == "?" {
		printHeader(mode)

		//TODO
		if target == "" {
			page_max = current_session.Vault.Service_Count / rows_per_page
			if current_session.Vault.Service_Count%rows_per_page != 0 {
				page_max = page_max + 1
			}

			page_lo = page * rows_per_page
			page_hi = page_lo + rows_per_page - 1
			if page_hi > current_session.Vault.Service_Count {
				page_hi = current_session.Vault.Service_Count
			}
		} else {
			// binarySearchByName(current_session, target, 0, current_session.Vault.Service_Count)
			//sequentialSearchByName(current_session, target)
		}

		// TODO		: Display table here
		fmt.Printf("%-30s%d ROWS PER PAGE\n", fmt.Sprintf("PAGE: %02d", page+1), rows_per_page)
		fmt.Printf("%-4v %-13v %-19v\n", "No.", "Service Name", "Last Modified")

		// i = lo
		// for i < hi {
		//	i = i + 1
		// }
		// Rows		: fmt.Printf("%03d. %-13v %-19v\n", i+1, current_session.Vault.Services[i].Name, current_session.Vault.Services[i].Last_Modified.Format("2006-01-02 15:04:05"))

		if !isSearching {
			fmt.Printf("\nNext [n], Previous [p], Search [s], Back [b]: ")
			fmt.Scanln(&option)
		} else {
			fmt.Printf("\nSearch: ")
			fmt.Scanln(&target)
			isSearching = false
		}

		switch option {
		case "n", "N":
			// TODO: Logic for next
			if page < page_max {
				page = page + 1
			}
		case "p", "P":
			// TODO: Logic for previous
			if page > 1 {
				page = page - 1
			}
		case "s", "S":
			// TODO: Logic to search
			isSearching = true
			option = "?"
		case "b", "B":
			// TODO: Logic to go back either to 1. reset search; 2. back to main menu
		default:
			option = "?"
		}
	}

	return "0"
}

func readOption(choice *string) {
	fmt.Printf("\nChoose an option: ")
	fmt.Scanln(choice)
}

func isOptionConfirmed(section, msg string) bool {
	var result bool = false

	is_error := false
	var response string = "0"
	for response == "0" {
		printHeader(section + " >> Verification")

		if is_error {
			fmt.Println("Error: Enter [y] for affirmative or [n] to abort.")
		}
		fmt.Printf("%v (y/N) > ", msg)
		fmt.Scanln(&response)

		switch response {
		case "y", "Y":
			result = true
		case "n", "N":
			result = false
		default:
			is_error = true
			response = "0"
		}
	}
	return result
}

func addUsn(mode string) string {
	var tusn string

	for tusn == "" {
		printHeader(mode)

		fmt.Printf("Username: ")
		fmt.Scanln(&tusn)

		i := 0
		usn_taken := false
		for (i < user_count) && (!usn_taken) {
			usn_taken = user_accounts[i].Data.Username == tusn
			i++
		}

		if usn_taken {
			printHeader(mode)
			displayPause("Error: Username has been taken.", "try again")
			tusn = ""
		}
	}
	return tusn
}

func addPwd(mode string, current_user string) string {
	var tpwd, tpwdc string

	for tpwd == "" || tpwdc == "" || tpwd != tpwdc {
		clearScreen()
		printHeader(mode)
		fmt.Printf("Username: %v\n", current_user)

		if tpwd == "" {
			fmt.Printf("Password: ")
			fmt.Scanln(&tpwd)

			if tpwd != "" && calculateStrength(tpwd) == '0' {
				fmt.Println()
				displayPause("Error: Passwords too weak. 8+ chars, Upper, Lower, Num, Symbol required.", "try again")
				tpwd = ""
			}

		} else {
			fmt.Printf("Password: %v\n", tpwd)

			fmt.Printf("Confirm Password: ")
			fmt.Scanln(&tpwdc)

			if tpwdc != "" && tpwd != tpwdc {
				fmt.Println()
				displayPause("Error: Passwords do not match.", "try again")
				tpwd = ""
				tpwdc = ""
			}
		}
	}
	return tpwd
}

func calculateStrength(pwd string) rune {
	var result rune = '0'
	var upper_count, lower_count, digit_count, symbol_count int
	length := len(pwd)

	if length >= 8 {
		for i := 0; i < length; i++ {
			c := pwd[i]
			switch {
			case c >= 'A' && c <= 'Z':
				upper_count++
			case c >= 'a' && c <= 'z':
				lower_count++
			case c >= '0' && c <= '9':
				digit_count++
			default:
				symbol_count++
			}
		}
	}

	if upper_count >= 1 && lower_count >= 1 && digit_count >= 1 && symbol_count >= 1 {
		if length >= 14 && upper_count >= 2 && lower_count >= 2 && digit_count >= 2 && symbol_count >= 2 {
			result = '3'
		} else if length >= 10 {
			result = '2'
		} else {
			result = '1'
		}
	}

	return result
}

func authUser() Session {
	var tusn, tpwd string

	var activeSession Session
	activeSession.A_Index = -1
	activeSession.Vault = nil

	var login_state rune = 'L'
	for login_state == 'L' {
		printHeader("Authenticating")
		fmt.Printf("Username: ")
		fmt.Scanln(&tusn)
		fmt.Printf("Password: ")
		fmt.Scanln(&tpwd)

		i := 0
		authenticated := false
		for i < user_count && !authenticated {
			if user_accounts[i].Data.Username == tusn && user_accounts[i].Secret == tpwd {
				activeSession.A_Index = i
				activeSession.Vault = &user_accounts[i].Data
				authenticated = true
			}
			i++
		}

		if authenticated {
			printHeader("Authenticating")
			displayPause("Login successful.", "continue")
			login_state = 'S'
		} else {
			printHeader("Authenticating")
			displayPause("Error: User information is incorrect.", "return to welcome menu")
			login_state = 'F'
		}
	}
	return activeSession
}

func registerUser() {
	var tusn, tpwd string
	var choice string

	is_confirmed := false
	for !is_confirmed {
		if tusn == "" {
			tusn = addUsn("Adding New User")
		}
		if tpwd == "" {
			tpwd = addPwd("Adding New User", tusn)
		}

		choice = "0"
		for choice == "0" {
			printHeader("Adding New User >> Confirming User Information")
			fmt.Printf("Username: %v\n", tusn)
			fmt.Printf("Password: %v\n\n", tpwd)

			fmt.Println("1. Edit username")
			fmt.Println("2. Edit password")
			fmt.Println("3. Confirm & Add User")
			fmt.Println("4. Cancel & Start Over")
			fmt.Println("5. Cancel & Return to Welcome Menu")
			readOption(&choice)

			switch choice {
			case "1":
				tusn = addUsn("Adding New User >> Editing Username")
				choice = "0"
			case "2":
				tpwd = addPwd("Adding New User >> Editing Password", tusn)
				choice = "0"
			case "3":
				if isOptionConfirmed("Adding New User >> Confirming User Information", "Proceed to confirm and add user?") {
					is_confirmed = true
				} else {
					choice = "0"
				}
			case "4":
				if isOptionConfirmed("Adding New User >> Confirming User Information", "Proceed to cancel and start over?") {
					tusn = ""
					tpwd = ""

					printHeader("Adding New User >> Resetting")
					displayPause("Registration reset.", "start over")
				} else {
					choice = "0"
				}
			case "5":
				if isOptionConfirmed("Adding New User >> Confirming User Information", "Proceed to cancel and return to welcome menu?") {
					is_confirmed = true
				} else {
					choice = "0"
				}
			default:
				choice = "0"
			}
		}
	}

	if choice == "5" {
		printHeader("Adding New User")
		displayPause("Aborting.", "return to welcome menu")
	} else if user_count >= UMAX {
		printHeader("Adding New User")
		displayPause("Error: Maximum user reached.", "return to welcome menu")
	} else {
		user_accounts[user_count].Secret = tpwd
		user_accounts[user_count].Data.Username = tusn
		user_accounts[user_count].Data.Service_Count = 0

		user_count = user_count + 1

		printHeader("Adding New User")
		displayPause("Registration successful! User added.", "return to welcome menu")
	}
}

func insertionSortByNameAs(current_session *Session) {
	if !(current_session.Vault == nil || current_session.Vault.Service_Count <= 1) {
		for i := 1; i < current_session.Vault.Service_Count; i++ {
			temp := current_session.Vault.Services[i]
			j := i - 1
			for j >= 0 && (current_session.Vault.Services[j].Name > temp.Name) {
				current_session.Vault.Services[j+1] = current_session.Vault.Services[j]
				j--
			}
			current_session.Vault.Services[j+1] = temp
		}
	}
}

func insertionSortByLMTDs(current_session *Session) {
	if !(current_session.Vault == nil || current_session.Vault.Service_Count <= 1) {
		for i := 1; i < current_session.Vault.Service_Count; i++ {
			temp := current_session.Vault.Services[i]
			j := i - 1
			for j >= 0 && (current_session.Vault.Services[j].Last_Modified.Before(temp.Last_Modified)) {
				current_session.Vault.Services[j+1] = current_session.Vault.Services[j]
				j--
			}
			current_session.Vault.Services[j+1] = temp
		}
	}
}

func binarySearchByName(current_session *Session, target_name string, lo, hi int) int {
	var found_index int = -1
	if current_session.Vault != nil {
		for lo <= hi && found_index == -1 {
			mid := lo + (hi-lo)/2

			if current_session.Vault.Services[mid].Name == target_name {
				found_index = mid
			} else if current_session.Vault.Services[mid].Name > target_name {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
	}
	return found_index
}

//func sequentialSearch(current_session *Session, target_name string) [SMAX]int {
//	var found_index [SMAX]int
//	if current_session.Vault != nil {
//		i := 0
//		j := 0
//		for i < current_session.Vault.Service_Count {
//			if (current_session.Vault.Services[i].Name == target_name) || (current_session.Vault.Services[i]).Last_Modified == found_index) {
//				found_index[j] = i
//				j = j + 1
//			}
//			i = i + 1
//		}
//	}
//	return found_index
//}

func welcomeMenu(current_session *Session) {
	var choice string = "0"
	for choice == "0" {
		printHeader("Welcome Menu")
		fmt.Printf("1. Login\n2. Register\n3. Exit\n")
		readOption(&choice)

		switch choice {
		case "1":
			*current_session = authUser()
			if current_session.A_Index == -1 {
				choice = "0"
			}
		case "2":
			registerUser()
			choice = "0"
		case "3":
			exit(current_session)
		default:
			choice = "0"
		}
	}
}

// TODO
func mainMenu(current_session *Session) {
	var choice string = "0"
	for choice == "0" {
		printHeader("Main Menu")
		fmt.Printf("1. Services Saved\n2. Add Service\n3. Log Out\n4. Log Out & Exit")
		readOption(&choice)

		switch choice {
		case "1":
			//TODO
			servicesSavedMenu(current_session)
		case "2":
			//TODO
		case "3":
			current_session.A_Index = -1
			current_session.Vault = nil
			printHeader("Main Menu >> Unauthenticating")
			displayPause("Logging out", "log out")
		case "4":
			exit(current_session)
		default:
			choice = "0"
		}
	}
}

func servicesSavedMenu(current_session *Session) {
	var choice string = "0"

	//if (current_session.Vault != nil) && (current_session.Vault.Service_Count > 0) {
	for choice == "0" {
		printHeader("Services Saved >> Options")
		fmt.Printf("1. Sorted by Name\n2. Sorted by Time\n3. Back to Main Menu")
		readOption(&choice)

		switch choice {
		case "1":
			insertionSortByNameAs(current_session)
			choice = displayPagedServices(current_session, "Services Saved >> Sorted by Name", "Name")
		case "2":
			insertionSortByLMTDs(current_session)
			choice = displayPagedServices(current_session, "Services Saved >> Sorted by Last Modified Time", "LMT")
		case "3":
			//TODO
		default:
			choice = "0"
		}
	}
	// } else {
	//	printHeader("Services Saved")
	//	displayPause("No services saved yet.", "go back to main menu")
	// }
}
