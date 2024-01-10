package psql

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/journal"
	"github.com/GianOrtiz/bean/pkg/money"
)

const (
	JOURNAL_ACCOUNT_ID    = "52ae3d3f-9b15-457b-a0a1-6474930097cd"
	INITIAL_BALANCE_CENTS = 0
)

func TestShouldCreateNewJournalAccount(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLJournalAccountRepository(db)

	stmt := mock.ExpectPrepare("INSERT INTO journal_account")
	now := time.Now()
	stmt.ExpectExec().WithArgs(JOURNAL_ACCOUNT_ID, INITIAL_BALANCE_CENTS, now, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repository.Create(journal.Account{
		ID:        JOURNAL_ACCOUNT_ID,
		Balance:   money.FromCents(INITIAL_BALANCE_CENTS),
		CreatedAt: now,
		UserID:    1,
	}); err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldRetrieveJournalAccountByID(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLJournalAccountRepository(db)

	rows := sqlmock.NewRows([]string{"id", "created_at", "balance"}).
		AddRow(JOURNAL_ACCOUNT_ID, time.Now(), INITIAL_BALANCE_CENTS)
	mock.ExpectQuery("SELECT id, created_at, balance FROM journal_account").
		WithArgs(JOURNAL_ACCOUNT_ID).
		WillReturnRows(rows)

	journalAccount, err := repository.GetByID(JOURNAL_ACCOUNT_ID)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
	}

	if journalAccount.ID != JOURNAL_ACCOUNT_ID {
		t.Errorf("didn't received same journal account request, id %q do not match requested %q",
			journalAccount.ID,
			JOURNAL_ACCOUNT_ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}

func TestShouldUpdateJournalAccountData(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	repository := NewPSQLJournalAccountRepository(db)

	updatedBalanceJournalAccountCents := 1000

	mock.ExpectPrepare("UPDATE journal_account").
		ExpectExec().
		WithArgs(updatedBalanceJournalAccountCents, JOURNAL_ACCOUNT_ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Update(&journal.Account{
		ID:        JOURNAL_ACCOUNT_ID,
		CreatedAt: time.Now(),
		Balance:   money.FromCents(updatedBalanceJournalAccountCents),
	})
	if err != nil {
		t.Errorf("expected no error on udpate, received: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations of executed statements were not met: %v", err)
	}
}
