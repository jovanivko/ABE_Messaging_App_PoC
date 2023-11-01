package models

import (
	"github.com/uptrace/bun"
	"time"
)

type MasterKey struct {
	bun.BaseModel `bun:"table:master_keys"`

	Scheme          string `json:"scheme" bun:"scheme,type:text,pk"` //either "fame" or "gpsw"
	PublicKey       []byte `json:"publicKey" bun:"public_key,type:bytea,notnull"`
	MasterSecretKey []byte `json:"msk" bun:"master_secret_key,type:bytea,notnull"` // '-' indicates that this field should not be included in JSON
}

type User struct {
	bun.BaseModel `bun:"table:users, alias:user"`

	Email       string `json:"email" bun:"email,type:text,pk"`
	Password    []byte `json:"password" bun:"password,notnull,type:bytea"`
	FirstName   string `json:"firstName" bun:"firstname,notnull,type:text"`
	LastName    string `json:"lastName" bun:"lastname,notnull,type:text"`
	Position    string `json:"position" bun:"position,notnull,type:text"`
	Department  string `json:"department" bun:"department,notnull,type:text"`
	PhoneNumber string `json:"phoneNumber" bun:"phone,notnull,type:text"`
	Salary      []byte `json:"salary,omitempty" bun:"salary,notnull,type:bytea"`
	Address     []byte `json:"address,omitempty" bun:"address,notnull,type:bytea"`
}

type Message struct {
	bun.BaseModel `bun:"table:messages"`

	ID         int       `json:"id" bun:"id,pk,autoincrement"`
	Sender     string    `json:"sender" bun:"sender,type:text,notnull"`
	Recipients []string  `json:"recipients" bun:"recipients,array,notnull"`
	Title      []byte    `json:"title,omitempty" bun:"title,type:bytea,notnull"`
	CreatedAt  time.Time `json:"createdAt" bun:"created_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
}

type MessageFragment struct {
	bun.BaseModel `bun:"table:message_fragments"`

	ID         int    `json:"id" bun:"pk,autoincrement"`
	MsgID      int    `json:"msg_id" bun:"msg_id,type:integer,notnull"`
	Content    []byte `json:"content,omitempty" bun:"content,type:bytea,notnull"`
	FragmentID int    `json:"fragment_id" bun:"fragment_id,type:integer,notnull"`
}
