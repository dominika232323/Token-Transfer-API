package db

type Wallet struct {
	ID      int64  `gorm:"primaryKey;autoIncrement"`
	Address string `gorm:"uniqueIndex;size:42;not null"`
	Balance int64  `gorm:"not null"`
}
