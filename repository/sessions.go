package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	u.db.Create(&session)
	return nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	u.db.Where("token = ?", tokenTarget).Delete(&model.Session{})
	return nil
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	u.db.Model(&model.Session{}).Where("username = ?", session.Username).Updates(session)
	return nil
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	var result model.Session
	tx := u.db.Raw("SELECT * FROM sessions WHERE token = ?", token).Scan(&result)
	if tx.Error != nil {
		return model.Session{}, tx.Error
	}

	session := result

	if u.TokenExpired(session) {
		err := u.DeleteSessions(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return session, nil
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	var result model.Session
	tx := u.db.Raw("SELECT * FROM sessions WHERE username = ?", name).Scan(&result)
	if tx.Error != nil {
		return model.Session{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return model.Session{}, errors.New("erorr guys")
	}
	return result, nil
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	var result model.Session
	if err := u.db.Where("token = ?", token).First(&result).Error; err != nil {
		return model.Session{}, errors.New("record not found")
	}
	return result, nil
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
