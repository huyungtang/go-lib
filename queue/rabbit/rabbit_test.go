package rabbit

import (
	"errors"
	"testing"
	"time"

	"github.com/huyungtang/go-lib/config"
	"github.com/huyungtang/go-lib/config/viper"
	"github.com/huyungtang/go-lib/file"
	"github.com/huyungtang/go-lib/queue"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestRabbit
// ****************************************************************************************************************************************
func TestRabbit(t *testing.T) {
	c, err := viper.Init(
		config.PathOption(file.PathWorking("../_testing")),
		config.EnvironmentOption("prod"),
	)
	if err != nil {
		t.Error(err)
	}

	cfg := &struct {
		TestStruct struct {
			EncString string `config:"key=Encrypt.EncString"`
		}
	}{}
	if err = c.GetStruct(cfg); err != nil {
		t.Error(err)
	}

	var db queue.Database
	if db, err = Init(
		cfg.TestStruct.EncString,
		queue.HandlerOption(func(mc queue.Message) {
			t.Log("database handler")
			mc.Next()
		}),
	); err != nil {
		t.Error(err)
	}
	defer db.Close()

	var exg queue.Exchange
	if exg, err = db.Exchange(
		"develop.testing",
		queue.AutoDeleteOption(true),
		queue.HandlerOption(func(mc queue.Message) {
			t.Log("exchange handler")
			mc.Next()
		}),
	); err != nil {
		t.Error(err)
	}
	defer exg.Close()

	if err = exg.QueueBind(
		"testing.queue",
		queue.AutoDeleteOption(true),
		queue.RoutingOption("testing.routing1"),
		queue.RoutingOption("testing.routing2"),
		queue.RoutingOption("testing.routing3"),
	); err != nil {
		t.Error(err)
	}

	errChan := make(chan error)
	go exg.Consume("testing.queue", errChan,
		queue.HandlerOption(func(mc queue.Message) {
			t.Log("consumer handler")
			t.Log(mc.Routing(), " | ", mc.String())
			if mc.Routing() == "testing.routing3" {
				errChan <- errors.New("done")
			}
		}),
	)

	go func() {
		time.Sleep(time.Second * 3)
		exg.Publish("testing.routing1",
			queue.StringMessageOption("testing.routing1.message1"),
			queue.StringMessageOption("testing.routing1.message2"),
			queue.StringMessageOption("testing.routing1.message3"),
			queue.StringMessageOption("testing.routing1.message4"),
		)
		time.Sleep(time.Second)
		exg.Publish("testing.routing2", queue.StringMessageOption("testing.routing2.message1"))
		time.Sleep(time.Second)
		exg.Publish("testing.routing3", queue.StringMessageOption("testing.routing3.message1"))
	}()

	if err = <-errChan; err != nil {
		if err.Error() != "done" {
			t.Error(err)
		}
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
