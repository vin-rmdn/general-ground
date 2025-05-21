package postgresql_test

import (
	"context"

	uuid "github.com/gofrs/uuid/v5"
)

func (s *postgreSQLTestSuite) TestCreateUser() {
	s.Run("when insertion fails, it returns an error", func() {
		userID, err := uuid.NewV4()
		s.Require().NoError(err)

		ctx := context.Background()
		ctxWithCancel, cancelFunc := context.WithCancel(ctx)
		cancelFunc()

		err = s.client.CreateUser(ctxWithCancel, userID.String())
		s.Require().EqualError(err, "failed to create user: timeout: context already done: context canceled")
	})

	s.Run("when user is created successfully", func() {
		userID, err := uuid.NewV4()
		s.Require().NoError(err)

		err = s.client.CreateUser(context.Background(), userID.String())
		s.Require().NoError(err)
	})
}
