package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ammac123/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, name)
	if (err != nil) || (user.Name != name) {
		return fmt.Errorf("error retrieving user")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	log.Printf("UserName: %+v", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	ctx := context.Background()
	now := time.Now()

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	})
	if err != nil {
		log.Printf("%v\n", err)
		return fmt.Errorf("could not create new user")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Println("User successfully created!")
	log.Printf("%+v\n", user)
	return nil
}
