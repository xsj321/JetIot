package note

type UserNote struct {
	Id       int    `json:"id"`
	NoteName string `json:"note_name"`
	NoteData string `json:"note_data"`
}
