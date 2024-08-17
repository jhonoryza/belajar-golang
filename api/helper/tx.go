package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		err2 := tx.Rollback()
		PanicIfError(err2)
		panic(err)
	} else {
		err3 := tx.Commit()
		PanicIfError(err3)
	}
}
