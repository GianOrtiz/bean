package usecase

import (
	"time"

	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/journal"
	"github.com/GianOrtiz/bean/pkg/money"
	"github.com/google/uuid"
)

type journalAccountUseCase struct {
	journalAccountRepository journal.AccountRepository
	journalEntryRepository   journal.EntryRepository
	db                       db.DBConn
}

func NewJournalAccountUseCase(journalAccountRepository journal.AccountRepository, journalEntryRepository journal.EntryRepository, db db.DBConn) journal.AccountUseCase {
	return &journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
		db:                       db,
	}
}

func (uc *journalAccountUseCase) Transact(fromUserID, toUserID int, amount money.Money) (err error) {
	dbTx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			dbTx.Rollback()
		} else {
			dbTx.Commit()
		}
	}()

	uc.journalAccountRepository.EnableTransaction(dbTx)
	uc.journalEntryRepository.EnableTransaction(dbTx)

	now := time.Now()
	fromAccount, err := uc.findJournalAccountOrCreateByUserID(fromUserID, now)
	if err != nil {
		return err
	}

	toAccount, err := uc.findJournalAccountOrCreateByUserID(toUserID, now)
	if err != nil {
		return err
	}

	fromAccount.Balance = money.Minus(fromAccount.Balance, amount)
	toAccount.Balance = money.Plus(toAccount.Balance, amount)

	fromAccountEntry := journal.Entry{
		JournalAccountID: fromAccount.ID,
		Amount:           money.Negative(amount),
	}
	toAccountEntry := journal.Entry{
		JournalAccountID: toAccount.ID,
		Amount:           amount,
	}
	transactionID := uuid.NewString()

	err = uc.journalEntryRepository.
		Create(transactionID, fromAccountEntry.JournalAccountID, fromAccountEntry.Amount, now)
	if err != nil {
		return err
	}

	err = uc.journalEntryRepository.
		Create(transactionID, toAccountEntry.JournalAccountID, toAccountEntry.Amount, now)
	if err != nil {
		return err
	}

	err = uc.journalAccountRepository.Update(fromAccount)
	if err != nil {
		return err
	}

	err = uc.journalAccountRepository.Update(toAccount)
	if err != nil {
		return err
	}

	return nil
}

func (uc *journalAccountUseCase) FindUserAccount(userID int) (*journal.Account, error) {
	return uc.journalAccountRepository.GetByUserID(userID)
}

func (uc *journalAccountUseCase) FindEntries(journalAccountID string) ([]*journal.Entry, error) {
	_, err := uc.journalAccountRepository.GetByID(journalAccountID)
	if err != nil {
		return nil, JournalAccountNotFoundErr
	}

	entries, err := uc.journalEntryRepository.GetByJournalAccountID(journalAccountID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (uc *journalAccountUseCase) findJournalAccountOrCreateByUserID(userID int, time time.Time) (*journal.Account, error) {
	fromAccount, err := uc.journalAccountRepository.GetByUserID(userID)
	if err != nil {
		journalAccountID := uuid.NewString()
		journalAccount := journal.Account{
			ID:        journalAccountID,
			CreatedAt: time,
			Balance:   money.FromCents(0),
			UserID:    userID,
		}
		if err := uc.journalAccountRepository.Create(journalAccount); err != nil {
			return nil, err
		}
		fromAccount, err = uc.journalAccountRepository.GetByID(journalAccountID)
		if err != nil {
			return nil, err
		}
	}
	return fromAccount, nil
}
