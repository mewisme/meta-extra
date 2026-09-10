package socket

/*
	type 3 = task
*/

var TaskLabels = map[string]string{
	"UpdatePresence":                "3",
	"ThreadMarkRead":                "21",
	"AcceptMessageRequestTask":      "66",
	"AddParticipantsTask":           "23",
	"UpdateAdminTask":               "25",
	"SetThreadApprovalModeTask":     "28",
	"SendReactionTask":              "29",
	"SearchUserTask":                "30",
	"SearchUserSecondaryTask":       "31",
	"RenameThreadTask":              "32",
	"DeleteMessageTask":             "33",
	"SetThreadImageTask":            "37",
	"SetThreadThemeTask":            "43",
	"SetThreadNicknameTask":         "44",
	"SendMessageTask":               "46",
	"ShareContactTask":              "359",
	"SetThreadEmojiTask":            "53",
	"ReportAppStateTask":            "123",
	"CreateGroupTask":               "130",
	"RemoveParticipantTask":         "140",
	"MuteThreadTask":                "144",
	"FetchThreadsTask":              "145",
	"DeleteThreadTask":              "146",
	"DeleteMessageMeOnlyTask":       "155",
	"CreatePollTask":                "163",
	"UpdatePollTask":                "164",
	"GetContactsFullTask":           "207",
	"CreateThreadTask":              "209",
	"FetchMessagesTask":             "228",
	"MuteThreadCallsTask":           "229",
	"UnarchiveThreadTask":           "242",
	"FetchCommunityMemberList":      "355",
	"CreateWhatsAppThreadTask":      "388",
	"PinMessageTask":                "430",
	"UnpinMessageTask":              "431",
	"GetContactsTask":               "452",
	"CommunityThreadHoleDetection":  "501",
	"FetchPollDetailsTask":          "545",
	"FetchReactionsV2UserList":      "577",
	"SendReactionV2":                "604",
	"DeleteCommunitySubThread":      "639",
	"CreateCommunitySubThread":      "665",
	"FetchAdditionalThreadData":     "733",
	"EditMessageTask":               "742",
	"SetThreadThemeWriterTask":      "1028",
	"SetThreadThemeAIGeneratedTask": "1013",
	"SetThreadThemeCustomTask":      "1037",
}

type Task interface {
	GetLabel() string
	Create() (payload any, queueName string)
}

type TaskData struct {
	FailureCount *int64 `json:"failure_count"`
	Label        string `json:"label,omitempty"`
	Payload      string `json:"payload,omitempty"`
	QueueName    string `json:"queue_name,omitempty"`
	TaskID       int64  `json:"task_id"`
}

type StatelessTaskData struct {
	Label   string `json:"label,omitempty"`
	Payload string `json:"payload,omitempty"`
	Version string `json:"version,omitempty"`
}
