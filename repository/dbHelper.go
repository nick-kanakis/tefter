package repository

import "github.com/jmoiron/sqlx"

const databaseDriver = "sqlite3"

func connect2DB(dbPath string) *sqlx.DB {
	var err error
	db := sqlx.MustConnect(databaseDriver, dbPath)
	tx := db.MustBegin()

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			//todo: log error
			tx.Rollback()
			db = nil
		}
	}()

	tx.MustExec(`CREATE TABLE IF NOT EXISTS notebook (
		id INTEGER NOT NULL,
		title TEXT NOT NULL,
		CONSTRAINT notebook_PK PRIMARY KEY(id))`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS note (
		id INTEGER NOT NULL,
		title TEXT NOT NULL,
		memo TEXT NOT NULL,
		created DATETIME NOT NULL,
		lastUpdated DATETIME NOT NULL,
		notebook_id INTEGER NOT NULL,
		CONSTRAINT note_PK PRIMARY KEY(id),
		CONSTRAINT notebook_id_FK FOREIGN KEY(notebook_id) REFERENCES notebook(id))`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS note_tag (
			note_id INTEGER NOT NULL,
			tag		TEXT NOT NULL,
			CONSTRAINT note_tag_PK PRIMARY KEY(tag, note_id),
			CONSTRAINT note_id_FK FOREIGN KEY(note_id) REFERENCES note(id))`)

	tx.MustExec(`CREATE VIRTUAL TABLE IF NOT EXISTS note_content USING fts4(title, memo)`)

	err = tx.Commit()
	if err != nil {
		//todo logs
		panic(err)
	}
	return db
}
