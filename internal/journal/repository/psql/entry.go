package psql

import (
	"time"

	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/journal"
	"github.com/GianOrtiz/bean/pkg/money"
)

type psqlJournalEntryRepository struct {
	conn db.Queryer
}

func NewPSQLJournalEntryRepository(db db.DBConn) journal.EntryRepository {
	return &psqlJournalEntryRepository{conn: db}
}

func (r *psqlJournalEntryRepository) Create(transactionID string, journalAccountID string, amount money.Money, time time.Time) error {
	query := `
		INSERT INTO
			journal_entry(
				transaction_id,
				journal_account_id,
				amount,
				created_at
			)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)
	`
	stmt, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(transactionID, journalAccountID, money.ToCents(amount), time)
	if err != nil {
		return err
	}
	return nil
}

func (r *psqlJournalEntryRepository) GetByJournalAccountID(journalAccountID string) ([]*journal.Entry, error) {
	query := `
		SELECT
			id,
			created_at,
			transaction_id,
			amount,
			journal_account_id
		FROM
			journal_entry
		WHERE
			journal_account_id=$1
	`
	rows, err := r.conn.Query(query, journalAccountID)
	if err != nil {
		return nil, err
	}
	var entries []*journal.Entry
	for rows.Next() {
		var entry journal.Entry
		var moneyCents int
		err = rows.Scan(
			&entry.ID,
			&entry.CreatedAt,
			&entry.TransactionID,
			&moneyCents,
			&entry.JournalAccountID,
		)
		if err != nil {
			return nil, err
		}
		entry.Amount = money.FromCents(moneyCents)
		entries = append(entries, &entry)
	}
	return entries, nil
}

func (r *psqlJournalEntryRepository) EnableTransaction(tx db.DBTx) {
	r.conn = tx
}
