package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountType represents the type of account in the ledger
type AccountType string

const (
	AccountTypeUser    AccountType = "USER"
	AccountTypeGoal    AccountType = "GOAL"
	AccountTypeEscrow  AccountType = "ESCROW"
	AccountTypeRevenue AccountType = "REVENUE"
)

// Account represents an account in the double-entry ledger system
type Account struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountType AccountType `gorm:"not null;size:20;index" json:"account_type"`
	EntityID    uuid.UUID   `gorm:"type:uuid;not null;index" json:"entity_id"` // References user, goal, etc.
	Currency    string      `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	CreatedAt   time.Time   `gorm:"not null" json:"created_at"`

	// Relationships
	LedgerEntries []LedgerEntry `gorm:"foreignKey:AccountID;constraint:OnDelete:RESTRICT" json:"ledger_entries,omitempty"`
}

// BeforeCreate sets UUID before creating account
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EntryType represents the type of ledger entry
type EntryType string

const (
	EntryTypeDebit  EntryType = "DEBIT"
	EntryTypeCredit EntryType = "CREDIT"
)

// LedgerEntry represents an immutable entry in the double-entry ledger
type LedgerEntry struct {
	ID            uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID     uuid.UUID              `gorm:"type:uuid;not null;index" json:"account_id"`
	TransactionID uuid.UUID              `gorm:"type:uuid;not null;index" json:"transaction_id"` // Groups related entries
	EntryType     EntryType              `gorm:"not null;size:10" json:"entry_type"`
	Amount        int64                  `gorm:"not null" json:"amount"` // Always positive, type determines debit/credit
	Currency      string                 `gorm:"not null;size:3;default:'NGN'" json:"currency"`
	Description   string                 `gorm:"not null;size:500" json:"description"`
	Metadata      map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt     time.Time              `gorm:"not null;index" json:"created_at"`

	// Relationships
	Account Account `gorm:"constraint:OnDelete:RESTRICT"`
}

// BeforeCreate sets UUID before creating ledger entry
func (le *LedgerEntry) BeforeCreate(tx *gorm.DB) error {
	if le.ID == uuid.Nil {
		le.ID = uuid.New()
	}
	return nil
}

// BalanceSnapshot represents a cached balance calculation for performance
type BalanceSnapshot struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"account_id"`
	Balance   int64     `gorm:"not null" json:"balance"`
	Currency  string    `gorm:"not null;size:3" json:"currency"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	// Relationships
	Account Account `gorm:"constraint:OnDelete:CASCADE"`
}

// BeforeCreate sets UUID before creating balance snapshot
func (bs *BalanceSnapshot) BeforeCreate(tx *gorm.DB) error {
	if bs.ID == uuid.Nil {
		bs.ID = uuid.New()
	}
	return nil
}

// Transaction represents a financial transaction (for grouping ledger entries)
type Transaction struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type        string                 `gorm:"not null;size:50;index" json:"type"` // CONTRIBUTION, WITHDRAWAL, etc.
	Description string                 `gorm:"not null;size:500" json:"description"`
	Amount      int64                  `gorm:"not null" json:"amount"`
	Currency    string                 `gorm:"not null;size:3" json:"currency"`
	Metadata    map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time              `gorm:"not null" json:"created_at"`

	// Relationships
	LedgerEntries []LedgerEntry `gorm:"foreignKey:TransactionID;constraint:OnDelete:RESTRICT" json:"ledger_entries,omitempty"`
}

// BeforeCreate sets UUID before creating transaction
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}