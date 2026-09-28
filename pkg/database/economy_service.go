package database

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

var economyMu sync.Mutex

func checkedAdd(current, delta int64) (int64, bool) {
	if delta > 0 && current > math.MaxInt64-delta {
		return 0, false
	}
	if delta < 0 && current < math.MinInt64-delta {
		return 0, false
	}
	return current + delta, true
}

var (
	ErrEconomyManagerNotInitialized = errors.New("economy data managers not initialized")
	ErrInsufficientFunds            = errors.New("insufficient funds")
	ErrBankFull                     = errors.New("bank capacity exceeded")
)

const (
	DefaultGlobalBankCapacity int64 = 1000
	DefaultLocalBankCapacity  int64 = 1000
)

// GetGlobalProfile retrieves a user's global economy profile (Stars)
func GetGlobalProfile(userID string) (*models.GlobalEconomyProfile, error) {
	if GlobalEconomyDM == nil {
		return nil, ErrEconomyManagerNotInitialized
	}

	query := bson.M{"_id": userID}
	profile, err := GlobalEconomyDM.Get(query)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		// Initialize a new profile
		profile = &models.GlobalEconomyProfile{
			UserID:       userID,
			StarsWallet:  0,
			StarsBank:    0,
			BankCapacity: DefaultGlobalBankCapacity,
			Inventory:    make(map[string]int),
			Cooldowns:    make(map[string]time.Time),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err = GlobalEconomyDM.Set(query, profile)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}

// GetLocalProfile retrieves a user's local economy profile (Server Coins)
func GetLocalProfile(guildID, userID string) (*models.LocalEconomyProfile, error) {
	if LocalEconomyDM == nil {
		return nil, ErrEconomyManagerNotInitialized
	}

	id := fmt.Sprintf("%s_%s", guildID, userID)
	query := bson.M{"_id": id}
	profile, err := LocalEconomyDM.Get(query)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		// Initialize a new profile
		profile = &models.LocalEconomyProfile{
			ID:           id,
			GuildID:      guildID,
			UserID:       userID,
			Wallet:       0,
			Bank:         0,
			BankCapacity: DefaultLocalBankCapacity,
			Inventory:    make(map[string]int),
			Cooldowns:    make(map[string]time.Time),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err = LocalEconomyDM.Set(query, profile)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}

// AddStars adds or removes Stars from a user's global profile
func AddStars(userID string, amount int64, toBank bool) (*models.GlobalEconomyProfile, error) {
	economyMu.Lock()
	defer economyMu.Unlock()
	return addStars(userID, amount, toBank)
}

func addStars(userID string, amount int64, toBank bool) (*models.GlobalEconomyProfile, error) {
	profile, err := GetGlobalProfile(userID)
	if err != nil {
		return nil, err
	}

	if toBank {
		newBank, ok := checkedAdd(profile.StarsBank, amount)
		if !ok {
			return nil, errors.New("balance overflow")
		}
		if amount > 0 && newBank > profile.BankCapacity {
			return nil, ErrBankFull
		}
		if newBank < 0 {
			return nil, ErrInsufficientFunds
		}
		profile.StarsBank = newBank
	} else {
		newWallet, ok := checkedAdd(profile.StarsWallet, amount)
		if !ok {
			return nil, errors.New("balance overflow")
		}
		if newWallet < 0 {
			return nil, ErrInsufficientFunds
		}
		profile.StarsWallet = newWallet
	}

	profile.UpdatedAt = time.Now()
	_, err = GlobalEconomyDM.Set(bson.M{"_id": userID}, profile)
	return profile, err
}

// AddLocalBalance adds or removes local balance from a user
func AddLocalBalance(guildID, userID string, amount int64, toBank bool) (*models.LocalEconomyProfile, error) {
	economyMu.Lock()
	defer economyMu.Unlock()
	return addLocalBalance(guildID, userID, amount, toBank)
}

func addLocalBalance(guildID, userID string, amount int64, toBank bool) (*models.LocalEconomyProfile, error) {
	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return nil, err
	}

	if toBank {
		newBank, ok := checkedAdd(profile.Bank, amount)
		if !ok {
			return nil, errors.New("balance overflow")
		}
		if amount > 0 && newBank > profile.BankCapacity {
			return nil, ErrBankFull
		}
		if newBank < 0 {
			return nil, ErrInsufficientFunds
		}
		profile.Bank = newBank
	} else {
		newWallet, ok := checkedAdd(profile.Wallet, amount)
		if !ok {
			return nil, errors.New("balance overflow")
		}
		if newWallet < 0 {
			return nil, ErrInsufficientFunds
		}
		profile.Wallet = newWallet
	}

	profile.UpdatedAt = time.Now()
	_, err = LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	return profile, err
}

// TransferStars transfers Stars between users globally
func TransferStars(fromUserID, toUserID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	economyMu.Lock()
	defer economyMu.Unlock()

	_, err := addStars(fromUserID, -amount, false)
	if err != nil {
		return err
	}

	_, err = addStars(toUserID, amount, false)
	if err != nil {
		// Rollback on failure
		_, _ = addStars(fromUserID, amount, false)
		return err
	}

	return nil
}

// TransferLocalBalance transfers local currency between users
func TransferLocalBalance(guildID, fromUserID, toUserID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	economyMu.Lock()
	defer economyMu.Unlock()

	_, err := addLocalBalance(guildID, fromUserID, -amount, false)
	if err != nil {
		return err
	}

	_, err = addLocalBalance(guildID, toUserID, amount, false)
	if err != nil {
		// Rollback on failure
		_, _ = addLocalBalance(guildID, fromUserID, amount, false)
		return err
	}

	return nil
}

// DepositLocal deposits money from Wallet to Bank
func DepositLocal(guildID, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return err
	}

	if profile.Wallet < amount {
		return ErrInsufficientFunds
	}
	if profile.Bank > profile.BankCapacity-amount {
		return ErrBankFull
	}

	profile.Wallet -= amount
	profile.Bank += amount
	profile.UpdatedAt = time.Now()
	_, err = LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	return err
}

// WithdrawLocal withdraws money from Bank to Wallet
func WithdrawLocal(guildID, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return err
	}

	if profile.Bank < amount {
		return ErrInsufficientFunds
	}

	profile.Bank -= amount
	profile.Wallet += amount
	profile.UpdatedAt = time.Now()
	_, err = LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	return err
}

// DepositStars deposits global Stars from Wallet to Bank
func DepositStars(userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetGlobalProfile(userID)
	if err != nil {
		return err
	}

	if profile.StarsWallet < amount {
		return ErrInsufficientFunds
	}
	if profile.StarsBank > profile.BankCapacity-amount {
		return ErrBankFull
	}

	profile.StarsWallet -= amount
	profile.StarsBank += amount
	profile.UpdatedAt = time.Now()
	_, err = GlobalEconomyDM.Set(bson.M{"_id": userID}, profile)
	return err
}

// WithdrawStars withdraws global Stars from Bank to Wallet
func WithdrawStars(userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetGlobalProfile(userID)
	if err != nil {
		return err
	}

	if profile.StarsBank < amount {
		return ErrInsufficientFunds
	}

	profile.StarsBank -= amount
	profile.StarsWallet += amount
	profile.UpdatedAt = time.Now()
	_, err = GlobalEconomyDM.Set(bson.M{"_id": userID}, profile)
	return err
}

// Cooldown checks if a cooldown has expired. Returns (true, 0) if ready, or (false, timeRemaining) if on cooldown.
func CooldownLocal(guildID, userID, command string, duration time.Duration) (bool, time.Duration, error) {
	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return false, 0, err
	}

	if profile.Cooldowns == nil {
		profile.Cooldowns = make(map[string]time.Time)
	}

	lastTime, exists := profile.Cooldowns[command]
	if exists {
		remaining := time.Until(lastTime.Add(duration))
		if remaining > 0 {
			return false, remaining, nil
		}
	}

	return true, 0, nil
}

// SetCooldownLocal updates the cooldown timestamp for a command
func SetCooldownLocal(guildID, userID, command string) error {
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return err
	}
	if profile.Cooldowns == nil {
		profile.Cooldowns = make(map[string]time.Time)
	}
	profile.Cooldowns[command] = time.Now()
	_, err = LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	return err
}

// CooldownStars checks if a global cooldown has expired.
func CooldownStars(userID, command string, duration time.Duration) (bool, time.Duration, error) {
	profile, err := GetGlobalProfile(userID)
	if err != nil {
		return false, 0, err
	}

	if profile.Cooldowns == nil {
		profile.Cooldowns = make(map[string]time.Time)
	}

	lastTime, exists := profile.Cooldowns[command]
	if exists {
		remaining := time.Until(lastTime.Add(duration))
		if remaining > 0 {
			return false, remaining, nil
		}
	}

	return true, 0, nil
}

// SetCooldownStars updates the cooldown timestamp for a global command
func SetCooldownStars(userID, command string) error {
	economyMu.Lock()
	defer economyMu.Unlock()
	profile, err := GetGlobalProfile(userID)
	if err != nil {
		return err
	}
	if profile.Cooldowns == nil {
		profile.Cooldowns = make(map[string]time.Time)
	}
	profile.Cooldowns[command] = time.Now()
	_, err = GlobalEconomyDM.Set(bson.M{"_id": userID}, profile)
	return err
}

// PurchaseItem deducts the balance and updates inventory in one serialized profile write.
func PurchaseItem(guildID, userID string, item models.Item, quantity int64) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if item.Price < 0 || (item.Price > 0 && quantity > math.MaxInt64/item.Price) {
		return errors.New("invalid item price or quantity")
	}

	economyMu.Lock()
	defer economyMu.Unlock()

	totalCost := item.Price * quantity
	if item.IsGlobal {
		profile, err := GetGlobalProfile(userID)
		if err != nil {
			return err
		}
		if profile.StarsWallet < totalCost {
			return ErrInsufficientFunds
		}
		if profile.Inventory == nil {
			profile.Inventory = make(map[string]int)
		}
		maxInt := int64(^uint(0) >> 1)
		if quantity > maxInt-int64(profile.Inventory[item.ID]) {
			return errors.New("inventory quantity overflow")
		}
		profile.StarsWallet -= totalCost
		profile.Inventory[item.ID] += int(quantity)
		profile.UpdatedAt = time.Now()
		_, err = GlobalEconomyDM.Set(bson.M{"_id": userID}, profile)
		return err
	}

	profile, err := GetLocalProfile(guildID, userID)
	if err != nil {
		return err
	}
	if profile.Wallet < totalCost {
		return ErrInsufficientFunds
	}
	if profile.Inventory == nil {
		profile.Inventory = make(map[string]int)
	}
	maxInt := int64(^uint(0) >> 1)
	if quantity > maxInt-int64(profile.Inventory[item.ID]) {
		return errors.New("inventory quantity overflow")
	}
	profile.Wallet -= totalCost
	profile.Inventory[item.ID] += int(quantity)
	profile.UpdatedAt = time.Now()
	_, err = LocalEconomyDM.Set(bson.M{"_id": profile.ID}, profile)
	return err
}

// GetItems returns all items (both global and for the specific guild)
func GetItems(guildID string) ([]models.Item, error) {
	if ItemDM == nil {
		return nil, ErrEconomyManagerNotInitialized
	}

	// Get global items
	queryGlobal := bson.M{"is_global": true}
	globals, err := ItemDM.GetAll(queryGlobal)
	if err != nil {
		return nil, err
	}

	var items []models.Item
	for _, doc := range globals {
		items = append(items, *doc)
	}

	// Get local items if guildID is provided
	if guildID != "" {
		queryLocal := bson.M{"guild_id": guildID, "is_global": false}
		locals, err := ItemDM.GetAll(queryLocal)
		if err != nil {
			return items, err
		}
		for _, doc := range locals {
			items = append(items, *doc)
		}
	}

	return items, nil
}

// SaveItem saves an item (global if GuildID is empty, local otherwise)
func SaveItem(item models.Item) error {
	if ItemDM == nil {
		return ErrEconomyManagerNotInitialized
	}
	_, err := ItemDM.Set(bson.M{"_id": item.ID}, &item)
	return err
}

// DeleteItem deletes an item from the database
func DeleteItem(itemID string) error {
	if ItemDM == nil {
		return ErrEconomyManagerNotInitialized
	}
	return ItemDM.Delete(bson.M{"_id": itemID})
}
