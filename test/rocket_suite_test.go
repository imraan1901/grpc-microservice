package test

import (
	"context"
	"testing"

	"github.com/imraan1901/grpc-proto/rocket/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RocketTestSuite struct {
	suite.Suite
}

func (s *RocketTestSuite) TestAddRocket() {
	s.T().Run("Adds a new rocket successfully", func(t *testing.T) {
		client := GetClient()
		resp, err := client.AddRocket(
			context.Background(),
			&rocket.AddRocketRequest{
				Rocket: &rocket.Rocket{
					Id:   "a2ffc12e-0b8e-4b07-a12c-df26634874c6",
					Name: "11",
					Type: "Apollo",
				},
			},
		)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "a2ffc12e-0b8e-4b07-a12c-df26634874c6", resp.Rocket.Id)
	})

	s.T().Run("Validates the uuid in the new rocket is a uuid", func(t *testing.T) {
		client := GetClient()
		_, err := client.AddRocket(
			context.Background(),
			&rocket.AddRocketRequest{
				Rocket: &rocket.Rocket{
					Id:   "not-a-valid-uuid",
					Name: "11",
					Type: "Apollo",
				},
			},
		)
		assert.Error(s.T(), err)
		st := status.Convert(err)
		assert.Equal(s.T(), codes.InvalidArgument, st.Code())

	})
}

func TestRocketService(t *testing.T) {
	suite.Run(t, new(RocketTestSuite))
}
