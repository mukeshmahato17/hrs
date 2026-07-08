package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/mukeshmahato17/hrs/db/fixture"
)

func TestAdminGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixture.AddUser(tdb.Store, "foo", "bar", false)
	hotel := fixture.AddHotel(&tdb.Store, "a", "b", nil, 5)
	room := fixture.AddRoom(&tdb.Store, "small", 45.9, hotel.ID)

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	booking := fixture.AddBooking(&tdb.Store, user.ID, room.ID, from, till)
	fmt.Println(booking)
}
