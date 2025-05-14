package tests

import (
	"context"
	"github.com/dominika232323/token-transfer-api/graph"
	"github.com/dominika232323/token-transfer-api/internal/db"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env")

	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	testDB.Exec("TRUNCATE TABLE wallets RESTART IDENTITY")

	code := m.Run()
	os.Exit(code)
}

func TestSuccessfulTransfer(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	receiverAddress := "0x0000000000000000000000000000000000000002"

	err := testDB.Create(&db.Wallet{Address: senderAddress, Balance: 1000}).Error
	assert.NoError(t, err)
	err = testDB.Create(&db.Wallet{Address: receiverAddress, Balance: 100}).Error
	assert.NoError(t, err)

	resolver := &graph.Resolver{DB: testDB}
	mutation := resolver.Mutation()

	newBalance, err := mutation.Transfer(context.Background(), senderAddress, receiverAddress, 200)

	assert.NoError(t, err)
	assert.Equal(t, int32(800), newBalance)

	var receiver db.Wallet
	testDB.First(&receiver, "address = ?", receiverAddress)
	assert.Equal(t, int64(300), receiver.Balance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(800), sender.Balance)
}
