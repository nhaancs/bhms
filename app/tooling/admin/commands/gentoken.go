package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nhaancs/realworld/business/core/event"
	"github.com/nhaancs/realworld/business/core/user"
	"github.com/nhaancs/realworld/business/core/user/stores/userdb"
	db "github.com/nhaancs/realworld/business/data/dbsql/pgx"
	"github.com/nhaancs/realworld/business/web/auth"
	"github.com/nhaancs/realworld/foundation/logger"
	"github.com/nhaancs/realworld/foundation/vault"
)

// GenToken generates a JWT for the specified user.
func GenToken(log *logger.Logger, dbConfig db.Config, vaultConfig vault.Config, userID uuid.UUID, kid string) error {
	if kid == "" {
		fmt.Println("help: gentoken <user_id> <kid>")
		return ErrHelp
	}

	db, err := db.Open(dbConfig)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	evnCore := event.NewCore(log)
	core := user.NewCore(log, evnCore, userdb.NewStore(log, db))

	usr, err := core.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("retrieve user: %w", err)
	}

	vault, err := vault.New(vaultConfig)
	if err != nil {
		return fmt.Errorf("new keystore: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		DB:        db,
		KeyLookup: vault,
	}

	a, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validate these claims. Dgraph does not support key rotate at
	// this time.
	token, err := a.GenerateToken(kid, claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", token)
	return nil
}
