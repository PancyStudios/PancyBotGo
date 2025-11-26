package database

import (
	"errors"
	"time"

	"github.com/PancyStudios/PancyBotGo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrPremiumManagerNotInitialized = errors.New("premium data manager not initialized")
	ErrPremiumDurationRequired      = errors.New("duration must be greater than zero for non-permanent premium")
)

func nowMillis() int64 {
	return time.Now().UnixMilli()
}

func getUserPremiumManager() (*DataManager[models.UserPremium], error) {
	if GlobalUserPremiumDM == nil {
		return nil, ErrPremiumManagerNotInitialized
	}
	return GlobalUserPremiumDM, nil
}

func getGuildPremiumManager() (*DataManager[models.GuildPremium], error) {
	if GlobalGuildPremiumDM == nil {
		return nil, ErrPremiumManagerNotInitialized
	}
	return GlobalGuildPremiumDM, nil
}

func isRecordExpired(permanent bool, expiresAt int64) bool {
	if permanent {
		return false
	}
	return expiresAt <= nowMillis()
}

func GetUserPremium(userID string) (*models.UserPremium, error) {
	dm, err := getUserPremiumManager()
	if err != nil {
		return nil, err
	}

	query := bson.M{"user": userID}
	record, err := dm.Get(query)
	if err != nil {
		return nil, err
	}

	if record != nil && isRecordExpired(record.Permanent, record.ExpiresAt) {
		_ = dm.Delete(query)
		return nil, nil
	}

	return record, nil
}

func IsUserPremium(userID string) (bool, *models.UserPremium, error) {
	record, err := GetUserPremium(userID)
	if err != nil {
		return false, nil, err
	}
	return record != nil, record, nil
}

func GrantUserPremium(userID string, duration time.Duration, permanent bool) (*models.UserPremium, error) {
	dm, err := getUserPremiumManager()
	if err != nil {
		return nil, err
	}

	expiresAt := int64(0)
	if !permanent {
		if duration <= 0 {
			return nil, ErrPremiumDurationRequired
		}
		expiresAt = nowMillis() + duration.Milliseconds()
	}

	payload := models.UserPremium{
		UserID:    userID,
		Permanent: permanent,
		ExpiresAt: expiresAt,
	}

	query := bson.M{"user": userID}
	updated, err := dm.Set(query, payload)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func RemoveUserPremium(userID string) error {
	dm, err := getUserPremiumManager()
	if err != nil {
		return err
	}
	return dm.Delete(bson.M{"user": userID})
}

func GetGuildPremium(guildID string) (*models.GuildPremium, error) {
	dm, err := getGuildPremiumManager()
	if err != nil {
		return nil, err
	}

	query := bson.M{"guild": guildID}
	record, err := dm.Get(query)
	if err != nil {
		return nil, err
	}

	if record != nil && isRecordExpired(record.Permanent, record.ExpiresAt) {
		_ = dm.Delete(query)
		return nil, nil
	}

	return record, nil
}

func IsGuildPremium(guildID string) (bool, *models.GuildPremium, error) {
	record, err := GetGuildPremium(guildID)
	if err != nil {
		return false, nil, err
	}
	return record != nil, record, nil
}

func GrantGuildPremium(guildID string, duration time.Duration, permanent bool) (*models.GuildPremium, error) {
	dm, err := getGuildPremiumManager()
	if err != nil {
		return nil, err
	}

	expiresAt := int64(0)
	if !permanent {
		if duration <= 0 {
			return nil, ErrPremiumDurationRequired
		}
		expiresAt = nowMillis() + duration.Milliseconds()
	}

	payload := models.GuildPremium{
		GuildID:   guildID,
		Permanent: permanent,
		ExpiresAt: expiresAt,
	}

	query := bson.M{"guild": guildID}
	updated, err := dm.Set(query, payload)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func RemoveGuildPremium(guildID string) error {
	dm, err := getGuildPremiumManager()
	if err != nil {
		return err
	}
	return dm.Delete(bson.M{"guild": guildID})
}

// ===== Premium Code Services =====

var (
	ErrCodeNotFound       = errors.New("código no encontrado")
	ErrCodeAlreadyClaimed = errors.New("código ya reclamado")
	ErrCodeExists         = errors.New("código ya existe")
)

func getPremiumCodeManager() (*DataManager[models.PremiumCode], error) {
	if GlobalPremiumCodeDM == nil {
		return nil, ErrPremiumManagerNotInitialized
	}
	return GlobalPremiumCodeDM, nil
}

// CreatePremiumCode crea un nuevo código premium
func CreatePremiumCode(code string, codeType models.PremiumCodeType, durationDays int, permanent bool, createdBy string) (*models.PremiumCode, error) {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return nil, err
	}

	// Verificar si el código ya existe
	existing, err := dm.Get(bson.M{"_id": code})
	if err == nil && existing != nil {
		return nil, ErrCodeExists
	}

	premiumCode := models.PremiumCode{
		Code:         code,
		Type:         codeType,
		DurationDays: durationDays,
		Permanent:    permanent,
		IsClaimed:    false,
		CreatedAt:    time.Now(),
		CreatedBy:    createdBy,
	}

	// Insertar el código
	result, err := dm.Set(bson.M{"_id": code}, premiumCode)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPremiumCode obtiene un código premium por su código
func GetPremiumCode(code string) (*models.PremiumCode, error) {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": code}
	premiumCode, err := dm.Get(query)
	if err != nil {
		return nil, err
	}

	if premiumCode == nil {
		return nil, ErrCodeNotFound
	}

	return premiumCode, nil
}

// RedeemPremiumCode redime un código premium
func RedeemPremiumCode(code string, claimedBy string) (*models.PremiumCode, error) {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return nil, err
	}

	// Obtener el código
	premiumCode, err := GetPremiumCode(code)
	if err != nil {
		return nil, err
	}

	// Verificar si ya fue reclamado
	if premiumCode.IsClaimed {
		return nil, ErrCodeAlreadyClaimed
	}

	// Marcar como reclamado
	premiumCode.IsClaimed = true
	premiumCode.ClaimedBy = claimedBy
	premiumCode.ClaimedAt = time.Now()

	// Actualizar en la base de datos
	query := bson.M{"_id": code}
	updated, err := dm.Set(query, *premiumCode)
	if err != nil {
		return nil, err
	}

	// Otorgar el premium correspondiente
	if premiumCode.Type == models.PremiumCodeTypeUser {
		duration := time.Duration(premiumCode.DurationDays) * 24 * time.Hour
		_, err = GrantUserPremium(claimedBy, duration, premiumCode.Permanent)
		if err != nil {
			return nil, err
		}
	}

	return updated, nil
}

// RedeemPremiumCodeForGuild redime un código premium de guild
func RedeemPremiumCodeForGuild(code string, guildID string, claimedBy string) (*models.PremiumCode, error) {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return nil, err
	}

	// Obtener el código
	premiumCode, err := GetPremiumCode(code)
	if err != nil {
		return nil, err
	}

	// Verificar si ya fue reclamado
	if premiumCode.IsClaimed {
		return nil, ErrCodeAlreadyClaimed
	}

	// Verificar que sea un código de guild
	if premiumCode.Type != models.PremiumCodeTypeGuild {
		return nil, errors.New("este código es solo para usuarios, no para servidores")
	}

	// Marcar como reclamado
	premiumCode.IsClaimed = true
	premiumCode.ClaimedBy = claimedBy
	premiumCode.ClaimedAt = time.Now()

	// Actualizar en la base de datos
	query := bson.M{"_id": code}
	updated, err := dm.Set(query, *premiumCode)
	if err != nil {
		return nil, err
	}

	// Otorgar el premium al guild
	duration := time.Duration(premiumCode.DurationDays) * 24 * time.Hour
	_, err = GrantGuildPremium(guildID, duration, premiumCode.Permanent)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// GetAllPremiumCodes obtiene todos los códigos premium (para admin)
func GetAllPremiumCodes() ([]*models.PremiumCode, error) {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return nil, err
	}

	return dm.GetAll(bson.M{})
}

// DeletePremiumCode elimina un código premium
func DeletePremiumCode(code string) error {
	dm, err := getPremiumCodeManager()
	if err != nil {
		return err
	}

	return dm.Delete(bson.M{"_id": code})
}
