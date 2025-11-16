package admin

import (
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/storage/admin/dbmodel"
)

func mapDBAdminToDomainAdmin(admin dbmodel.Admin) model.Admin {
	return model.Admin{
		ID:           admin.ID,
		PasswordHash: admin.PasswordHash,
	}
}
