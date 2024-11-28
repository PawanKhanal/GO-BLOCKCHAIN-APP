package main

import (
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

// Blockchain stores the chain of blocks
type Blockchain struct {
	Chain               []Block
	CurrentTransactions []string
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

// Initialize the blockchain with the genesis block
func initBlockchain() Blockchain {
	blockchain := Blockchain{}
	blockchain.Chain = append(blockchain.Chain, CreateGenesisBlock())                              // Add the genesis block
	blockchain.CurrentTransactions = append(blockchain.CurrentTransactions, "Initial transaction") // Initial transaction
	return blockchain
}

// Handler to return the blockchain data as JSON
func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain.Chain)
}

// Global blockchain variable
var blockchain Blockchain

func main() {
	// Initialize the blockchain with the genesis block
	blockchain = initBlockchain()

	// Example block creation
	newBlock := CreateBlock(blockchain.Chain[len(blockchain.Chain)-1], []string{"New transaction"}, 12345)
	blockchain.Chain = append(blockchain.Chain, newBlock)

	// Print the initialized blockchain and the new block
	fmt.Println("Blockchain initialized with Genesis Block: ", blockchain.Chain)
	fmt.Println("New block added:", newBlock)

	// Set up CORS to allow communication with the frontend
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow React app
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Register the blockchain API endpoint
	http.HandleFunc("/api/blockchain", getBlockchainHandler)

	// Start the HTTP server and wrap it with the CORS handler
	http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux))
}
