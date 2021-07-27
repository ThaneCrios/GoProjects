package proto

type UsersFilter struct {
	UserUUID  string `json:"user_uuid"`
	UserState string `json:"user_state"`
}
