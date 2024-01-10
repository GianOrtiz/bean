CREATE TABLE journal_entry(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id TEXT NOT NULL,
    journal_account_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at DATETIME NOT NULL,

    FOREIGN KEY (journal_account_id) REFERENCES journal_account(id)
);
