package tests

import (
	"context"
	"fmt"
	"github.com/dominika232323/token-transfer-api/graph"
	"github.com/dominika232323/token-transfer-api/internal/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"testing"
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

func TestTransferWithNegativeAmount(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 1000)
	_, err = mutation.Transfer(context.Background(), senderAddress, recipientAddress, -200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount cannot be negative")
}

func TestTransferWithZeroAmount(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 100)
	newBalance, err := mutation.Transfer(context.Background(), senderAddress, senderAddress, 0)

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), newBalance)

	var recipient db.Wallet
	testDB.First(&recipient, "address = ?", recipientAddress)
	assert.Equal(t, int64(100), recipient.Balance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(1000), sender.Balance)
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
	unknowRecipientAddress := "0x0000000000000000000000000000000000000002"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, "", 0)
	newBalance, err := mutation.Transfer(context.Background(), senderAddress, unknowRecipientAddress, 200)

	assert.NoError(t, err)
	assert.Equal(t, int32(800), newBalance)

	var recipient db.Wallet
	testDB.First(&recipient, "address = ?", unknowRecipientAddress)
	assert.Equal(t, int64(200), recipient.Balance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(800), sender.Balance)
}

func TestTransferFromUnknownSender(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"
	unknowSenderAddress := "0x0000000000000000000000000000000000000003"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, recipientAddress, 100)
	_, err = mutation.Transfer(context.Background(), unknowSenderAddress, recipientAddress, 200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance")
}

func TestTransferToSelf(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, "", 0)
	newBalance, err := mutation.Transfer(context.Background(), senderAddress, senderAddress, 200)

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), newBalance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(1000), sender.Balance)
}

func TestTransferWithNegativeAmountToSelf(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, "", 0)
	_, err = mutation.Transfer(context.Background(), senderAddress, senderAddress, -200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount cannot be negative")
}

func TestTransferWithZeroAmountToSelf(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"

	err, mutation := SetUpDatabase(t, senderAddress, 1000, "", 0)
	newBalance, err := mutation.Transfer(context.Background(), senderAddress, senderAddress, 0)

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), newBalance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(1000), sender.Balance)
}

func TestTransferToSelfWithInsufficientBalance(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"

	err, mutation := SetUpDatabase(t, senderAddress, 100, "", 0)
	_, err = mutation.Transfer(context.Background(), senderAddress, senderAddress, 200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance")
}

func TestTransferToNonExistentSelf(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"

	err, mutation := SetUpDatabase(t, "", 0, "", 0)
	_, err = mutation.Transfer(context.Background(), senderAddress, senderAddress, 200)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance")
}

func TestConcurrentTransfers(t *testing.T) {
	wallet1Address := "0x0000000000000000000000000000000000000001"
	wallet2Address := "0x0000000000000000000000000000000000000002"

	_, mutation := SetUpDatabase(t, wallet1Address, 10, wallet2Address, 10)

	var wg sync.WaitGroup
	wg.Add(3)
	start := make(chan struct{})

	results := make([]error, 3)
	transfers := []int32{-4, -7, 1}

	for i, amount := range transfers {
		go func(i int, amount int32) {
			defer wg.Done()

			<-start

			if amount < 0 {
				_, err := mutation.Transfer(context.Background(), wallet1Address, wallet2Address, -1*amount)
				results[i] = err
			} else {
				_, err := mutation.Transfer(context.Background(), wallet2Address, wallet1Address, amount)
				results[i] = err
			}

		}(i, amount)
	}

	close(start)
	wg.Wait()

	var wallet1 db.Wallet
	testDB.First(&wallet1, "address = ?", wallet1Address)

	var wallet2 db.Wallet
	testDB.First(&wallet2, "address = ?", wallet2Address)

	var wallet1Received int32
	var wallet2Received int32

	for i, err := range results {
		if err == nil {
			if transfers[i] < 0 {
				wallet2Received += -1 * transfers[i]
			} else {
				wallet1Received += transfers[i]
			}
		}
	}

	expectedFinalWallet1Balance := 10 - int64(wallet2Received) + int64(wallet1Received)
	expectedFinalWallet2Balance := 10 - int64(wallet1Received) + int64(wallet2Received)

	assert.Equal(t, expectedFinalWallet1Balance, wallet1.Balance)
	assert.Equal(t, expectedFinalWallet2Balance, wallet2.Balance)

	assert.GreaterOrEqual(t, wallet1.Balance, int64(0))
	assert.GreaterOrEqual(t, wallet2.Balance, int64(0))
}

func TestConcurrentTransfers_MultipleRuns(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Run-%d", i), TestConcurrentTransfers)
	}
}

func TestBidirectionalConcurrentTransfers(t *testing.T) {
	walletA := "0x0000000000000000000000000000000000000001"
	walletB := "0x0000000000000000000000000000000000000002"

	_, mutation := SetUpDatabase(t, walletA, 1000, walletB, 1000)

	var wg sync.WaitGroup
	wg.Add(2)

	start := make(chan struct{})

	var err1, err2 error

	go func() {
		defer wg.Done()
		<-start
		_, err1 = mutation.Transfer(context.Background(), walletA, walletB, 100)
	}()

	go func() {
		defer wg.Done()
		<-start
		_, err2 = mutation.Transfer(context.Background(), walletB, walletA, 150)
	}()

	close(start)
	wg.Wait()

	if err1 != nil {
		assert.NotContains(t, err1.Error(), "deadlock")
	}
	if err2 != nil {
		assert.NotContains(t, err2.Error(), "deadlock")
	}

	var a db.Wallet
	var b db.Wallet

	testDB.First(&a, "address = ?", walletA)
	testDB.First(&b, "address = ?", walletB)

	total := a.Balance + b.Balance

	assert.Equal(t, int64(2000), total, "Total balance should remain constant")
	assert.Equal(t, a.Balance, int64(1050))
	assert.Equal(t, b.Balance, int64(950))
}

func TestConcurrentWalletCreation(t *testing.T) {
	senderAddress := "0x0000000000000000000000000000000000000001"
	recipientAddress := "0x0000000000000000000000000000000000000002"

	_, mutation := SetUpDatabase(t, senderAddress, 1000, "", 0)

	var wg sync.WaitGroup
	numTransfers := 5
	wg.Add(numTransfers)
	start := make(chan struct{})

	transferAmount := int32(100)
	errors := make([]error, numTransfers)

	for i := 0; i < numTransfers; i++ {
		go func(i int) {
			defer wg.Done()

			<-start

			_, err := mutation.Transfer(context.Background(), senderAddress, recipientAddress, transferAmount)
			errors[i] = err
		}(i)
	}

	close(start)
	wg.Wait()

	successfulTransfers := 0

	for _, err := range errors {
		if err == nil {
			successfulTransfers++
		}
	}

	assert.Greater(t, successfulTransfers, 0, "At least one transfer should succeed")

	var recipient db.Wallet
	result := testDB.First(&recipient, "address = ?", recipientAddress)
	assert.NoError(t, result.Error, "New wallet should exist")

	expectedRecipientBalance := int64(successfulTransfers) * int64(transferAmount)
	assert.Equal(t, expectedRecipientBalance, recipient.Balance)

	var sender db.Wallet
	testDB.First(&sender, "address = ?", senderAddress)
	assert.Equal(t, int64(1000)-expectedRecipientBalance, sender.Balance)
}

func SetUpDatabase(t *testing.T, senderAddress string, senderBalance int64, recipientAddress string, recipientBalance int64) (error, graph.MutationResolver) {
	RestartDatabase()

	err := CreateWallet(t, senderAddress, senderBalance)
	assert.NoError(t, err)

	if recipientAddress != "" {
		err = CreateWallet(t, recipientAddress, recipientBalance)
		assert.NoError(t, err)
	}

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
