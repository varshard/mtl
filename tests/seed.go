package tests

import "gorm.io/gorm"

func SeedDB(db *gorm.DB) error {
	return db.Exec(`
INSERT INTO vote_item VALUES (1,'item1','description',1),(2,'item 2','Created by John',2),(3,'voted item','it has some votes',1);
`).Error
}

func Truncate(db *gorm.DB) error {
	return db.Exec(`DELETE FROM vote_item;`).Error
}
