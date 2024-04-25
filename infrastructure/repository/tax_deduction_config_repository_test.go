package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thitiphum-bluesage/assessment-tax/domains"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMock() (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock: " + err.Error())
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		db.Close()
		panic("failed to open gorm database: " + err.Error())
	} 
	cleanup := func() {
		db.Close()
	}

	return gdb, mock, cleanup
}

func TestGetConfig(t *testing.T) {
	gdb, mock, cleanup := setupMock()
	defer cleanup()

	taxRepo := NewTaxDeductionConfigRepository(gdb)

	expectedConfig := domains.TaxDeductionConfig{
		ConfigName:           "MainConfig",
		PersonalDeduction:    60000,
		KReceiptDeductionMax: 50000,
		DonationDeductionMax: 100000,
	}

	rows := sqlmock.NewRows([]string{"config_name", "personal_deduction", "k_receipt_deduction_max", "donation_deduction_max"}).
		AddRow(expectedConfig.ConfigName, expectedConfig.PersonalDeduction, expectedConfig.KReceiptDeductionMax, expectedConfig.DonationDeductionMax)

	mock.ExpectQuery(`SELECT \* FROM "tax_deduction_configs" WHERE config_name = \$1 ORDER BY "tax_deduction_configs"\."config_name" LIMIT \$2`).
		WithArgs("MainConfig", 1).
		WillReturnRows(rows)

	config, err := taxRepo.GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, expectedConfig, *config)

	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestUpdatePersonalDeduction(t *testing.T) {
	gdb, mock, cleanup := setupMock()
	defer cleanup()

	taxRepo := NewTaxDeductionConfigRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "tax_deduction_configs" SET "personal_deduction"=\$1 WHERE config_name = \$2`).
		WithArgs(70000.0, "MainConfig").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := taxRepo.UpdatePersonalDeduction(70000)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "tax_deduction_configs" SET "personal_deduction"=\$1 WHERE config_name = \$2`).
		WithArgs(70000.0, "MainConfig").
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = taxRepo.UpdatePersonalDeduction(70000)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateKReceiptDeductionMax(t *testing.T) {
	gdb, mock, cleanup := setupMock()
	defer cleanup()

	taxRepo := NewTaxDeductionConfigRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "tax_deduction_configs" SET "k_receipt_deduction_max"=\$1 WHERE config_name = \$2`).
		WithArgs(45000.0, "MainConfig").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := taxRepo.UpdateKReceiptDeductionMax(45000)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "tax_deduction_configs" SET "k_receipt_deduction_max"=\$1 WHERE config_name = \$2`).
		WithArgs(45000.0, "MainConfig").
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = taxRepo.UpdateKReceiptDeductionMax(45000)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}


