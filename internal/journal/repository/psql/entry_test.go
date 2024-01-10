package psql

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/money"
)

const (
	TRANSACTION_ID           = "f126a75e-1f83-48f9-aa28-dc1326961180"
	AMOUNT_TRANSACTION_CENTS = 100
)

func TestShouldCreateJournalEntry(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLJournalEntryRepository(db)

	stmt := mock.ExpectPrepare("INSERT INTO journal_entry")
	stmt.ExpectExec().
		WithArgs(TRANSACTION_ID, JOURNAL_ACCOUNT_ID, AMOUNT_TRANSACTION_CENTS).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repository.Create(
		TRANSACTION_ID,
		JOURNAL_ACCOUNT_ID,
		money.FromCents(AMOUNT_TRANSACTION_CENTS),
	); err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldGetEntriesFromJournalAccount(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLJournalEntryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "created_at", "transaction_id", "amount", "journal_account_id"}).
		AddRow(1, time.Now(), TRANSACTION_ID, AMOUNT_TRANSACTION_CENTS, JOURNAL_ACCOUNT_ID)
	mock.ExpectQuery("SELECT id, created_at, transaction_id, amount, journal_account_id FROM journal_entry").
		WithArgs(JOURNAL_ACCOUNT_ID).
		WillReturnRows(rows)

	journalEntries, err := repository.GetByJournalAccountID(JOURNAL_ACCOUNT_ID)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if len(journalEntries) != 1 {
		t.Errorf(
			"didn't received same journal entries that were expected, received len %d, expected len %d",
			len(journalEntries),
			1,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}
