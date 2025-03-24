package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	if p.FromDate.After(p.TillDate) {
		return fmt.Errorf("fromDate cannot be greater than tillDate")
	}

	return nil
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID     primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons int                `bson:"numPersons" json:"numPersons"`
	FromDate   time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}
