package postgresql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vin-rmdn/general-ground/database/postgresql"
)

type postgreSQLTestSuite struct {
	suite.Suite

	client *postgresql.PostgreSQL
}

func TestPostgreSQL(t *testing.T) {
	suite.Run(t, new(postgreSQLTestSuite))
}

func (s *postgreSQLTestSuite) SetupSuite() {
	client, err := postgresql.New(context.Background(), postgresql.DatabaseConfig{
		Host:     "localhost",
		Port:     2001,
		User:     "general-ground",
		Password: "password",
		Database: "general-ground",
	})
	s.Require().NoError(err)

	s.client = client
}
