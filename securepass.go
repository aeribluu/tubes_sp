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
	ID                          int
	Service_Name, Login_ID      string
	Password, Password_Strength string
	Date_Created, Date_Modified time.Time
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

type ServiceTableWidths struct {
	No_Width,
	Service_Name_Width,
	Login_ID_Width,
	Password_Width,
	Date_Created_Width,
	Date_Modified_Width int
}

var user_accounts [UMAX]UserAccount
var user_count int

func main() {
	var current_session Session
	current_session.A_Index = StateUnauthorized
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
	current_session.A_Index = StateExit
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

func displayPagedServices(current_session *Session, mode, sorted_by string) string {
	var target, option, sort_order string
	var page, page_max, rows_per_page, page_lo, page_hi, service_no, found_count int
	var found_idx [SMAX]int
	var show_pwd bool

	option = "?"
	rows_per_page = 5
	if sorted_by == "" {
		sorted_by = "ID"
	}
	sort_order = "Asc"
	applySort(current_session, sorted_by, sort_order)

	for option == "?" {
		printHeader(mode)

		if current_session.Vault == nil || current_session.Vault.Service_Count == 0 {
			displayPause("No services saved yet.", "go back")
			option = "0"
		} else {
			if target == "" {
				found_count = current_session.Vault.Service_Count
				i := 0
				for i < found_count {
					found_idx[i] = i
					i = i + 1
				}
			} else {
				if sorted_by == "Exact Service Name Search Mode" {
					idx := binarySearch(current_session, target)
					found_count = 0

					if idx != -1 {
						left := idx

						for left > 0 && toLower(current_session.Vault.Services[left-1].Service_Name) == toLower(target) {
							left = left - 1
						}

						for left < current_session.Vault.Service_Count && toLower(current_session.Vault.Services[left].Service_Name) == toLower(target) {
							found_idx[found_count] = left
							found_count = found_count + 1
							left = left + 1
						}
					}
				} else {
					sequentialSearch(current_session, target, &found_idx, &found_count)
				}
			}

			page_max = found_count / rows_per_page
			if found_count%rows_per_page != 0 {
				page_max = page_max + 1
			}
			if page >= page_max {
				page = page_max - 1
			}
			if page < 0 {
				page = 0
			}

			page_lo = page * rows_per_page
			page_hi = page_lo + rows_per_page
			if page_hi > found_count {
				page_hi = found_count
			}

			table_widths := getDynamicServiceTableWidths(current_session, &found_idx, found_count, show_pwd)
			display_sort := fmt.Sprintf("%v %v", sorted_by, sort_order)
			if sorted_by == "ID" {
				display_sort = "Default"
			} else if sorted_by == "Exact Service Name Search Mode" {
				display_sort = "Exact Service Name Search Mode"
			}

			displayServiceTableStatus(page, page_max, rows_per_page, getServiceTableWidth(table_widths, show_pwd), display_sort)
			if target != "" {
				fmt.Printf("Search: %v (%d result(s))\n", target, found_count)
			}
			displayServiceTableHeader(table_widths, show_pwd)

			if found_count == 0 {
				displayServiceEmptyRow(table_widths, show_pwd)
			} else {
				i := page_lo
				for i < page_hi {
					service_idx := found_idx[i]
					displayServiceRow(service_idx+1, current_session.Vault.Services[service_idx], table_widths, show_pwd)
					i = i + 1
				}
			}
			displayServiceTableBorder(table_widths, show_pwd)

			fmt.Printf("\n")
			if page+1 < page_max {
				fmt.Printf("Next [n], ")
			}
			if page > 0 {
				fmt.Printf("Previous [p], ")
			}
			fmt.Printf("Search [s], Sort [o], Rows [r], Detail/Modify [i], Delete [d]")
			if target != "" {
				fmt.Printf(", Clear Search [c]")
			}
			if show_pwd {
				fmt.Printf(", Hide Password [h]")
			} else {
				fmt.Printf(", Show Password [v]")
			}
			fmt.Printf(", Back [b]: ")
			fmt.Scanln(&option)

			switch option {
			case "n", "N":
				if page+1 < page_max {
					page = page + 1
				}
				option = "?"
			case "p", "P":
				if page > 0 {
					page = page - 1
				}
				option = "?"
			case "s", "S":
				fmt.Printf("Search: ")
				fmt.Scanln(&target)
				page = 0
				option = "?"
			case "c", "C":
				target = ""
				page = 0
				option = "?"
			case "o", "O":
				chooseSortMenu(current_session, mode, &sorted_by, &sort_order, show_pwd)
				page = 0
				option = "?"
			case "r", "R":
				changeRowsPerPage(mode, &rows_per_page)
				page = 0
				option = "?"
			case "i", "I":
				fmt.Printf("Service Number: ")
				fmt.Scanln(&service_no)
				viewServiceDetail(current_session, mode, service_no, &show_pwd)
				applySort(current_session, sorted_by, sort_order)
				service_no = 0
				option = "?"
			case "d", "D":
				deleteService(current_session, mode)
				applySort(current_session, sorted_by, sort_order)
				page = 0
				option = "?"
			case "v", "V":
				if !show_pwd {
					if verifySessionPassword(current_session, mode) {
						show_pwd = true
						printHeader(mode + " >> Password Column")
						displayPause("Password column is now visible.", "continue")
					} else {
						printHeader(mode + " >> Password Column")
						displayPause("Error: User password is incorrect.", "go back")
					}
				}
				option = "?"
			case "h", "H":
				if show_pwd {
					show_pwd = false
					if sorted_by == "Password" {
						sorted_by = "ID"
						sort_order = "Asc"
						applySort(current_session, sorted_by, sort_order)
					}
				}
				option = "?"
			case "b", "B":
				if target != "" {
					target = ""
					page = 0
					option = "?"
				} else {
					option = "0"
				}
			default:
				option = "?"
			}
		}
	}

	return "0"
}

func displayServiceDetail(service Service, show_pwd, show_pwd_strength, show_dates bool) {
	fmt.Printf("Service Name: %v\n", service.Service_Name)
	fmt.Printf("Login ID: %v\n", service.Login_ID)

	if show_pwd {
		fmt.Printf("Password: %v\n", service.Password)
	} else {
		fmt.Printf("Password: ********\n")
	}
	if show_pwd_strength {
		fmt.Printf("Password Strength: %v\n", service.Password_Strength)
	}
	if show_dates {
		fmt.Printf("Date Created: %v\n", service.Date_Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Date Modified: %v\n", service.Date_Modified.Format("2006-01-02 15:04:05"))
	}
}

func displayServiceTableStatus(page, page_max, rows_per_page, table_width int, display_sort string) {
	left_text := fmt.Sprintf("PAGE: %02d/%02d", page+1, page_max)
	mid_text := fmt.Sprintf("SORT: %v", display_sort)
	right_text := fmt.Sprintf("%d ROWS PER PAGE", rows_per_page)
	left_mid_text := fmt.Sprintf("%v    %v", left_text, mid_text)

	if len(left_mid_text)+len(right_text) >= table_width {
		fmt.Printf("%v %v %v\n", left_text, mid_text, right_text)
	} else {
		fmt.Printf("%-*v%v\n", table_width-len(right_text), left_mid_text, right_text)
	}
}

func changeRowsPerPage(mode string, rows_per_page *int) {
	var new_rows_per_page int

	printHeader(mode + " >> Rows Per Page")
	fmt.Printf("Rows per page (1-%d): ", SMAX)
	fmt.Scanln(&new_rows_per_page)

	if new_rows_per_page < 1 || new_rows_per_page > SMAX {
		printHeader(mode + " >> Rows Per Page")
		displayPause("Error: Invalid rows per page.", "go back")
	} else {
		*rows_per_page = new_rows_per_page
	}
}

func getDynamicServiceTableWidths(current_session *Session, found_idx *[SMAX]int, found_count int, show_pwd bool) ServiceTableWidths {
	var table_widths ServiceTableWidths

	// MIN WIDTH
	table_widths.No_Width = 4
	table_widths.Service_Name_Width = len("Service Name")
	table_widths.Login_ID_Width = len("Login ID")
	table_widths.Password_Width = len("Password")
	table_widths.Date_Created_Width = len("Date Created")
	table_widths.Date_Modified_Width = len("Date Modified")

	if current_session != nil && current_session.Vault != nil {
		i := 0
		for i < found_count {
			service_idx := found_idx[i]
			if service_idx >= 0 && service_idx < current_session.Vault.Service_Count {
				service := current_session.Vault.Services[service_idx]
				date_created := service.Date_Created.Format("2006-01-02 15:04:05")
				last_modified := service.Date_Modified.Format("2006-01-02 15:04:05")

				table_widths.Service_Name_Width = getMaxColumnWidth(table_widths.Service_Name_Width, len(service.Service_Name))
				table_widths.Login_ID_Width = getMaxColumnWidth(table_widths.Login_ID_Width, len(service.Login_ID))
				if show_pwd {
					table_widths.Password_Width = getMaxColumnWidth(table_widths.Password_Width, len(service.Password))
				}
				table_widths.Date_Created_Width = getMaxColumnWidth(table_widths.Date_Created_Width, len(date_created))
				table_widths.Date_Modified_Width = getMaxColumnWidth(table_widths.Date_Modified_Width, len(last_modified))
			}
			i = i + 1
		}
	}

	return table_widths
}

func getMaxColumnWidth(current_width, text_length int) int {
	var result int

	result = current_width
	if text_length > result {
		result = text_length
	}
	if result > 25 {
		result = 25
	}

	return result
}

func displayServiceTableHeader(table_widths ServiceTableWidths, show_pwd bool) {
	displayServiceTableBorder(table_widths, show_pwd)
	printServiceTableColumn("No.", table_widths.No_Width)
	printServiceTableColumn("Service Name", table_widths.Service_Name_Width)
	printServiceTableColumn("Login ID", table_widths.Login_ID_Width)
	if show_pwd {
		printServiceTableColumn("Password", table_widths.Password_Width)
	}
	printServiceTableColumn("Date Created", table_widths.Date_Created_Width)
	printServiceTableColumn("Date Modified", table_widths.Date_Modified_Width)
	fmt.Println("|")
	displayServiceTableBorder(table_widths, show_pwd)
}

func displayServiceRow(no int, service Service, table_widths ServiceTableWidths, show_pwd bool) {
	no_text := fmt.Sprintf("%03d", no)
	date_created := service.Date_Created.Format("2006-01-02 15:04:05")
	last_modified := service.Date_Modified.Format("2006-01-02 15:04:05")

	printServiceTableColumn(no_text, table_widths.No_Width)
	printServiceTableColumn(fitText(service.Service_Name, table_widths.Service_Name_Width), table_widths.Service_Name_Width)
	printServiceTableColumn(fitText(service.Login_ID, table_widths.Login_ID_Width), table_widths.Login_ID_Width)
	if show_pwd {
		printServiceTableColumn(fitText(service.Password, table_widths.Password_Width), table_widths.Password_Width)
	}
	printServiceTableColumn(fitText(date_created, table_widths.Date_Created_Width), table_widths.Date_Created_Width)
	printServiceTableColumn(fitText(last_modified, table_widths.Date_Modified_Width), table_widths.Date_Modified_Width)
	fmt.Println("|")
}

func displayServiceEmptyRow(table_widths ServiceTableWidths, show_pwd bool) {
	inner_width := getServiceTableWidth(table_widths, show_pwd) - 4
	fmt.Printf("| %-*v |\n", inner_width, fitText("No matching services found.", inner_width))
}

func displayServiceTableBorder(table_widths ServiceTableWidths, show_pwd bool) {
	fmt.Print("+")
	printBorderPart(table_widths.No_Width)
	printBorderPart(table_widths.Service_Name_Width)
	printBorderPart(table_widths.Login_ID_Width)
	if show_pwd {
		printBorderPart(table_widths.Password_Width)
	}
	printBorderPart(table_widths.Date_Created_Width)
	printBorderPart(table_widths.Date_Modified_Width)
	fmt.Println()
}

func printServiceTableColumn(text string, width int) {
	fmt.Printf("| %-*v ", width, text)
}

func printBorderPart(width int) {
	i := 0
	for i < width+2 {
		fmt.Print("-")
		i = i + 1
	}
	fmt.Print("+")
}

func getServiceTableWidth(table_widths ServiceTableWidths, show_pwd bool) int {
	var table_width int

	table_width = 1
	table_width = table_width + table_widths.No_Width + 3
	table_width = table_width + table_widths.Service_Name_Width + 3
	table_width = table_width + table_widths.Login_ID_Width + 3
	if show_pwd {
		table_width = table_width + table_widths.Password_Width + 3
	}
	table_width = table_width + table_widths.Date_Created_Width + 3
	table_width = table_width + table_widths.Date_Modified_Width + 3

	return table_width
}

func fitText(text string, width int) string {
	var ntext string
	if width <= 0 {
		ntext = ""
	} else if len(text) <= width {
		ntext = text
	} else if width <= 3 {
		i := 0
		for i < width {
			ntext = ntext + "."
			i = i + 1
		}
	} else {
		i := 0
		for i < width-3 {
			ntext = ntext + string(text[i])
			i = i + 1
		}
		ntext = ntext + "..."
	}

	return ntext
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

func verifySessionPassword(current_session *Session, mode string) bool {
	var result bool
	var tpwd string

	printHeader(mode + " >> Password Verification")
	if current_session == nil || current_session.Vault == nil || current_session.A_Index < 0 || current_session.A_Index >= user_count {
		displayPause("Error: No active session.", "go back")
	} else {
		fmt.Printf("User Password: ")
		fmt.Scanln(&tpwd)
		if user_accounts[current_session.A_Index].Secret == tpwd {
			result = true
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
			i = i + 1
		}

		if usn_taken {
			printHeader(mode)
			displayPause("Error: Username has been taken.", "try again")
			tusn = ""
		}
	}
	return tusn
}

func addPwd(mode, current_label, current_value, extra_label, extra_value string, check_strength bool) string {
	var tpwd, tpwdc string

	for tpwd == "" || tpwdc == "" || tpwd != tpwdc {
		clearScreen()
		printHeader(mode)

		if current_label != "" {
			fmt.Printf("%v: %v\n", current_label, current_value)
		}

		if extra_label != "" {
			fmt.Printf("%v: %v\n", extra_label, extra_value)
		}

		if tpwd == "" {
			fmt.Printf("Password: ")
			fmt.Scanln(&tpwd)

			if tpwd != "" && check_strength && calculateStrength(tpwd) == '0' {
				fmt.Println()
				displayPause("Error: Password too weak. 8+ chars, Upper, Lower, Num, Symbol required.", "try again")
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

func addServiceName(mode string) string {
	var service_name string

	for service_name == "" {
		printHeader(mode)

		fmt.Printf("Service Name: ")
		fmt.Scanln(&service_name)
	}
	return service_name
}

func addServiceLoginID(mode string, current_service string) string {
	var login_id string

	for login_id == "" {
		printHeader(mode)
		fmt.Printf("Service Name: %v\n", current_service)

		fmt.Printf("Login ID: ")
		fmt.Scanln(&login_id)
	}
	return login_id
}

func addService(current_session *Session) {
	var tservice Service
	var choice string

	mode := "Adding New Service"
	is_confirmed := false
	for !is_confirmed {
		if tservice.Service_Name == "" {
			tservice.Service_Name = addServiceName(mode)
		}
		if tservice.Login_ID == "" {
			tservice.Login_ID = addServiceLoginID(mode, tservice.Service_Name)
		}
		if tservice.Password == "" {
			tservice.Password = addPwd(mode, "Service Name", tservice.Service_Name, "Login ID", tservice.Login_ID, false)
		}

		choice = "0"
		for choice == "0" {
			printHeader(mode + " >> Confirming Service Information")
			displayServiceDetail(tservice, true, false, false)

			fmt.Println()
			fmt.Println("1. Edit Service name")
			fmt.Println("2. Edit Login ID")
			fmt.Println("3. Edit Password")
			fmt.Println("4. Confirm & Add Service")
			fmt.Println("5. Cancel & Start Over")
			fmt.Println("6. Cancel & Return to Welcome Menu")
			readOption(&choice)

			switch choice {
			case "1":
				tservice.Service_Name = addServiceName(mode + " >> Editing Service Name")
				choice = "0"
			case "2":
				tservice.Login_ID = addServiceLoginID(mode+" >> Editing Login ID", tservice.Service_Name)
				choice = "0"
			case "3":
				tservice.Password = addPwd(mode+" >> Editing Password", "Service Name", tservice.Service_Name, "Login ID", tservice.Login_ID, false)
				choice = "0"
			case "4":
				if isOptionConfirmed(mode+" >> Confirming Service Information", "Proceed to confirm and add service?") {
					is_confirmed = true
				} else {
					choice = "0"
				}
			case "5":
				if isOptionConfirmed(mode+" >> Confirming Service Information", "Proceed to cancel and start over?") {
					tservice.Service_Name = ""
					tservice.Login_ID = ""
					tservice.Password = ""

					printHeader(mode + " >> Resetting")
					displayPause("Service entry reset.", "start over")
				} else {
					choice = "0"
				}
			case "6":
				if isOptionConfirmed(mode+" >> Confirming Service Information", "Proceed to cancel and return to main menu?") {
					is_confirmed = true
				} else {
					choice = "0"
				}
			default:
				choice = "0"
			}
		}
	}

	if choice == "6" {
		printHeader(mode)
		displayPause("Aborting.", "return to main menu")
	} else if current_session.Vault == nil {
		printHeader(mode)
		displayPause("Error: No active session.", "return to main menu")
	} else if current_session.Vault.Service_Count >= SMAX {
		printHeader(mode)
		displayPause("Error: Maximum service reached.", "return to main menu")
	} else {
		now := time.Now()
		tservice.ID = getNextServiceID(current_session)
		tservice.Date_Created = now
		tservice.Date_Modified = now

		switch calculateStrength(tservice.Password) {
		case '1':
			tservice.Password_Strength = "Weak"
		case '2':
			tservice.Password_Strength = "Good"
		case '3':
			tservice.Password_Strength = "Strong"
		default:
			tservice.Password_Strength = "Very Weak"
		}

		service_idx := current_session.Vault.Service_Count
		current_session.Vault.Services[service_idx] = tservice

		current_session.Vault.Service_Count = current_session.Vault.Service_Count + 1
		applySort(current_session, "ID", "Asc")

		printHeader(mode)
		displayPause("Registration successful! Service added.", "return to main menu")
	}
}

func getNextServiceID(current_session *Session) int {
	var next_id int = 1

	if current_session != nil && current_session.Vault != nil {
		i := 0
		for i < current_session.Vault.Service_Count {
			if current_session.Vault.Services[i].ID >= next_id {
				next_id = current_session.Vault.Services[i].ID + 1
			}
			i = i + 1
		}
	}
	return next_id
}

func deleteService(current_session *Session, mode string) {
	var delete_no, delete_idx, last_idx int
	var is_valid bool

	printHeader(mode + " >> Delete Service")
	if current_session == nil || current_session.Vault == nil || current_session.Vault.Service_Count == 0 {
		displayPause("No services saved yet.", "go back")
	} else {
		fmt.Printf("Delete No.: ")
		fmt.Scanln(&delete_no)
		delete_idx = delete_no - 1
		is_valid = delete_idx >= 0 && delete_idx < current_session.Vault.Service_Count

		if !is_valid {
			printHeader(mode + " >> Delete Service")
			displayPause("Error: Invalid service number.", "go back")
		} else {
			printHeader(mode + " >> Delete Service")
			fmt.Printf("No.: %03d\n", delete_no)
			fmt.Printf("Service Name: %v\n", current_session.Vault.Services[delete_idx].Service_Name)
			fmt.Printf("Login ID.: %v\n\n", current_session.Vault.Services[delete_idx].Login_ID)

			if isOptionConfirmed(mode+" >> Delete Service", "Proceed to delete this service?") {
				last_idx = current_session.Vault.Service_Count - 1
				current_session.Vault.Services[delete_idx] = current_session.Vault.Services[last_idx]
				current_session.Vault.Service_Count = current_session.Vault.Service_Count - 1

				printHeader(mode + " >> Delete Service")
				displayPause("Service deleted successfully.", "go back")
			} else {
				printHeader(mode + " >> Delete Service")
				displayPause("Deletion cancelled.", "go back")
			}
		}
	}
}

func viewServiceDetail(current_session *Session, mode string, service_no int, show_pwd *bool) {
	var service_idx int
	var option string

	service_idx = service_no - 1
	if service_idx < 0 || service_idx >= current_session.Vault.Service_Count {
		printHeader(mode + " >> Service Detail")
		displayPause("Error: Invalid service number.", "go back")
	} else {
		option = "?"
		for option == "?" {
			printHeader(mode + fmt.Sprintf(" >> Service No. %03d Details", service_no))

			displayServiceDetail(current_session.Vault.Services[service_idx], *show_pwd, *show_pwd, true)
			fmt.Println()
			if *show_pwd {
				fmt.Printf("Hide Password [h], ")
			} else {
				fmt.Printf("View Password [v], ")
			}
			fmt.Printf("Modify [m], Back[b]: ")
			fmt.Scanln(&option)

			switch option {
			case "v", "V":
				if !*show_pwd {
					if verifySessionPassword(current_session, mode+" >> Service Detail") {
						*show_pwd = true
					} else {
						printHeader(mode + " >> Service Detail")
						displayPause("Error: User password is incorrect.", "go back")
					}
				}
				option = "?"
			case "h", "H":
				*show_pwd = false
				option = "?"
			case "m", "M":
				modifyService(current_session, mode, service_idx)
				option = "?"
			case "b", "B":
				option = "0"
			default:
				option = "?"
			}
		}
	}
}

func modifyService(current_session *Session, mode string, service_idx int) {
	var choice string
	var modified_service Service
	var show_pwd bool

	modified_service = current_session.Vault.Services[service_idx]

	mode = mode + " >> Modify Service"
	choice = "0"
	for choice == "0" {
		printHeader(mode)

		displayServiceDetail(modified_service, show_pwd, false, false)

		fmt.Println()
		fmt.Println("1. Edit Service name")
		fmt.Println("2. Edit Login ID")
		fmt.Println("3. Edit Password")
		fmt.Println("4. Confirm & Save Changes")
		fmt.Println("5. Cancel Changes")
		readOption(&choice)

		switch choice {
		case "1":
			modified_service.Service_Name = addServiceName(mode + " >> Editing Service Name")
			choice = "0"
		case "2":
			modified_service.Login_ID = addServiceLoginID(mode+" >> Editing Login ID", modified_service.Service_Name)
			choice = "0"
		case "3":
			if verifySessionPassword(current_session, mode) {
				modified_service.Password = addPwd(mode+" >> Editing Password", "Service Name", modified_service.Service_Name, "Login ID", modified_service.Login_ID, false)
				switch calculateStrength(modified_service.Password) {
				case '1':
					modified_service.Password_Strength = "Weak"
				case '2':
					modified_service.Password_Strength = "Good"
				case '3':
					modified_service.Password_Strength = "Strong"
				default:
					modified_service.Password_Strength = "Very Weak"
				}
				show_pwd = true
			} else {
				printHeader(mode)
				displayPause("Error: User password is incorrect.", "go back")
			}
			choice = "0"
		case "4":
			if isOptionConfirmed(mode, "Proceeed to save changes?") {
				modified_service.Date_Modified = time.Now()
				current_session.Vault.Services[service_idx] = modified_service

				printHeader(mode)
				displayPause("Service modified successfully.", "go back")
			} else {
				choice = "0"
			}
		case "5":
			if isOptionConfirmed(mode, "Proceeed to cancel changes?") {
				printHeader(mode)
				displayPause("Modification cancelled.", "go back")
			} else {
				choice = "0"
			}
		default:
			choice = "0"
		}
	}
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
	activeSession.A_Index = StateUnauthorized
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
			i = i + 1
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

	mode := "Adding New User"
	is_confirmed := false
	for !is_confirmed {
		if tusn == "" {
			tusn = addUsn(mode)
		}
		if tpwd == "" {
			tpwd = addPwd(mode, "Username", tusn, "", "", true)
		}

		choice = "0"
		for choice == "0" {
			printHeader(mode + " >> Confirming User Information")
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
				tusn = addUsn(mode + " >> Editing Username")
				choice = "0"
			case "2":
				tpwd = addPwd(mode+" >> Editing Password", "Username", tusn, "", "", true)
				choice = "0"
			case "3":
				if isOptionConfirmed(mode+" >> Confirming User Information", "Proceed to confirm and add user?") {
					is_confirmed = true
				} else {
					choice = "0"
				}
			case "4":
				if isOptionConfirmed(mode+" >> Confirming User Information", "Proceed to cancel and start over?") {
					tusn = ""
					tpwd = ""

					printHeader(mode + " >> Resetting")
					displayPause("Registration reset.", "start over")
				} else {
					choice = "0"
				}
			case "5":
				if isOptionConfirmed(mode+" >> Confirming User Information", "Proceed to cancel and return to welcome menu?") {
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
		printHeader(mode + "")
		displayPause("Aborting.", "return to welcome menu")
	} else if user_count >= UMAX {
		printHeader(mode + "")
		displayPause("Error: Maximum user reached.", "return to welcome menu")
	} else {
		user_accounts[user_count].Secret = tpwd
		user_accounts[user_count].Data.Username = tusn
		user_accounts[user_count].Data.Service_Count = 0

		user_count = user_count + 1

		printHeader(mode + "")
		displayPause("Registration successful! User added.", "return to welcome menu")
	}
}

func applySort(current_session *Session, sorted_by string, sort_order string) {
	if sorted_by == "Exact Service Name Search Mode" {
		insertionSortServices(current_session, "Service Name", "Asc")
	} else if sort_order == "Desc" {
		selectionSortServices(current_session, sorted_by, sort_order)
	} else {
		insertionSortServices(current_session, sorted_by, sort_order)
	}
}

func insertionSortServices(current_session *Session, sorted_by, sort_order string) {
	if !(current_session == nil || current_session.Vault == nil || current_session.Vault.Service_Count <= 1) {
		i := 1
		for i < current_session.Vault.Service_Count {
			temp := current_session.Vault.Services[i]
			j := i - 1
			for j >= 0 && shouldServiceMoveAfter(current_session.Vault.Services[j], temp, sorted_by, sort_order) {
				current_session.Vault.Services[j+1] = current_session.Vault.Services[j]
				j = j - 1
			}
			current_session.Vault.Services[j+1] = temp
			i = i + 1
		}
	}
}

func selectionSortServices(current_session *Session, sorted_by, sort_order string) {
	if !(current_session == nil || current_session.Vault == nil || current_session.Vault.Service_Count <= 1) {
		i := 0
		for i < current_session.Vault.Service_Count-1 {
			xtr := i
			j := i + 1
			for j < current_session.Vault.Service_Count {
				if shouldServiceMoveAfter(current_session.Vault.Services[xtr], current_session.Vault.Services[j], sorted_by, sort_order) {
					xtr = j
				}
				j = j + 1
			}
			if xtr != i {
				temp := current_session.Vault.Services[i]
				current_session.Vault.Services[i] = current_session.Vault.Services[xtr]
				current_session.Vault.Services[xtr] = temp
			}
			i = i + 1
		}
	}
}

func shouldServiceMoveAfter(left_service, right_service Service, sorted_by, sort_order string) bool {
	var result bool

	compare_result := compareService(left_service, right_service, sorted_by)
	if sort_order == "Desc" {
		result = compare_result < 0
	} else {
		result = compare_result > 0
	}

	return result
}

func compareService(left_service, right_service Service, sorted_by string) int {
	var result int

	switch sorted_by {
	case "ID":
		if left_service.ID > right_service.ID {
			result = 1
		} else if left_service.ID < right_service.ID {
			result = -1
		}
	case "Service Name":
		result = compareText(left_service.Service_Name, right_service.Service_Name)
	case "Login ID":
		result = compareText(left_service.Login_ID, right_service.Login_ID)
	case "Date Created":
		result = compareTime(left_service.Date_Created, right_service.Date_Created)
	case "Date Modified":
		result = compareTime(left_service.Date_Modified, right_service.Date_Modified)
	default:
		if left_service.ID > right_service.ID {
			result = 1
		} else if left_service.ID < right_service.ID {
			result = -1
		}
	}

	if result == 0 && sorted_by != "ID" {
		if left_service.ID > right_service.ID {
			result = 1
		} else if left_service.ID < right_service.ID {
			result = -1
		}
	}
	return result
}

func compareText(left_text, right_text string) int {
	var result int = 0
	left_text = toLower(left_text)
	right_text = toLower(right_text)

	if left_text > right_text {
		result = 1
	} else if left_text < right_text {
		result = -1
	}
	return result
}

func compareTime(left_time time.Time, right_time time.Time) int {
	var result int = 0

	if left_time.After(right_time) {
		result = 1
	} else if left_time.Before(right_time) {
		result = -1
	}
	return result
}

func chooseSortMenu(current_session *Session, mode string, sorted_by *string, sort_order *string, show_pwd bool) {
	var choice string = "0"
	var sort_changed bool

	for choice == "0" {
		sort_changed = false
		printHeader(mode + " >> Sort")
		fmt.Println("1. Default Order")
		fmt.Println("2. Service Name Asc")
		fmt.Println("3. Service Name Desc")
		fmt.Println("4. Login ID Asc")
		fmt.Println("5. Login ID Desc")
		fmt.Println("6. Date Created Asc")
		fmt.Println("7. Date Created Desc")
		fmt.Println("8. Date Modified Asc")
		fmt.Println("9. Date Modified Desc")
		fmt.Println("10. Exact Service Name Search Mode")
		fmt.Println("11. Back")
		readOption(&choice)

		switch choice {
		case "1":
			*sorted_by = "ID"
			*sort_order = "Asc"
			sort_changed = true
		case "2":
			*sorted_by = "Service Name"
			*sort_order = "Asc"
			sort_changed = true
		case "3":
			*sorted_by = "Service Name"
			*sort_order = "Desc"
			sort_changed = true
		case "4":
			*sorted_by = "Login ID"
			*sort_order = "Asc"
			sort_changed = true
		case "5":
			*sorted_by = "Login ID"
			*sort_order = "Desc"
			sort_changed = true
		case "6":
			*sorted_by = "Date Created"
			*sort_order = "Asc"
			sort_changed = true
		case "7":
			*sorted_by = "Date Created"
			*sort_order = "Desc"
			sort_changed = true
		case "8":
			*sorted_by = "Date Modified"
			*sort_order = "Asc"
			sort_changed = true
		case "9":
			*sorted_by = "Date Modified"
			*sort_order = "Desc"
			sort_changed = true
		case "10":
			*sorted_by = "Exact Service Name Search Mode"
			*sort_order = ""
			sort_changed = true
		case "11":
			choice = "11"
		default:
			choice = "0"
		}

		if sort_changed {
			applySort(current_session, *sorted_by, *sort_order)
		}
	}
}

func binarySearch(current_session *Session, target_name string) int {
	var found_idx int = -1
	var lo, hi, mid int
	var mid_name, target_lower string

	if current_session.Vault != nil && current_session.Vault.Service_Count > 0 {
		lo = 0
		hi = current_session.Vault.Service_Count - 1
		target_lower = toLower(target_name)

		for lo <= hi && found_idx == -1 {
			mid = lo + (hi-lo)/2
			mid_name = toLower(current_session.Vault.Services[mid].Service_Name)

			if mid_name == target_lower {
				found_idx = mid
			} else if mid_name > target_lower {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
	}

	return found_idx
}

func sequentialSearch(current_session *Session, target_name string, found_idx *[SMAX]int, found_count *int) {
	*found_count = 0

	if current_session.Vault != nil && target_name != "" {
		i := 0
		for i < current_session.Vault.Service_Count {
			service := current_session.Vault.Services[i]
			date_created := service.Date_Created.Format("2006-01-02 15:04:05")
			date_modified := service.Date_Modified.Format("2006-01-02 15:04:05")
			if containsFoldASCII(service.Service_Name, target_name) || containsFoldASCII(service.Login_ID, target_name) || containsFoldASCII(date_created, target_name) || containsFoldASCII(date_modified, target_name) {
				found_idx[*found_count] = i
				*found_count = *found_count + 1
			}
			i = i + 1
		}
	}
}

func containsFoldASCII(text, target_name string) bool {
	var result bool

	text = toLower(text)
	target_name = toLower(target_name)

	if target_name == "" {
		result = true
	} else if len(target_name) <= len(text) {
		i := 0
		for i <= len(text)-len(target_name) && !result {
			j := 0
			for j < len(target_name) && text[i+j] == target_name[j] {
				j = j + 1
			}
			if j == len(target_name) {
				result = true
			}
			i = i + 1
		}
	}
	return result
}

func toLower(text string) string {
	var result string

	i := 0
	for i < len(text) {
		c := text[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result = result + string(c)
		i = i + 1
	}
	return result
}

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
			if isOptionConfirmed("Welcome Menu >> Exitting", "Proceed to exit the program?") {
				exit(current_session)
			} else {
				choice = "0"
			}
		default:
			choice = "0"
		}
	}
}

func mainMenu(current_session *Session) {
	var choice string = "0"
	for choice == "0" {
		printHeader("Main Menu")
		fmt.Printf("1. Services Saved\n2. Add Service\n3. Log Out\n4. Log Out & Exit")
		readOption(&choice)

		switch choice {
		case "1":
			if current_session.Vault == nil || current_session.Vault.Service_Count == 0 {
				printHeader("Services Saved")
				displayPause("No services saved yet.", "go back to main menu")
			} else {
				displayPagedServices(current_session, "Services Saved", "ID")
			}
			choice = "0"
		case "2":
			addService(current_session)
			choice = "0"
		case "3":
			current_session.A_Index = StateUnauthorized
			current_session.Vault = nil
			printHeader("Main Menu >> Unauthenticating")
			displayPause("Logging out", "log out")
		case "4":
			if isOptionConfirmed("Main Menu >> Exitting", "Proceed to exit the program?") {
				exit(current_session)
			} else {
				choice = "0"
			}
		default:
			choice = "0"
		}
	}
}
