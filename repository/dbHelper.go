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
			//TODO: log error
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

	tx.MustExec(`CREATE VIRTUAL TABLE IF NOT EXISTS note_content USING fts4(content="note", title, memo)`)
	
	tx.MustExec(`CREATE TRIGGER note_bu BEFORE UPDATE ON note BEGIN
				 DELETE FROM note_content WHERE docid = old.rowid;
				 END;`)
	tx.MustExec(`CREATE TRIGGER note_bd BEFORE DELETE ON note BEGIN
				DELETE FROM note_content WHERE docid = old.rowid;
				END`)

	tx.MustExec(`CREATE TRIGGER note_au AFTER UPDATE ON note BEGIN
				 INSERT INTO note_content(docid, title, memo) VALUES(new.rowid, new.title, new.memo);
				 END;`)
	tx.MustExec(`CREATE TRIGGER note_ad AFTER DELETE ON note BEGIN
				INSERT INTO note_content(docid, title, memo) VALUES(new.rowid, new.title, new.memo);
				END;`)


	err = tx.Commit()
	checkError(err)
	return db
}

func checkError(err error) {
	if err != nil {
		//todo logs
		panic(err)
	}
}

func removeDups(integers []int64) []int64{
	seen := make(map[int64]struct{}, len(integers))
    j := 0
    for _, v := range integers {
        if _, ok := seen[v]; ok {
            continue
        }
        seen[v] = struct{}{}
        integers[j] = v
        j++
    }
    return integers[:j]
}
