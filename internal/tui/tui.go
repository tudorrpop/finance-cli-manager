package tui

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"finance-cli-manager/internal/db"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Global variables to manage state
var (
	database    *sql.DB
	app         *tview.Application
	pages       *tview.Pages
	budgetTable *tview.Table
	transTable  *tview.Table
	reportView  *tview.TextView
	searchInput *tview.InputField
)

func RunTUI() error {
	var err error
	database, err = db.ConnectDB()
	if err != nil {
		return err
	}
	defer database.Close()

	app = tview.NewApplication()
	pages = tview.NewPages()

	// -- HEADER & TABS --
	tabNames := []string{"Budgets", "Transactions", "Notifications"}
	currentTab := 0

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(`
[::b][green]╔══════════════════════════════════╗
[::b][green]║[yellow] Your Personal Finance Manager    [green]║
[::b][green]╚══════════════════════════════════╝
`)

	tabBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Function to update the tab display
	updateTabBar := func() {
		text := ""
		for i, name := range tabNames {
			if i == currentTab {
				text += fmt.Sprintf("[white:blue] %s [-:-]  ", name)
			} else {
				text += fmt.Sprintf("[gray] %s [-:-]  ", name)
			}
		}
		tabBar.SetText(text)
	}
	updateTabBar()

	legend := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[yellow]Tab[-] Switch View   [yellow]d[-] Delete   [yellow]n[-] New   [yellow]i[-] Import CSV   [yellow]Esc[-] Back")

	// -- PAGE 1: BUDGETS --
	budgetTable = tview.NewTable().SetBorders(true).SetSelectable(true, false)
	refreshBudgetTable()

	// -- PAGE 2: TRANSACTIONS --
	searchInput = tview.NewInputField().SetLabel("Search: ").SetFieldWidth(30)
	searchInput.SetChangedFunc(func(text string) {
		refreshTransTable(text)
	})

	transTable = tview.NewTable().SetBorders(true).SetSelectable(true, false)
	transFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchInput, 3, 1, false).
		AddItem(transTable, 0, 8, true)

	refreshTransTable("")

	// -- PAGE 3: REPORTS --
	reportView = tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	refreshReport()

	// -- ADD PAGES TO MANAGER --
	pages.AddPage("Budgets", budgetTable, true, true)
	pages.AddPage("Transactions", transFlex, true, false)
	pages.AddPage("Notifications", reportView, true, false)

	// -- KEY INPUTS & NAVIGATION --

	// 1. Budget Table Inputs
	budgetTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'n' {
			showBudgetForm(0, "", 0.0)
			return nil
		} else if event.Rune() == 'd' {
			handleBudgetDelete()
			return nil
		} else if event.Key() == tcell.KeyEnter {
			r, _ := budgetTable.GetSelection()
			if r > 0 {
				id, _ := strconv.Atoi(budgetTable.GetCell(r, 0).Text)
				cat := budgetTable.GetCell(r, 1).Text
				amt, _ := strconv.ParseFloat(budgetTable.GetCell(r, 2).Text, 64)
				showBudgetForm(id, cat, amt)
			}
			return nil
		}
		return event
	})

	// 2. Transaction Page Inputs
	transFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if searchInput.HasFocus() {
			if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEsc {
				app.SetFocus(transTable)
				return nil
			}
			return event
		}

		if event.Rune() == 'i' {
			showImportForm()
			return nil
		} else if event.Rune() == 'n' {
			showTransactionForm()
			return nil
		} else if event.Rune() == 'd' {
			handleTransactionDelete()
			return nil
		} else if event.Rune() == '/' {
			app.SetFocus(searchInput)
			return nil
		}
		return event
	})

	// 3. Global Tab Switching
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		frontPage, _ := pages.GetFrontPage()
		if frontPage == "Form" || frontPage == "Modal" {
			return event
		}

		if event.Key() == tcell.KeyTab {
			pages.HidePage(tabNames[currentTab])
			currentTab = (currentTab + 1) % len(tabNames)
			pages.ShowPage(tabNames[currentTab])
			updateTabBar()

			if currentTab == 0 {
				app.SetFocus(budgetTable)
			} else if currentTab == 1 {
				refreshTransTable("")
				app.SetFocus(transTable)
			} else if currentTab == 2 {
				refreshReport()
			}
			return nil
		}
		return event
	})

	// -- LAYOUT --
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(tabBar, 2, 1, false).
		AddItem(pages, 0, 10, true).
		AddItem(legend, 1, 1, false)

	app.SetRoot(mainFlex, true).EnableMouse(true)
	return app.Run()
}

// ---------------------------
// HELPER FUNCTIONS & FORMS
// ---------------------------

func refreshBudgetTable() {
	budgetTable.Clear()
	headers := []string{"ID", "Category", "Monthly Limit"}
	for i, h := range headers {
		budgetTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}

	budgets, _ := db.ListBudgets(database)
	for i, b := range budgets {
		row := i + 1
		budgetTable.SetCell(row, 0, tview.NewTableCell(strconv.Itoa(b.ID)).SetAlign(tview.AlignCenter))
		budgetTable.SetCell(row, 1, tview.NewTableCell(b.Category).SetAlign(tview.AlignCenter))
		budgetTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.2f", b.Amount)).SetAlign(tview.AlignCenter))
	}
}

func refreshTransTable(search string) {
	transTable.Clear()
	headers := []string{"ID", "Date", "Payee", "Category", "Amount"}
	for i, h := range headers {
		transTable.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}

	trans, _ := db.ListTransactions(database, search)
	for i, t := range trans {
		row := i + 1
		transTable.SetCell(row, 0, tview.NewTableCell(strconv.Itoa(t.ID)))
		transTable.SetCell(row, 1, tview.NewTableCell(t.Date))
		transTable.SetCell(row, 2, tview.NewTableCell(t.Payee))
		transTable.SetCell(row, 3, tview.NewTableCell(t.Category))

		color := tcell.ColorWhite
		if t.Amount < 0 {
			color = tcell.ColorRed
		}
		transTable.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%.2f", t.Amount)).SetTextColor(color).SetAlign(tview.AlignRight))
	}
}

func refreshReport() {
	reportView.Clear()
	reports, _ := db.GetBudgetReport(database)

	fmt.Fprintln(reportView, "\n [::b]MONTHLY PERFORMANCE (Actual vs Budget)[::-]\n")

	for _, r := range reports {
		percent := 0.0
		if r.Limit > 0 {
			percent = (math.Abs(r.Spent) / r.Limit) * 100
		}

		color := "[green]"
		if percent > 80 {
			color = "[yellow]"
		}
		if percent > 100 {
			color = "[red]"
		}

		barWidth := 30
		filled := int((percent / 100) * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

		fmt.Fprintf(reportView, " %s%-15s %s [white]$%.0f / $%.0f (%.0f%%)\n\n",
			color, r.Category, bar, r.Spent, r.Limit, percent)
	}
}

// -- FORMS --

func showBudgetForm(id int, cat string, amt float64) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Manage Budget ").SetTitleAlign(tview.AlignCenter)

	form.AddInputField("Category", cat, 20, nil, nil)
	form.AddInputField("Amount", fmt.Sprintf("%.2f", amt), 20, nil, nil)

	form.AddButton("Save", func() {
		newCat := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
		newAmtStr := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		newAmt, _ := strconv.ParseFloat(newAmtStr, 64)

		if id == 0 {
			db.AddBudget(database, newCat, newAmt)
		} else {
			db.UpdateBudget(database, id, newCat, newAmt)
		}

		pages.RemovePage("Form")
		refreshBudgetTable()
		app.SetFocus(budgetTable)
	})

	form.AddButton("Cancel", func() {
		pages.RemovePage("Form")
		app.SetFocus(budgetTable)
	})

	pages.AddPage("Form", modalCenter(form, 40, 10), true, true)
}

func showTransactionForm() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" New Transaction ").SetTitleAlign(tview.AlignCenter)

	today := time.Now().Format("2006-01-02")
	form.AddInputField("Date (YYYY-MM-DD)", today, 20, nil, nil)
	form.AddInputField("Payee", "", 20, nil, nil)
	form.AddInputField("Category", "Uncategorized", 20, nil, nil)
	form.AddInputField("Amount", "", 20, nil, nil)

	form.AddButton("Save", func() {
		date := form.GetFormItemByLabel("Date (YYYY-MM-DD)").(*tview.InputField).GetText()
		payee := form.GetFormItemByLabel("Payee").(*tview.InputField).GetText()
		cat := form.GetFormItemByLabel("Category").(*tview.InputField).GetText()
		amtStr := form.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amt, _ := strconv.ParseFloat(amtStr, 64)

		if err := db.AddTransaction(database, date, payee, cat, amt); err != nil {
			log.Printf("Error adding transaction: %v", err)
		}

		pages.RemovePage("Form")
		refreshTransTable("")
		app.SetFocus(transTable)
	})

	form.AddButton("Cancel", func() {
		pages.RemovePage("Form")
		app.SetFocus(transTable)
	})

	pages.AddPage("Form", modalCenter(form, 50, 15), true, true)
}

func showImportForm() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Import CSV ").SetTitleAlign(tview.AlignCenter)

	form.AddInputField("File Path", "statement.csv", 35, nil, nil)

	form.AddButton("Import", func() {
		path := form.GetFormItemByLabel("File Path").(*tview.InputField).GetText()
		count, skipped, _, err := db.ProcessCSVTransactions(database, path)

		msg := fmt.Sprintf("Imported %d transactions!", count)
		if err != nil {
			msg = fmt.Sprintf("Error: %v", err)
		} else if len(skipped) > 0 {
			msg += fmt.Sprintf("\nSkipped auto-cat for: %v", skipped[:min(3, len(skipped))])
		}

		modal := tview.NewModal().
			SetText(msg).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(i int, l string) {
				pages.RemovePage("Modal")
				pages.RemovePage("Form")
				refreshTransTable("")
				app.SetFocus(transTable)
			})
		pages.AddPage("Modal", modal, true, true)
	})

	form.AddButton("Cancel", func() {
		pages.RemovePage("Form")
		app.SetFocus(transTable)
	})

	pages.AddPage("Form", modalCenter(form, 50, 9), true, true)
}

func handleBudgetDelete() {
	row, _ := budgetTable.GetSelection()
	if row <= 0 {
		return
	}

	idStr := budgetTable.GetCell(row, 0).Text
	catName := budgetTable.GetCell(row, 1).Text
	id, _ := strconv.Atoi(idStr)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Delete budget for %s?", catName)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				db.DeleteBudget(database, id)
				refreshBudgetTable()
			}
			pages.RemovePage("Modal")
			app.SetFocus(budgetTable)
		})

	pages.AddPage("Modal", modal, true, true)
}

func modalCenter(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func handleTransactionDelete() {
	row, _ := transTable.GetSelection()
	if row <= 0 {
		return
	}

	idStr := transTable.GetCell(row, 0).Text
	payee := transTable.GetCell(row, 2).Text
	id, _ := strconv.Atoi(idStr)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Delete transaction for %s?", payee)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				if err := db.DeleteTransaction(database, id); err != nil {
					log.Printf("Error deleting transaction: %v", err)
				}
				refreshTransTable(searchInput.GetText())
				refreshReport()
			}
			pages.RemovePage("Modal")
			app.SetFocus(transTable)
		})

	pages.AddPage("Modal", modal, true, true)
}
