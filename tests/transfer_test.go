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

	code := m.Run()
	os.Exit(code)
}

func TestSuccessfulTransfer(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 100)
	newBalance, err := mutation.Transfer(context.Background(), senderAddress, recipientAddress, 200)

	assert.NoError(t, err)
	assert.Equal(t, int32(800), newBalance)

	var recipient db.Wallet
	testDB.First(&recipient, "address = ?", recipientAddress)
	assert.Equal(t, int64(300), recipient.Balance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(800), sender.Balance)
}

func TestTransferInsufficientBalance(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"

	err, mutation := SetUpDatabase(t, senderAddress, 100, recipientAddress, 100)
	_, err = mutation.Transfer(context.Background(), senderAddress, recipientAddress, 200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance")
}

func TestTransferToUnknownRecipient(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"
	unknowRecipientAddress := "0x0000000000000000000000000000000000000003"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 100)
	_, err = mutation.Transfer(context.Background(), senderAddress, unknowRecipientAddress, 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recipient not found")
}

func TestTransferFromUnknownSender(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"
	unknowSenderAddress := "0x0000000000000000000000000000000000000003"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 100)
	_, err = mutation.Transfer(context.Background(), unknowSenderAddress, recipientAddress, 100)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sender not found")
}

func SetUpDatabase(t *testing.T, senderAddress string, senderBalance int64, recipientAddress string, recipientBalance int64) (error, graph.MutationResolver) {
	RestartDatabase()

	err := CreateWallet(t, senderAddress, senderBalance)
	assert.NoError(t, err)
	err = CreateWallet(t, recipientAddress, recipientBalance)
	assert.NoError(t, err)

	mutation := CreateMutationResolver()
	return err, mutation
}

func RestartDatabase() *gorm.DB {
	return testDB.Exec("TRUNCATE TABLE wallets RESTART IDENTITY")
}

func CreateMutationResolver() graph.MutationResolver {
	resolver := &graph.Resolver{DB: testDB}
	mutation := resolver.Mutation()
	return mutation
}

func CreateWallet(t *testing.T, senderAddress string, balance int64) error {
	err := testDB.Create(&db.Wallet{Address: senderAddress, Balance: balance}).Error
	assert.NoError(t, err)
	return err
}
