package sm

import (
	"errors"
	"fmt"

	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
)

const (
	Err_SubscriptionDataIsDamaged     = "subscription data is damaged"
	ErrF_SubscriptionRecordIsNotFound = "subscription record is not found, record=%v"
)

// checkDatabaseConsistency checks consistency of thread subscription records
// and user subscription records. This function is used in the scheduler and is
// also run once during the server's start.
func (srv *Server) checkDatabaseConsistency() (err error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var tsrs []sm.ThreadSubscriptionsRecord
	tsrs, err = srv.dbo.GetAllThreadSubscriptions()
	if err != nil {
		return err
	}

	var usrs []sm.UserSubscriptionsRecord
	usrs, err = srv.dbo.GetAllUserSubscriptions()
	if err != nil {
		return err
	}

	// Check.
	var subscriptionsA = make(map[sm.Subscription]struct{})
	for _, tsr := range tsrs {
		if tsr.Users.Size() > 0 {
			for _, userId := range *tsr.Users {
				subscriptionsA[sm.Subscription{ThreadId: tsr.ThreadId, UserId: userId}] = struct{}{}
			}
		}
	}

	var subscriptionsB = make(map[sm.Subscription]struct{})
	for _, usr := range usrs {
		if usr.Threads.Size() > 0 {
			for _, threadId := range *usr.Threads {
				subscriptionsB[sm.Subscription{ThreadId: threadId, UserId: usr.UserId}] = struct{}{}
			}
		}
	}

	if len(subscriptionsA) != len(subscriptionsB) {
		return errors.New(Err_SubscriptionDataIsDamaged)
	}

	var recordExists bool
	for key, _ := range subscriptionsA {
		_, recordExists = subscriptionsB[key]
		if !recordExists {
			return fmt.Errorf(ErrF_SubscriptionRecordIsNotFound, key)
		}
	}

	return nil
}
