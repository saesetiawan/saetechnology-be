package helpers

import (
	"database/sql"
	"fmt"
)

func CommitOrRollback(tx *sql.Tx) {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		fmt.Println("error", p.(error).Error())
	} else {
		_ = tx.Commit()
	}
}
