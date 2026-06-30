package vault_test

import (
	"pvault/errors"
	"pvault/vault"
	"pvault/vault/data"
	"pvault/vault/data/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	suite.Suite
	DatabaseMock data.DatabaseMock
	Vault        vault.Vault
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, &UpgradeTestSuite{})
}

func (s *UpgradeTestSuite) SetupTest() {
	s.DatabaseMock = data.DatabaseMock{
		Version: vault.CURRENT_VERSION - 1,
	}

	s.Vault = vault.Vault{
		Path:     file.NewPath(s.T(), ""),
		Database: &s.DatabaseMock,
	}
}

func (s *UpgradeTestSuite) TestIsOutOfDateWhereVersionLessThanCurrentReturnsTrue() {
	//-- act
	res := s.Vault.IsOutOfDate()

	//-- arrange
	s.Require().True(res)
}

func (s *UpgradeTestSuite) TestUpgradeWhereVaultIsNotOutOfDateReturnsError() {
	//-- arrange
	s.DatabaseMock.Version = vault.CURRENT_VERSION

	//-- act
	res := s.Vault.Upgrade()

	//-- arrange
	s.Require().ErrorContains(res, "vault is up-to-date")
}

func (s *UpgradeTestSuite) TestUpgradeWhereDatabaseUpgradeFailsReturnsError() {
	//-- arrange
	s.DatabaseMock.UpgradeError = errors.New("")

	//-- act
	res := s.Vault.Upgrade()

	//-- arrange
	s.Require().ErrorContains(res, "error upgrading database")
}

func (s *UpgradeTestSuite) TestUpgradeValidRunsUpgrade() {
	//-- act
	res := s.Vault.Upgrade()

	//-- arrange
	s.Require().NoError(res)

	s.Require().IsType(version2.Database{}, s.Vault.Database)
	s.Assert().FileExists(s.Vault.Database.(version2.Database).IndexPath())
}
