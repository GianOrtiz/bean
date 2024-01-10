package usecase

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/journal"
	mock_journal "github.com/GianOrtiz/bean/pkg/journal/mock"
	"github.com/GianOrtiz/bean/pkg/money"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

const (
	JOURNAL_ACCOUNT_ID         = "b0a4d93d-b940-4ade-a359-7f9f131b3b1b"
	FROM_ACCOUNT_ID            = "e126a63b-7e63-4367-a20f-b485cd8b5561"
	FROM_USER_ID               = 1
	TO_USER_ID                 = 2
	FROM_ACCOUNT_BALANCE_CENTS = 1000
	TO_ACCOUNT_ID              = "97526968-6976-411a-b425-708f8dd826d0"
	TO_ACCOUNT_BALANCE_CENTS   = 2000
	TRANSACTION_AMOUNT_CENTS   = 200
)

func TestShouldFindEntriesOnAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	journalAccountRepository := mock_journal.NewMockAccountRepository(ctrl)
	journalEntryRepository := mock_journal.NewMockEntryRepository(ctrl)
	usecase := journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
	}

	journalAccountRepository.
		EXPECT().
		GetByID(JOURNAL_ACCOUNT_ID).
		Return(&journal.Account{
			ID:        JOURNAL_ACCOUNT_ID,
			CreatedAt: time.Now(),
			Balance:   money.FromCents(1000),
		}, nil)
	journalEntryRepository.
		EXPECT().
		GetByJournalAccountID(JOURNAL_ACCOUNT_ID).
		Return([]*journal.Entry{
			{
				ID:               1,
				CreatedAt:        time.Now(),
				JournalAccountID: JOURNAL_ACCOUNT_ID,
				TransactionID:    uuid.NewString(),
				Amount:           money.FromCents(-1100),
			},
			{
				ID:               2,
				CreatedAt:        time.Now(),
				JournalAccountID: JOURNAL_ACCOUNT_ID,
				TransactionID:    uuid.NewString(),
				Amount:           money.FromCents(2100),
			},
		}, nil)

	entries, err := usecase.FindEntries(JOURNAL_ACCOUNT_ID)
	if err != nil {
		t.Errorf("expected to receive no error on find entries, got %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected to receive 2 entries, received %d", len(entries))
	}
}

func TestShouldThrowErrorIfAccountIsNotFoundWhenFindEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	journalAccountRepository := mock_journal.NewMockAccountRepository(ctrl)
	journalEntryRepository := mock_journal.NewMockEntryRepository(ctrl)
	usecase := journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
	}

	journalAccountRepository.
		EXPECT().
		GetByID(JOURNAL_ACCOUNT_ID).
		Return(nil, sql.ErrNoRows)

	_, err := usecase.FindEntries(JOURNAL_ACCOUNT_ID)
	if err != JournalAccountNotFoundErr {
		t.Errorf("expected to receive error %v, received: %v", JournalAccountNotFoundErr, err)
	}
}

func TestShouldErrorWhenEntriesDoNotSumUpToZeroOnTransact(t *testing.T) {
	entries := []journal.TransactEntry{
		{
			Amount: money.FromCents(200),
			UserID: 1,
		},
		{
			Amount: money.FromCents(0),
			UserID: 2,
		},
	}

	ctrl := gomock.NewController(t)
	journalAccountRepository := mock_journal.NewMockAccountRepository(ctrl)
	journalEntryRepository := mock_journal.NewMockEntryRepository(ctrl)
	dbConn, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	usecase := journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
		db:                       db,
	}

	err = usecase.Transact(entries)
	if err == nil {
		t.Error("expected an error, received nil")
	}
}

func TestShouldErrorWhenUserIsRepeatedInEntriesOnTransact(t *testing.T) {
	entries := []journal.TransactEntry{
		{
			Amount: money.FromCents(200),
			UserID: 1,
		},
		{
			Amount: money.FromCents(-200),
			UserID: 1,
		},
	}

	ctrl := gomock.NewController(t)
	journalAccountRepository := mock_journal.NewMockAccountRepository(ctrl)
	journalEntryRepository := mock_journal.NewMockEntryRepository(ctrl)
	dbConn, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	db := db.NewSqlDB(dbConn)
	usecase := journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
		db:                       db,
	}

	err = usecase.Transact(entries)
	if err == nil {
		t.Error("expected an error, received nil")
	}
}

func TestShouldTransact(t *testing.T) {
	entries := []journal.TransactEntry{
		{
			Amount: money.FromCents(-1 * TRANSACTION_AMOUNT_CENTS),
			UserID: FROM_USER_ID,
		},
		{
			Amount: money.FromCents(TRANSACTION_AMOUNT_CENTS),
			UserID: TO_USER_ID,
		},
	}

	ctrl := gomock.NewController(t)
	journalAccountRepository := mock_journal.NewMockAccountRepository(ctrl)
	journalEntryRepository := mock_journal.NewMockEntryRepository(ctrl)
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not stablish a new connection with database mock: %v", err)
	}
	defer dbConn.Close()

	now := time.Now()
	db := db.NewSqlDB(dbConn)
	usecase := journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
		db:                       db,
	}

	mock.ExpectBegin()
	journalAccountRepository.
		EXPECT().
		EnableTransaction(gomock.Any())
	journalEntryRepository.
		EXPECT().
		EnableTransaction(gomock.Any())

	journalAccountRepository.
		EXPECT().
		GetByUserID(FROM_USER_ID).
		Return(nil, sql.ErrNoRows)
	journalAccountRepository.
		EXPECT().
		Create(gomock.Any())
	journalAccountRepository.
		EXPECT().
		GetByID(gomock.Any()).
		Return(&journal.Account{
			ID:        FROM_ACCOUNT_ID,
			CreatedAt: now,
			Balance:   money.FromCents(FROM_ACCOUNT_BALANCE_CENTS),
		}, nil)

	journalAccountRepository.
		EXPECT().
		GetByUserID(TO_USER_ID).
		Return(nil, sql.ErrNoRows)
	journalAccountRepository.
		EXPECT().
		Create(gomock.Any())
	journalAccountRepository.
		EXPECT().
		GetByID(gomock.Any()).
		Return(&journal.Account{
			ID:        TO_ACCOUNT_ID,
			CreatedAt: now,
			Balance:   money.FromCents(TO_ACCOUNT_BALANCE_CENTS),
		}, nil)

	journalEntryRepository.
		EXPECT().
		Create(gomock.Any(), FROM_ACCOUNT_ID, money.Negative(TRANSACTION_AMOUNT_CENTS), gomock.Any())
	journalEntryRepository.
		EXPECT().
		Create(gomock.Any(), TO_ACCOUNT_ID, money.FromCents(TRANSACTION_AMOUNT_CENTS), gomock.Any())

	journalAccountRepository.
		EXPECT().
		Update(&journal.Account{
			ID:        FROM_ACCOUNT_ID,
			CreatedAt: now,
			Balance: money.Minus(
				money.FromCents(FROM_ACCOUNT_BALANCE_CENTS),
				money.FromCents(TRANSACTION_AMOUNT_CENTS)),
		})
	journalAccountRepository.
		EXPECT().
		Update(&journal.Account{
			ID:        TO_ACCOUNT_ID,
			CreatedAt: now,
			Balance: money.Plus(
				money.FromCents(TO_ACCOUNT_BALANCE_CENTS),
				money.FromCents(TRANSACTION_AMOUNT_CENTS)),
		})
	err = usecase.Transact(entries)
	if err != nil {
		t.Errorf("expected to receive nil error on transact, received: %v", err)
	}
}
