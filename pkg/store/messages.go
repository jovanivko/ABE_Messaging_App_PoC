package store

import (
	"abeProofOfConcept/pkg/store/models"
	"context"
)

type Message = models.Message
type MsgFrg = models.MessageFragment

func GetMessages(ctx context.Context, email string) ([]Message, error) {
	var msgs []Message
	err := DB.NewSelect().Model(&msgs).Where("?=ANY(recipients)", email).WhereOr(
		"sender=?", email,
	).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func CreateMessage(ctx context.Context, sender string, recipients []string, title []byte) (*Message, error) {
	msg := &Message{
		Sender:     sender,
		Recipients: recipients,
		Title:      title,
	}
	_, err := DB.NewInsert().Model(msg).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func GetMessageFragments(ctx context.Context, msgID int) ([]MsgFrg, error) {
	var frgs []MsgFrg
	err := DB.NewSelect().Model(&frgs).Where("msg_id=?", msgID).Order("fragment_id ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return frgs, nil
}

func CreateMessageFragment(ctx context.Context, msgID int, content []byte, fragmentID int) (*MsgFrg, error) {
	frg := &MsgFrg{
		MsgID:      msgID,
		Content:    content,
		FragmentID: fragmentID,
	}

	_, err := DB.NewInsert().Model(frg).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return frg, nil
}
