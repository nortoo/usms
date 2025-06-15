package usm

import (
	"github.com/nortoo/usm"
	"gorm.io/gorm"
)

var client *usm.Client

func Client() *usm.Client {
	return client
}

func Init(store, casbinStore *gorm.DB, casbinPolicyFile string) error {
	var err error
	client, err = usm.New(&usm.Options{
		Store: store,
		CasbinOptions: &usm.CasbinOptions{
			Store:      casbinStore,
			PolicyPath: casbinPolicyFile,
		},
	})
	return err
}
