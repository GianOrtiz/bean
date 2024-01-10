package usecase

import (
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

func NewJournalAccountUseCase(journalAccountRepository journal.AccountRepository, journalEntryRepository journal.EntryRepository) journal.AccountUseCase {
	return &journalAccountUseCase{
		journalAccountRepository: journalAccountRepository,
		journalEntryRepository:   journalEntryRepository,
	}
}

func (uc *journalAccountUseCase) Transact(fromAccountID string, toAccountID string, amount money.Money) (err error) {
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

	fromAccount, err := uc.findJournalAccountOrCreate(fromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := uc.findJournalAccountOrCreate(toAccountID)
	if err != nil {
		return err
	}

	fromAccount.Balance = money.Minus(fromAccount.Balance, amount)
	toAccount.Balance = money.Plus(toAccount.Balance, amount)

	fromAccountEntry := journal.Entry{
		JournalAccountID: fromAccountID,
		Amount:           money.Negative(amount),
	}
	toAccountEntry := journal.Entry{
		JournalAccountID: toAccountID,
		Amount:           amount,
	}
	transactionID := uuid.NewString()

	err = uc.journalEntryRepository.
		Create(transactionID, fromAccountEntry.JournalAccountID, fromAccountEntry.Amount)
	if err != nil {
		return err
	}

	err = uc.journalEntryRepository.
		Create(transactionID, toAccountEntry.JournalAccountID, toAccountEntry.Amount)
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

func (uc *journalAccountUseCase) findJournalAccountOrCreate(journalAccountID string) (*journal.Account, error) {
	fromAccount, err := uc.journalAccountRepository.GetByID(journalAccountID)
	if err != nil {
		if err := uc.journalAccountRepository.Create(journalAccountID, money.FromCents(0)); err != nil {
			return nil, err
		}
		fromAccount, err = uc.journalAccountRepository.GetByID(journalAccountID)
		if err != nil {
			return nil, err
		}
	}
	return fromAccount, nil
}
