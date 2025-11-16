package tui

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"finance-cli-manager/internal/db"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunTUI() error {
	database, err := db.ConnectDB()
	if err != nil {
		return err
	}
	defer database.Close()

	app := tview.NewApplication()

	tabNames := []string{"Budgets", "Transactions", "Notifications"}
	currentTab := 0

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(`
[::b][green]╔══════════════════════════════════╗
[::b][green]║[yellow] Your (Accredited) Personal Finance Manager [green]║
[::b][green]╚══════════════════════════════════╝
`)

	tabBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	updateTabBar := func() {
		text := ""
		for i, name := range tabNames {
			if i == currentTab {
				text += fmt.Sprintf("[white:blue] %s [-:-]  ", name)
			} else {
				text += fmt.Sprintf("[black:white] %s [-:-]  ", name)
			}
		}
		tabBar.SetText(text)
	}
	updateTabBar()

	legend := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("----------------------------------------------------\n[green]a[-] Add   [yellow]d[-] Delete   [cyan]Enter[-] Update   [red]Esc[-] Back   [magenta]Tab[-] Switch Tab")

	pages := tview.NewPages()

	budgetTable := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	refreshBudgetTable(budgetTable, database)

	if budgetTable.GetRowCount() > 1 {
		budgetTable.Select(1, 0)
	}

	transView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Transactions view coming soon")
	pages.AddPage("Transactions", transView, true, false)

	notifView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Notifications view coming soon")
	pages.AddPage("Notifications", notifView, true, false)

	pages.AddPage("Budgets", budgetTable, true, true)

	var addForm, updateForm *tview.Form
	var deleteModal *tview.Modal
	var updateID, deleteID int

	addForm = tview.NewForm().
		AddInputField("Category", "", 20, nil, nil).
		AddInputField("Amount", "", 20, nil, nil)
	addForm.AddButton("Save", func() {
		cat := addForm.GetFormItemByLabel("Category").(*tview.InputField).GetText()
		amtStr := addForm.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			log.Printf("Invalid amount: %v", err)
			return
		}
		if err := db.AddBudget(database, cat, amt); err != nil {
			log.Printf("AddBudget error: %v", err)
		}
		refreshBudgetTable(budgetTable, database)
		pages.SwitchToPage("Budgets")
		app.SetFocus(budgetTable)
	})
	addForm.AddButton("Cancel", func() {
		pages.SwitchToPage("Budgets")
		app.SetFocus(budgetTable)
	})
	addForm.SetBorder(true).SetTitle("Add Budget").SetTitleAlign(tview.AlignCenter)
	pages.AddPage("addForm", addForm, true, false)

	updateForm = tview.NewForm().
		AddInputField("Category", "", 20, nil, nil).
		AddInputField("Amount", "", 20, nil, nil)
	updateForm.AddButton("Save", func() {
		cat := updateForm.GetFormItemByLabel("Category").(*tview.InputField).GetText()
		amtStr := updateForm.GetFormItemByLabel("Amount").(*tview.InputField).GetText()
		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			log.Printf("Invalid amount: %v", err)
			return
		}
		if err := db.UpdateBudget(database, updateID, cat, amt); err != nil {
			log.Printf("UpdateBudget error: %v", err)
		}
		refreshBudgetTable(budgetTable, database)
		pages.SwitchToPage("Budgets")
		app.SetFocus(budgetTable)
	})
	updateForm.AddButton("Cancel", func() {
		pages.SwitchToPage("Budgets")
		app.SetFocus(budgetTable)
	})
	updateForm.SetBorder(true).SetTitle("Update Budget").SetTitleAlign(tview.AlignCenter)
	pages.AddPage("updateForm", updateForm, true, false)

	deleteModal = tview.NewModal().
		SetText("Are you sure you want to delete this budget?").
		AddButtons([]string{"Cancel", "Delete"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				if err := db.DeleteBudget(database, deleteID); err != nil {
					log.Printf("DeleteBudget error: %v", err)
				}
				refreshBudgetTable(budgetTable, database)
			}
			pages.SwitchToPage("Budgets")
			app.SetFocus(budgetTable)
		})
	pages.AddPage("deleteModal", deleteModal, true, false)

	openUpdateForm := func(row int) {
		if row == 0 {
			return
		}
		idCell := budgetTable.GetCell(row, 0)
		catCell := budgetTable.GetCell(row, 1)
		amtCell := budgetTable.GetCell(row, 2)
		id, _ := strconv.Atoi(idCell.Text)
		updateID = id
		updateForm.GetFormItemByLabel("Category").(*tview.InputField).SetText(catCell.Text)
		updateForm.GetFormItemByLabel("Amount").(*tview.InputField).SetText(amtCell.Text)
		pages.SwitchToPage("updateForm")
		app.SetFocus(updateForm)
	}

	budgetTable.SetSelectedFunc(func(row, column int) {
		openUpdateForm(row)
	})

	budgetTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := budgetTable.GetSelection()
		switch event.Rune() {
		case 'a':
			pages.SwitchToPage("addForm")
			app.SetFocus(addForm)
			return nil
		case 'd':
			if row > 0 {
				idCell := budgetTable.GetCell(row, 0)
				deleteID, _ = strconv.Atoi(idCell.Text)
				pages.ShowPage("deleteModal")
				app.SetFocus(deleteModal)
			}
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			openUpdateForm(row)
			return nil
		}
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(budgetTable)
			return nil
		}
		return event
	})

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 5, 1, false).
		AddItem(tabBar, 1, 1, false).
		AddItem(legend, 2, 1, false).
		AddItem(pages, 0, 10, true)

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		frontPage, _ := pages.GetFrontPage()
		dialogActive := frontPage == "addForm" || frontPage == "updateForm" || frontPage == "deleteModal"

		if dialogActive {
			return event
		}

		switch event.Key() {
		case tcell.KeyTab:
			pages.HidePage(tabNames[currentTab])
			currentTab = (currentTab + 1) % len(tabNames)
			pages.ShowPage(tabNames[currentTab])
			updateTabBar()
			if currentTab == 0 {
				app.SetFocus(budgetTable)
			}
			return nil
		case tcell.KeyBacktab:
			pages.HidePage(tabNames[currentTab])
			currentTab = (currentTab - 1 + len(tabNames)) % len(tabNames)
			pages.ShowPage(tabNames[currentTab])
			updateTabBar()
			if currentTab == 0 {
				app.SetFocus(budgetTable)
			}
			return nil
		}
		return event
	})

	app.SetRoot(mainFlex, true).EnableMouse(true)
	return app.Run()
}

func refreshBudgetTable(table *tview.Table, dbConn *sql.DB) {
	table.Clear()
	table.SetCell(0, 0, tview.NewTableCell("ID").SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Category").SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Amount").SetSelectable(false).SetAlign(tview.AlignCenter))

	budgets, err := db.ListBudgets(dbConn)
	if err != nil {
		log.Printf("Failed to list budgets: %v", err)
		return
	}

	for i, b := range budgets {
		table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", b.ID)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(b.Category).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", b.Amount)).SetAlign(tview.AlignCenter))
	}
}
