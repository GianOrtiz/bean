package usecase

import (
	"errors"
	"log"
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

func (uc *journalAccountUseCase) Transact(entries []journal.TransactEntry) (err error) {
	sumOfEntries := 0
	var userIDs []int
	for _, entry := range entries {
		sumOfEntries += money.ToCents(entry.Amount)

		for _, userID := range userIDs {
			if entry.UserID == userID {
				return errors.New("repeated user in entries")
			}
		}
		userIDs = append(userIDs, entry.UserID)
	}

	if sumOfEntries != 0 {
		return errors.New("sum of entries must be equal to zero")
	}

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
	transactionID := uuid.NewString()
	for _, entry := range entries {
		log.Printf("%+v\n", entry)
		account, err := uc.findJournalAccountOrCreateFromTransactEntry(entry, now)
		if err != nil {
			return err
		}
		account.Balance = money.Plus(account.Balance, entry.Amount)
		err = uc.journalAccountRepository.Update(account)
		if err != nil {
			return err
		}

		err = uc.journalEntryRepository.Create(transactionID, account.ID, entry.Amount, now)
		if err != nil {
			return err
		}
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

func (uc *journalAccountUseCase) findJournalAccountOrCreateFromTransactEntry(entry journal.TransactEntry, time time.Time) (*journal.Account, error) {
	fromAccount, err := uc.journalAccountRepository.GetByUserID(entry.UserID)
	if err != nil {
		journalAccountID := uuid.NewString()
		journalAccount := journal.Account{
			ID:        journalAccountID,
			CreatedAt: time,
			Balance:   money.FromCents(0),
			UserID:    entry.UserID,
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
