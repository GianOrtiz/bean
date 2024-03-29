package psql

import (
	"github.com/GianOrtiz/bean/pkg/db"
	"github.com/GianOrtiz/bean/pkg/journal"
	"github.com/GianOrtiz/bean/pkg/money"
)

type psqlJournalAccountRepository struct {
	conn db.Queryer
}

func NewPSQLJournalAccountRepository(db db.DBConn) journal.AccountRepository {
	return &psqlJournalAccountRepository{conn: db}
}

func (r *psqlJournalAccountRepository) Create(account journal.Account) error {
	query := `
		INSERT INTO
			journal_account(
				id,
				balance,
				created_at,
				user_id
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
	_, err = stmt.Exec(account.ID, money.ToCents(account.Balance), account.CreatedAt, account.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *psqlJournalAccountRepository) GetByID(id string) (*journal.Account, error) {
	query := `
		SELECT
			id,
			created_at,
			balance
		FROM
			journal_account
		WHERE
			id=$1
	`
	row := r.conn.QueryRow(query, id)
	var account journal.Account
	var moneyCents int
	err := row.Scan(
		&account.ID,
		&account.CreatedAt,
		&moneyCents,
	)
	account.Balance = money.FromCents(moneyCents)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *psqlJournalAccountRepository) GetByUserID(userID int) (*journal.Account, error) {
	query := `
		SELECT
			id,
			created_at,
			balance
		FROM
			journal_account
		WHERE
			user_id=$1
	`
	row := r.conn.QueryRow(query, userID)
	var account journal.Account
	var moneyCents int
	err := row.Scan(
		&account.ID,
		&account.CreatedAt,
		&moneyCents,
	)
	account.Balance = money.FromCents(moneyCents)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *psqlJournalAccountRepository) Update(account *journal.Account) error {
	query := `
		UPDATE
			journal_account
		SET
			balance=$1
		WHERE
			id=$2
	`
	stmt, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(money.ToCents(account.Balance), account.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *psqlJournalAccountRepository) EnableTransaction(tx db.DBTx) {
	r.conn = tx
}
