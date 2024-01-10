package db

type TXEnabler interface {
	EnableTransaction(DBTx)
}
