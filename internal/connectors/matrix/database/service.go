package database

import (
	"log/slog"

	"gorm.io/gorm"
)

type service struct {
	db     *gorm.DB
	logger *slog.Logger
}

func New(db *gorm.DB, logger *slog.Logger) (Service, error) {
	service := service{
		db:     db,
		logger: logger,
	}

	err := service.migrate()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (service *service) migrate() error {
	models := []any{
		MatrixRoom{},
		MatrixUser{},
		MatrixMessage{},
		MatrixEvent{},
	}

	for _, m := range models {
		err := service.db.AutoMigrate(m)
		if err != nil {
			return err
		}
	}

	err := service.fixMatrixMessageEventFK()
	if err != nil {
		return err
	}

	return nil
}

func (service *service) fixMatrixMessageEventFK() error {
	var deleteRule string

	err := service.db.Raw(`SELECT DELETE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_NAME = 'fk_matrix_messages_event'`).Row().Scan(&deleteRule)
	if err != nil {
		return err
	}

	if deleteRule != "SET NULL" {
		service.logger.Info("recreating matrix_messages foreign key relation")

		err = service.db.Exec(`ALTER TABLE matrix_messages DROP FOREIGN KEY fk_matrix_messages_event;`).Error
		if err != nil {
			return err
		}

		err = service.db.Exec(`ALTER TABLE matrix_messages ADD CONSTRAINT fk_matrix_messages_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL ON UPDATE RESTRICT`).Error
		if err != nil {
			return err
		}
	}

	return nil
}
