//go:generate mockgen -source=./journal.go -destination=mock/journal.go

package journal

import (
	"time"

	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/money"
)

// Account represents an account in an accountant book.
type Account struct {
	// ID is the unique identifier of the account.
	ID string `json:"id"`
	// CreatedAt is the date the account was created.
	CreatedAt time.Time `json:"createdAt"`
	// Balance is the current balance of the account.
	Balance money.Money `json:"balance"`
}

type Entry struct {
	// ID is the unique identifier of the entry.
	ID int `json:"id"`
	// CreatedAt is the date of creation of the entry.
	CreatedAt time.Time `json:"createdAt"`
	// JournalAccountID is the identifier of the account associated to
	// this entry.
	JournalAccountID string `json:"journalAccountID"`
	// TransactionID is the unique identifier of the transaction between
	// two related entries.
	TransactionID string `json:"transactionId"`
	// Amount is the amount on the transaction of this entry.
	Amount money.Money `json:"money"`
}

// EntryRepository is the abstract representation of data access to
// journal entry.
type EntryRepository interface {
	// Create creates a new journal entry.
	Create(transactionID string, journalAccountID string, amount money.Money) error
	// GetByJournalAccountID retrieves journal entries from a journal account.
	GetByJournalAccountID(journalAccountID string) ([]*Entry, error)
	db.TXEnabler
}

// AccountRepository is the abstract representation of data access
// to journal account.
type AccountRepository interface {
	// Create creates a new journal account.
	Create(id string, initialBalance money.Money) error
	// GetByID retrieves a journal account by its id.
	GetByID(id string) (*Account, error)
	// UpdateBalance updates a journal account.
	Update(account *Account) error
	db.TXEnabler
}

// AccountUseCase is the representation of use cases for the journal
// account data.
type AccountUseCase interface {
	// Transact transacts an amount from one account to another, following
	// a transaction model.
	Transact(fromAccountID string, toAccountID string, amount money.Money) error
	// FindEntries retrieve all journal entries associated to an account.
	FindEntries(journalAccountID string) ([]*Entry, error)
}
