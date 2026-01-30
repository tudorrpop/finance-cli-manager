# Personal Finance CLI Manager

## What I’ve successfully built so far:

* A functional TUI dashboard using tview
* A tab-based layout:
  * Budgets - view, edit, add, and delete budgets
  * Notifications - placeholder for budget alerts and limit warnings
  * Transactions - currently working on CSV/OFX import & processing
* Fully working CRUD dialogs inside the TUI (add, update, delete)
* Smart Tab key behavior:
  * Switches between tabs when no dialogs are open
  * Switches between fields when inside a dialog
* Layout - some things centered nicely, some… well, they resisted :)
* I used mysql, but will switch to sqlite
* I started with the CLI before doing the TUI thing, so I also have working commands like budget add, budget list, etc.

<img width="977" height="288" alt="Captură de ecran din 2025-11-16 la 20 11 59" src="https://github.com/user-attachments/assets/7678b6f9-2ab9-4fd5-9ee0-fd3c4bf704c3" />
<img width="977" height="288" alt="Captură de ecran din 2025-11-16 la 20 12 14" src="https://github.com/user-attachments/assets/3f7e46b1-eda5-4636-a9ae-8c1188e0e1ee" />
<img width="977" height="288" alt="Captură de ecran din 2025-11-16 la 20 12 06" src="https://github.com/user-attachments/assets/3a70a3e0-bc15-4ff4-b8ce-f33a3e855e66" />
<img width="977" height="413" alt="Captură de ecran din 2025-11-16 la 20 12 21" src="https://github.com/user-attachments/assets/a1768ebb-97a6-47bc-8c09-3d3b1d10591f" />

## What I’ve modified in the last part:

* SQLite Migration: Switched the backend to SQLite for easier local setup and zero-configuration storage.
* TUI Overhaul: Improved the terminal interface layout with responsive tabs, data tables, and search.
* CSV Import: Implemented a fully functional CSV parser to batch import transactions.
* Feature Completion: All core requirements (Budgets, Reports, Search and Manual Entry) are now implemented.

<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 41 27" src="https://github.com/user-attachments/assets/0aac9a9f-189e-4d9c-a8e3-e26bb6093d24" />
<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 42 57" src="https://github.com/user-attachments/assets/0cecacfb-650d-4b02-b569-de9e058bd505" />
<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 41 57" src="https://github.com/user-attachments/assets/33455487-80c8-4fcc-8fcc-5051d4e597bf" />
<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 50 02" src="https://github.com/user-attachments/assets/af67f650-5b9f-412c-a937-46a31ad8229b" />
<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 50 09" src="https://github.com/user-attachments/assets/084c3265-d1f7-4167-934d-0bfef4771666" />
<img width="970" height="481" alt="Captură de ecran din 2026-01-31 la 00 50 13" src="https://github.com/user-attachments/assets/fc9f1dee-c545-4736-820d-0cc086b78de5" />


