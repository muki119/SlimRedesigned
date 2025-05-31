package Models

// type ModelErrors struct {
// 	groupTableError            error
// 	userTableError             error
// 	groupMessageTableError     error
// 	groupParticipantTableError error
// }

func InitialiseModels() {
	groupTableError := CreateGroupTable()
	userTableError := CreateUserTable()
	groupMessageTableError := CreateGroupMessageTable()
	groupParticipantTableError := CreateGroupParticipantTable()

	if groupTableError != nil || userTableError != nil || groupMessageTableError != nil || groupParticipantTableError != nil {
		return
	}
}
