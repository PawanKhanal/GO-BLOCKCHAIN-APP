package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/cors"
)

// Block defines the structure of a block in the blockchain
type Block struct {
	Index        int      `json:"index"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
	Proof        int      `json:"proof"`
	PrevHash     string   `json:"previous_hash"`
}

// Blockchain stores the chain of blocks and the transaction pool
type Blockchain struct {
	Chain               []Block
	CurrentTransactions []string
	TransactionPool     []string // Added TransactionPool to store pending transactions
}

// Transaction defines the structure of a transaction
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}

// CalculateHash generates a SHA-256 hash of a block
func (b *Block) CalculateHash() string {
	record := strconv.Itoa(b.Index) + b.Timestamp + fmt.Sprint(b.Transactions) + strconv.Itoa(b.Proof) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// CreateGenesisBlock creates the first block in the blockchain
func CreateGenesisBlock() Block {
	return Block{Index: 0, Timestamp: time.Now().String(), Transactions: nil, Proof: 100, PrevHash: ""}
}

// CreateBlock creates a new block based on the previous block
func CreateBlock(prevBlock Block, transactions []string, proof int) Block {
	block := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		Proof:        proof,
		PrevHash:     prevBlock.CalculateHash(),
	}
	return block
}

// ProofOfWork implements the proof-of-work algorithm
func (bc *Blockchain) ProofOfWork(lastProof int) int {
	proof := 0
	for !bc.IsValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// IsValidProof checks if the proof is valid
func (bc *Blockchain) IsValidProof(lastProof, proof int) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.New()
	guessHash.Write([]byte(guess))
	guessResult := guessHash.Sum(nil)
	// Valid proof must result in a hash with leading zeros (a very simple condition)
	return bytes.Equal(guessResult[:4], []byte{0, 0, 0, 0}) // Use bytes.Equal for slice comparison
}

// AddTransactionToPool adds a transaction to the transaction pool
func (bc *Blockchain) AddTransactionToPool(transaction string) {
	bc.TransactionPool = append(bc.TransactionPool, transaction)
}

// ClearTransactionPool clears the transaction pool after mining a block
func (bc *Blockchain) ClearTransactionPool() {
	bc.TransactionPool = []string{}
}

// Initialize the blockchain with the genesis block
func initBlockchain() Blockchain {
	blockchain := Blockchain{}
	blockchain.Chain = append(blockchain.Chain, CreateGenesisBlock()) // Add the genesis block
	return blockchain
}

// Handler to return the blockchain data as JSON
func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain.Chain)
}

// Handler to add a new transaction to the pool
func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid transaction format", http.StatusBadRequest)
		return
	}

	// Add the transaction to the pool
	transactionDetails := fmt.Sprintf("Sender: %s, Receiver: %s, Amount: %d", transaction.Sender, transaction.Receiver, transaction.Amount)
	blockchain.AddTransactionToPool(transactionDetails)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transaction added to the pool")
}

// Handler to mine a new block
func mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	lastProof := blockchain.Chain[len(blockchain.Chain)-1].Proof
	proof := blockchain.ProofOfWork(lastProof)

	// Add the mined block to the blockchain
	newBlock := CreateBlock(blockchain.Chain[len(blockchain.Chain)-1], blockchain.TransactionPool, proof)
	blockchain.Chain = append(blockchain.Chain, newBlock)

	// Clear the transaction pool after mining the block
	blockchain.ClearTransactionPool()

	// Respond with the new block
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBlock)
}

// Global blockchain variable
var blockchain Blockchain

func main() {
	// Initialize the blockchain with the genesis block
	blockchain = initBlockchain()

	// Set up CORS to allow communication with the frontend
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow React app
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Register the blockchain API endpoints
	http.HandleFunc("/api/blockchain", getBlockchainHandler)
	http.HandleFunc("/api/transaction", addTransactionHandler) // Endpoint to add a transaction
	http.HandleFunc("/api/mine", mineBlockHandler)             // Endpoint to mine a new block

	// Start the HTTP server and wrap it with the CORS handler
	http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux))
}
