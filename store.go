package main

import (
	"context"

	"gorm.io/gorm"
)

// Contact represents user data
type Contact struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" form:"name"`
	Phone string `json:"phone" form:"phone" gorm:"index:unique"`
}

// Store is interface for contacts storage
type Store interface {
	Contact(ctx context.Context, id uint) (Contact, error)
	Contacts(ctx context.Context) ([]Contact, error)
	CreateContact(ctx context.Context, c *Contact) error
	EditContact(ctx context.Context, id uint, c Contact) error
	DeleteContact(ctx context.Context, id uint) error
}

type model struct {
	db *gorm.DB
}

func (m *model) Contact(ctx context.Context, id uint) (Contact, error) {
	var contact Contact

	if err := m.db.WithContext(ctx).First(&contact, id).Error; err != nil {
		return Contact{}, err
	}

	return contact, nil
}

func (m *model) Contacts(ctx context.Context) ([]Contact, error) {
	var contacts []Contact

	if err := m.db.WithContext(ctx).Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

func (m *model) CreateContact(ctx context.Context, c *Contact) error {
	return m.db.WithContext(ctx).Create(&c).Error
}

func (m *model) EditContact(ctx context.Context, id uint, c Contact) error {
	c.ID = id
	return m.db.WithContext(ctx).Updates(&c).Error
}

func (m *model) DeleteContact(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Delete(&Contact{}, id).Error
}
