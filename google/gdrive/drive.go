package gdrive

import (
	"context"

	"github.com/huyungtang/go-lib/google"
	base "google.golang.org/api/drive/v3"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(opts ...google.Options) (serv Service, err error) {
	cfg := new(google.Option).
		ApplyOptions(opts)

	var cal *base.Service
	if cal, err = base.NewService(context.Background(), cfg.GetClientOption()); err != nil {
		return
	}

	return &service{cal, cfg}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// service ********************************************************************************************************************************
type service struct {
	*base.Service
	*google.Option
}

// Service
// ****************************************************************************************************************************************
type Service interface {
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
