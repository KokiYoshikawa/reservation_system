package dto

type CreateReservationRequest struct {
	ServiceID int64  `json:"serviceId"`
	SlotID    int64  `json:"slotId"`
	Note      string `json:"note"`
}
