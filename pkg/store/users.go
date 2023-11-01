package store

import (
	"abeProofOfConcept/pkg/store/models"
	"errors"
	"golang.org/x/net/context"
)

type User = models.User

func GetUser(ctx context.Context, email string) (*User, error) {
	user := new(User)
	err := DB.NewSelect().Model(user).Where("email=?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if user.Email == "" {
		return nil, errors.New("no user found with the provided email")
	}
	return user, err
}

func CreateUser(
	ctx context.Context, email string, password []byte, firstName string, lastName string,
	position string, department string, phoneNumber string, salary []byte, address []byte,
) (*User, error) {
	user := &User{
		Email:       email,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		Position:    position,
		Department:  department,
		PhoneNumber: phoneNumber,
		Salary:      salary,
		Address:     address,
	}
	_, err := DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetPositions(ctx context.Context) ([]string, error) {
	var positions []string
	err := DB.NewSelect().Model((*User)(nil)).Column("position").Distinct().Scan(ctx, &positions)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func GetDepartments(ctx context.Context) ([]string, error) {
	var dpts []string
	err := DB.NewSelect().Model((*User)(nil)).Column("department").Distinct().Scan(ctx, &dpts)
	if err != nil {
		return nil, err
	}
	return dpts, nil
}

func GetAllUserEmails(ctx context.Context) ([]string, error) {
	var emails []string
	err := DB.NewSelect().Model((*User)(nil)).Column("email").Scan(ctx, &emails)
	if err != nil {
		return nil, err
	}
	return emails, nil
}
