CREATE TABLE journal_account(
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    balance INTEGER NOT NULL,
    created_at DATETIME NOT NULL,

    FOREIGN KEY (user_id) REFERENCES user(id)
);
