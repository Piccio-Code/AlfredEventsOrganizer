package waha

const (
	StatusWorking    = "WORKING"
	StatusScanQRCode = "SCAN_QR_CODE"
)

// I tipi esportati qui sotto sono il contratto verso il resto del backend.
// Le struct new* piu' in basso sono lo strato di trasporto: servono solo a
// deserializzare le risposte di WAHA e vengono mappate su questi.

type SessionResponse struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Phone    string `json:"phone"`
	PushName string `json:"pushName"`
}

type ParingCodeResponse struct {
	Code string `json:"code"`
}

type GroupResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GroupDetailResponse struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Owner        string             `json:"owner"`
	IsReadOnly   bool               `json:"isReadOnly"`
	IsAnnounce   bool               `json:"isAnnounce"`
	Participants []GroupParticipant `json:"participants"`
}

type GroupParticipant struct {
	Id           string `json:"id"`
	Number       string `json:"number"`
	Name         string `json:"name"`
	IsAdmin      bool   `json:"isAdmin"`
	IsSuperAdmin bool   `json:"isSuperAdmin"`
}

type ContactInfo struct {
	Id       string `json:"id"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	PushName string `json:"pushName"`
}

type SessionWebHook struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Session   string `json:"session"`
	Metadata  struct {
		UserId    string `json:"user.id"`
		UserEmail string `json:"user.email"`
	} `json:"metadata"`
	Engine  string `json:"engine"`
	Event   string `json:"event"`
	Payload struct {
		Name     string      `json:"name"`
		Data     interface{} `json:"data"`
		Status   string      `json:"status"`
		Statuses []struct {
			Status    string `json:"status"`
			Timestamp int    `json:"timestamp"`
		} `json:"statuses"`
	} `json:"payload"`
	Me struct {
		Id               string `json:"id"`
		Lid              string `json:"lid"`
		Jid              string `json:"jid"`
		ReachoutTimelock struct {
			EnforcementType     string `json:"enforcementType"`
			IsActive            bool   `json:"isActive"`
			TimeEnforcementEnds int    `json:"timeEnforcementEnds"`
		} `json:"reachoutTimelock"`
		MessageCapping struct {
			CappingStatus string `json:"cappingStatus"`
			TotalQuota    int    `json:"totalQuota"`
			UsedQuota     int    `json:"usedQuota"`
			CycleStart    int    `json:"cycleStart"`
			CycleEnd      int    `json:"cycleEnd"`
			MvStatus      string `json:"mvStatus"`
			OteStatus     string `json:"oteStatus"`
		} `json:"messageCapping"`
		PushName string `json:"pushName"`
	} `json:"me"`
	Environment struct {
		Version  string `json:"version"`
		Engine   string `json:"engine"`
		Tier     string `json:"tier"`
		Browser  string `json:"browser"`
		Platform string `json:"platform"`
		Worker   struct {
			Id string `json:"id"`
		} `json:"worker"`
	} `json:"environment"`
}

// previousStatus restituisce lo stato precedente a quello corrente. In
// Payload.Statuses WAHA manda la cronologia e l'ultimo elemento e' lo stato
// corrente, quindi il precedente e' il penultimo.
func (s SessionWebHook) previousStatus() string {
	if len(s.Payload.Statuses) < 2 {
		return ""
	}

	return s.Payload.Statuses[len(s.Payload.Statuses)-2].Status
}

// BecameWorking indica che la sessione e' appena diventata operativa, cosi' la
// notifica non parte a ogni evento WORKING che WAHA rimanda.
func (s SessionWebHook) BecameWorking() bool {
	return s.Payload.Status == StatusWorking && s.previousStatus() != StatusWorking
}

// NeedsPairing indica che WAHA sta aspettando QR code o pairing code.
func (s SessionWebHook) NeedsPairing() bool {
	return s.Payload.Status == StatusScanQRCode
}

type newSessionResponse struct {
	Name string `json:"name"`
	Apps []struct {
		Enabled bool   `json:"enabled"`
		Id      string `json:"id"`
		Session string `json:"session"`
		App     string `json:"app"`
		Config  struct {
		} `json:"config"`
	} `json:"apps"`
	Me struct {
		Id               string `json:"id"`
		Lid              string `json:"lid"`
		Jid              string `json:"jid"`
		ReachoutTimelock struct {
			EnforcementType     string `json:"enforcementType"`
			IsActive            bool   `json:"isActive"`
			TimeEnforcementEnds int    `json:"timeEnforcementEnds"`
		} `json:"reachoutTimelock"`
		MessageCapping struct {
			CappingStatus string `json:"cappingStatus"`
			TotalQuota    int    `json:"totalQuota"`
			UsedQuota     int    `json:"usedQuota"`
			CycleStart    int    `json:"cycleStart"`
			CycleEnd      int    `json:"cycleEnd"`
			MvStatus      string `json:"mvStatus"`
			OteStatus     string `json:"oteStatus"`
		} `json:"messageCapping"`
		PushName string `json:"pushName"`
	} `json:"me"`
	AssignedWorker string `json:"assignedWorker"`
	Presence       struct {
	} `json:"presence"`
	Timestamps struct {
		Activity int `json:"activity"`
	} `json:"timestamps"`
	Status string `json:"status"`
	Config struct {
		Metadata struct {
			UserId    string `json:"user.id"`
			UserEmail string `json:"user.email"`
		} `json:"metadata"`
		Proxy  interface{} `json:"proxy"`
		Debug  bool        `json:"debug"`
		Ignore struct {
			Status   interface{} `json:"status"`
			Groups   interface{} `json:"groups"`
			Channels interface{} `json:"channels"`
		} `json:"ignore"`
		Client struct {
			BrowserName string `json:"browserName"`
			DeviceName  string `json:"deviceName"`
		} `json:"client"`
		Noweb struct {
			Store struct {
				Enabled  bool `json:"enabled"`
				FullSync bool `json:"fullSync"`
			} `json:"store"`
		} `json:"noweb"`
		Gows struct {
			Storage struct {
				Messages       bool `json:"messages"`
				Groups         bool `json:"groups"`
				Chats          bool `json:"chats"`
				Labels         bool `json:"labels"`
				Contacts       bool `json:"contacts"`
				MessageSecrets bool `json:"messageSecrets"`
			} `json:"storage"`
		} `json:"gows"`
		Webjs struct {
			TagsEventsOn bool `json:"tagsEventsOn"`
		} `json:"webjs"`
		Webhooks []struct {
			Url           string      `json:"url"`
			Events        []string    `json:"events"`
			Hmac          interface{} `json:"hmac"`
			Retries       interface{} `json:"retries"`
			CustomHeaders interface{} `json:"customHeaders"`
		} `json:"webhooks"`
	} `json:"config"`
}

type newParingCode struct {
	Code string `json:"code"`
}

type newGroupResponse struct {
	GroupMetadata struct {
		Id struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"id"`
		Creation int `json:"creation"`
		Owner    struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"owner"`
		Subject   string `json:"subject"`
		Desc      string `json:"desc"`
		DescId    string `json:"descId"`
		DescTime  int    `json:"descTime"`
		DescOwner struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"descOwner"`
		Restrict                  bool        `json:"restrict"`
		Announce                  bool        `json:"announce"`
		NoFrequentlyForwarded     bool        `json:"noFrequentlyForwarded"`
		EphemeralDuration         int         `json:"ephemeralDuration"`
		AfterReadDuration         interface{} `json:"afterReadDuration"`
		MembershipApprovalMode    bool        `json:"membershipApprovalMode"`
		MemberAddMode             string      `json:"memberAddMode"`
		ReportToAdminMode         bool        `json:"reportToAdminMode"`
		Support                   bool        `json:"support"`
		Suspended                 bool        `json:"suspended"`
		Terminated                bool        `json:"terminated"`
		SuspendAppealStatus       interface{} `json:"suspendAppealStatus"`
		SuspendAppealUpdateTime   interface{} `json:"suspendAppealUpdateTime"`
		SuspendAppealApprovedSeen bool        `json:"suspendAppealApprovedSeen"`
		UniqueShortNameMap        struct {
		} `json:"uniqueShortNameMap"`
		IsLidAddressingMode           bool        `json:"isLidAddressingMode"`
		IsParentGroup                 bool        `json:"isParentGroup"`
		IsParentGroupClosed           bool        `json:"isParentGroupClosed"`
		DefaultSubgroup               bool        `json:"defaultSubgroup"`
		GeneralSubgroup               bool        `json:"generalSubgroup"`
		HiddenSubgroup                bool        `json:"hiddenSubgroup"`
		GroupSafetyCheck              bool        `json:"groupSafetyCheck"`
		GeneralChatAutoAddDisabled    bool        `json:"generalChatAutoAddDisabled"`
		AllowNonAdminSubGroupCreation bool        `json:"allowNonAdminSubGroupCreation"`
		LastReportToAdminTimestamp    interface{} `json:"lastReportToAdminTimestamp"`
		HasCapi                       bool        `json:"hasCapi"`
		MemberShareGroupHistoryMode   string      `json:"memberShareGroupHistoryMode"`
		Participants                  []struct {
			Id struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"id"`
			IsAdmin               bool `json:"isAdmin"`
			IsSuperAdmin          bool `json:"isSuperAdmin"`
			GroupHistorySentState int  `json:"groupHistorySentState"`
		} `json:"participants"`
		PendingParticipants []interface{} `json:"pendingParticipants"`
		PastParticipants    []struct {
			Id struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"id"`
			LeaveTs     int    `json:"leaveTs"`
			LeaveReason string `json:"leaveReason"`
		} `json:"pastParticipants"`
		MembershipApprovalRequests []interface{} `json:"membershipApprovalRequests"`
		SubgroupSuggestions        []interface{} `json:"subgroupSuggestions"`
	} `json:"groupMetadata"`
	Id struct {
		Server     string `json:"server"`
		User       string `json:"user"`
		Serialized string `json:"_serialized"`
	} `json:"id"`
	Name           string `json:"name"`
	IsGroup        bool   `json:"isGroup"`
	IsReadOnly     bool   `json:"isReadOnly"`
	UnreadCount    int    `json:"unreadCount"`
	Timestamp      int    `json:"timestamp"`
	Archived       bool   `json:"archived"`
	Pinned         bool   `json:"pinned"`
	IsLocked       bool   `json:"isLocked"`
	IsMuted        bool   `json:"isMuted"`
	MuteExpiration int    `json:"muteExpiration"`
	LastMessage    struct {
		Data struct {
			Id struct {
				FromMe      bool   `json:"fromMe"`
				Remote      string `json:"remote"`
				Id          string `json:"id"`
				Participant struct {
					Server     string `json:"server"`
					User       string `json:"user"`
					Serialized string `json:"_serialized"`
				} `json:"participant"`
				Field5     string `json:"$1"`
				Serialized string `json:"_serialized"`
			} `json:"id"`
			Viewed     bool   `json:"viewed"`
			Body       string `json:"body"`
			Type       string `json:"type"`
			T          int    `json:"t"`
			NotifyName string `json:"notifyName"`
			From       struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"from"`
			To struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"to"`
			Author struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"author"`
			Ack                         int           `json:"ack"`
			Invis                       bool          `json:"invis"`
			Star                        bool          `json:"star"`
			KicNotified                 bool          `json:"kicNotified"`
			IsFromTemplate              bool          `json:"isFromTemplate"`
			IsAdsMedia                  bool          `json:"isAdsMedia"`
			PollInvalidated             bool          `json:"pollInvalidated"`
			IsSentCagPollCreation       bool          `json:"isSentCagPollCreation"`
			LatestEditMsgKey            interface{}   `json:"latestEditMsgKey"`
			LatestEditSenderTimestampMs interface{}   `json:"latestEditSenderTimestampMs"`
			Broadcast                   bool          `json:"broadcast"`
			MentionedJidList            []interface{} `json:"mentionedJidList"`
			GroupMentions               []interface{} `json:"groupMentions"`
			IsEventCanceled             bool          `json:"isEventCanceled"`
			EventInvalidated            bool          `json:"eventInvalidated"`
			IsVcardOverMmsDocument      bool          `json:"isVcardOverMmsDocument"`
			IsForwarded                 bool          `json:"isForwarded"`
			IsQuestion                  bool          `json:"isQuestion"`
			IsSpoiler                   bool          `json:"isSpoiler"`
			QuestionReplyQuotedMessage  interface{}   `json:"questionReplyQuotedMessage"`
			QuestionResponsesCount      int           `json:"questionResponsesCount"`
			ReadQuestionResponsesCount  int           `json:"readQuestionResponsesCount"`
			ForwardsCount               int           `json:"forwardsCount"`
			Labels                      []interface{} `json:"labels"`
			HasReaction                 bool          `json:"hasReaction"`
			NewsletterAdminProfile      interface{}   `json:"newsletterAdminProfile"`
			ExpiredTimestamp            interface{}   `json:"expiredTimestamp"`
			ViewMode                    string        `json:"viewMode"`
			MessageSecret               struct {
				Field1  int `json:"0"`
				Field2  int `json:"1"`
				Field3  int `json:"2"`
				Field4  int `json:"3"`
				Field5  int `json:"4"`
				Field6  int `json:"5"`
				Field7  int `json:"6"`
				Field8  int `json:"7"`
				Field9  int `json:"8"`
				Field10 int `json:"9"`
				Field11 int `json:"10"`
				Field12 int `json:"11"`
				Field13 int `json:"12"`
				Field14 int `json:"13"`
				Field15 int `json:"14"`
				Field16 int `json:"15"`
				Field17 int `json:"16"`
				Field18 int `json:"17"`
				Field19 int `json:"18"`
				Field20 int `json:"19"`
				Field21 int `json:"20"`
				Field22 int `json:"21"`
				Field23 int `json:"22"`
				Field24 int `json:"23"`
				Field25 int `json:"24"`
				Field26 int `json:"25"`
				Field27 int `json:"26"`
				Field28 int `json:"27"`
				Field29 int `json:"28"`
				Field30 int `json:"29"`
				Field31 int `json:"30"`
				Field32 int `json:"31"`
			} `json:"messageSecret"`
			ProductHeaderImageRejected   bool        `json:"productHeaderImageRejected"`
			LastPlaybackProgress         int         `json:"lastPlaybackProgress"`
			IsDynamicReplyButtonsMsg     bool        `json:"isDynamicReplyButtonsMsg"`
			IsCarouselCard               bool        `json:"isCarouselCard"`
			ParentMsgId                  interface{} `json:"parentMsgId"`
			CallSilenceReason            interface{} `json:"callSilenceReason"`
			IsVideoCall                  bool        `json:"isVideoCall"`
			CallDuration                 interface{} `json:"callDuration"`
			CallCreator                  interface{} `json:"callCreator"`
			CallParticipants             interface{} `json:"callParticipants"`
			IsCallLink                   interface{} `json:"isCallLink"`
			CallLinkToken                interface{} `json:"callLinkToken"`
			TerminatedByDeviceSwitch     interface{} `json:"terminatedByDeviceSwitch"`
			SelfOtherDeviceConnected     interface{} `json:"selfOtherDeviceConnected"`
			BytesSent                    interface{} `json:"bytesSent"`
			BytesReceived                interface{} `json:"bytesReceived"`
			IsMdHistoryMsg               bool        `json:"isMdHistoryMsg"`
			StickerSentTs                int         `json:"stickerSentTs"`
			IsAvatar                     bool        `json:"isAvatar"`
			LastUpdateFromServerTs       int         `json:"lastUpdateFromServerTs"`
			InvokedBotWid                interface{} `json:"invokedBotWid"`
			BotTargetSenderJid           interface{} `json:"botTargetSenderJid"`
			MetaFrom                     interface{} `json:"metaFrom"`
			BizBotType                   interface{} `json:"bizBotType"`
			BotDeepLinkToken             interface{} `json:"botDeepLinkToken"`
			AiMediaCollectionInfo        interface{} `json:"aiMediaCollectionInfo"`
			ExpectedImageCount           interface{} `json:"expectedImageCount"`
			ExpectedVideoCount           interface{} `json:"expectedVideoCount"`
			BotResponseTargetId          interface{} `json:"botResponseTargetId"`
			BotPluginType                interface{} `json:"botPluginType"`
			BotPluginReferenceIndex      interface{} `json:"botPluginReferenceIndex"`
			BotPluginSearchProvider      interface{} `json:"botPluginSearchProvider"`
			BotPluginSearchUrl           interface{} `json:"botPluginSearchUrl"`
			BotPluginSearchQuery         interface{} `json:"botPluginSearchQuery"`
			BotPluginMaybeParent         bool        `json:"botPluginMaybeParent"`
			BotReelPluginThumbnailCdnUrl interface{} `json:"botReelPluginThumbnailCdnUrl"`
			BotMessageDisclaimerText     interface{} `json:"botMessageDisclaimerText"`
			BotSessionTransparencyType   interface{} `json:"botSessionTransparencyType"`
			BotMsgBodyType               interface{} `json:"botMsgBodyType"`
			BotModeSelection             interface{} `json:"botModeSelection"`
			BotModeOverride              interface{} `json:"botModeOverride"`
			ReportingTokenInfo           struct {
				ReportingTag struct {
				} `json:"reportingTag"`
				Version int `json:"version"`
			} `json:"reportingTokenInfo"`
			HatchMetadataSync                     interface{}   `json:"hatchMetadataSync"`
			RequiresDirectConnection              bool          `json:"requiresDirectConnection"`
			BizContentPlaceholderType             interface{}   `json:"bizContentPlaceholderType"`
			HostedBizEncStateMismatch             bool          `json:"hostedBizEncStateMismatch"`
			SenderOrRecipientAccountTypeHosted    bool          `json:"senderOrRecipientAccountTypeHosted"`
			PlaceholderCreatedWhenAccountIsHosted bool          `json:"placeholderCreatedWhenAccountIsHosted"`
			GroupHistoryBundleMessageKey          interface{}   `json:"groupHistoryBundleMessageKey"`
			GroupHistoryBundleMetadata            interface{}   `json:"groupHistoryBundleMetadata"`
			GroupHistoryIndividualMessageInfo     interface{}   `json:"groupHistoryIndividualMessageInfo"`
			NonJidMentions                        interface{}   `json:"nonJidMentions"`
			Links                                 []interface{} `json:"links"`
		} `json:"_data"`
		Id struct {
			FromMe      bool   `json:"fromMe"`
			Remote      string `json:"remote"`
			Id          string `json:"id"`
			Participant struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"participant"`
			Field5     string `json:"$1"`
			Serialized string `json:"_serialized"`
		} `json:"id"`
		Ack             int           `json:"ack"`
		HasMedia        bool          `json:"hasMedia"`
		Body            string        `json:"body"`
		Type            string        `json:"type"`
		Timestamp       int           `json:"timestamp"`
		From            string        `json:"from"`
		To              string        `json:"to"`
		Author          string        `json:"author"`
		DeviceType      string        `json:"deviceType"`
		IsForwarded     bool          `json:"isForwarded"`
		ForwardingScore int           `json:"forwardingScore"`
		IsStatus        bool          `json:"isStatus"`
		IsStarred       bool          `json:"isStarred"`
		Broadcast       bool          `json:"broadcast"`
		FromMe          bool          `json:"fromMe"`
		HasQuotedMsg    bool          `json:"hasQuotedMsg"`
		HasReaction     bool          `json:"hasReaction"`
		VCards          []interface{} `json:"vCards"`
		MentionedIds    []interface{} `json:"mentionedIds"`
		GroupMentions   []interface{} `json:"groupMentions"`
		IsGif           bool          `json:"isGif"`
		Links           []interface{} `json:"links"`
	} `json:"lastMessage"`
}

type newGroupDetailResponse struct {
	GroupMetadata struct {
		Id struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"id"`
		Creation int `json:"creation"`
		Owner    struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"owner"`
		Subject   string `json:"subject"`
		Desc      string `json:"desc"`
		DescId    string `json:"descId"`
		DescTime  int    `json:"descTime"`
		DescOwner struct {
			Server     string `json:"server"`
			User       string `json:"user"`
			Serialized string `json:"_serialized"`
		} `json:"descOwner"`
		Restrict                  bool        `json:"restrict"`
		Announce                  bool        `json:"announce"`
		NoFrequentlyForwarded     bool        `json:"noFrequentlyForwarded"`
		EphemeralDuration         int         `json:"ephemeralDuration"`
		AfterReadDuration         interface{} `json:"afterReadDuration"`
		MembershipApprovalMode    bool        `json:"membershipApprovalMode"`
		MemberAddMode             string      `json:"memberAddMode"`
		ReportToAdminMode         bool        `json:"reportToAdminMode"`
		Support                   bool        `json:"support"`
		Suspended                 bool        `json:"suspended"`
		Terminated                bool        `json:"terminated"`
		SuspendAppealStatus       interface{} `json:"suspendAppealStatus"`
		SuspendAppealUpdateTime   interface{} `json:"suspendAppealUpdateTime"`
		SuspendAppealApprovedSeen bool        `json:"suspendAppealApprovedSeen"`
		UniqueShortNameMap        struct {
		} `json:"uniqueShortNameMap"`
		IsLidAddressingMode           bool                  `json:"isLidAddressingMode"`
		IsParentGroup                 bool                  `json:"isParentGroup"`
		IsParentGroupClosed           bool                  `json:"isParentGroupClosed"`
		DefaultSubgroup               bool                  `json:"defaultSubgroup"`
		GeneralSubgroup               bool                  `json:"generalSubgroup"`
		HiddenSubgroup                bool                  `json:"hiddenSubgroup"`
		GroupSafetyCheck              bool                  `json:"groupSafetyCheck"`
		GeneralChatAutoAddDisabled    bool                  `json:"generalChatAutoAddDisabled"`
		AllowNonAdminSubGroupCreation bool                  `json:"allowNonAdminSubGroupCreation"`
		LastReportToAdminTimestamp    interface{}           `json:"lastReportToAdminTimestamp"`
		HasCapi                       bool                  `json:"hasCapi"`
		MemberShareGroupHistoryMode   string                `json:"memberShareGroupHistoryMode"`
		Participants                  []newGroupParticipant `json:"participants"`
		PendingParticipants           []interface{}         `json:"pendingParticipants"`
		PastParticipants              []struct {
			Id struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"id"`
			LeaveTs     int    `json:"leaveTs"`
			LeaveReason string `json:"leaveReason"`
		} `json:"pastParticipants"`
		MembershipApprovalRequests []interface{} `json:"membershipApprovalRequests"`
		SubgroupSuggestions        []interface{} `json:"subgroupSuggestions"`
	} `json:"groupMetadata"`
	Id struct {
		Server     string `json:"server"`
		User       string `json:"user"`
		Serialized string `json:"_serialized"`
	} `json:"id"`
	Name           string `json:"name"`
	IsGroup        bool   `json:"isGroup"`
	IsReadOnly     bool   `json:"isReadOnly"`
	UnreadCount    int    `json:"unreadCount"`
	Timestamp      int    `json:"timestamp"`
	Archived       bool   `json:"archived"`
	Pinned         bool   `json:"pinned"`
	IsLocked       bool   `json:"isLocked"`
	IsMuted        bool   `json:"isMuted"`
	MuteExpiration int    `json:"muteExpiration"`
	LastMessage    struct {
		Data struct {
			Id struct {
				FromMe      bool   `json:"fromMe"`
				Remote      string `json:"remote"`
				Id          string `json:"id"`
				Participant struct {
					Server     string `json:"server"`
					User       string `json:"user"`
					Serialized string `json:"_serialized"`
				} `json:"participant"`
				Field5     string `json:"$1"`
				Serialized string `json:"_serialized"`
			} `json:"id"`
			Viewed     bool   `json:"viewed"`
			Body       string `json:"body"`
			Type       string `json:"type"`
			T          int    `json:"t"`
			NotifyName string `json:"notifyName"`
			From       struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"from"`
			To struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"to"`
			Author struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"author"`
			Ack                         int           `json:"ack"`
			Invis                       bool          `json:"invis"`
			Star                        bool          `json:"star"`
			KicNotified                 bool          `json:"kicNotified"`
			IsFromTemplate              bool          `json:"isFromTemplate"`
			IsAdsMedia                  bool          `json:"isAdsMedia"`
			PollInvalidated             bool          `json:"pollInvalidated"`
			IsSentCagPollCreation       bool          `json:"isSentCagPollCreation"`
			LatestEditMsgKey            interface{}   `json:"latestEditMsgKey"`
			LatestEditSenderTimestampMs interface{}   `json:"latestEditSenderTimestampMs"`
			Broadcast                   bool          `json:"broadcast"`
			MentionedJidList            []interface{} `json:"mentionedJidList"`
			GroupMentions               []interface{} `json:"groupMentions"`
			IsEventCanceled             bool          `json:"isEventCanceled"`
			EventInvalidated            bool          `json:"eventInvalidated"`
			IsVcardOverMmsDocument      bool          `json:"isVcardOverMmsDocument"`
			IsForwarded                 bool          `json:"isForwarded"`
			IsQuestion                  bool          `json:"isQuestion"`
			IsSpoiler                   bool          `json:"isSpoiler"`
			QuestionReplyQuotedMessage  interface{}   `json:"questionReplyQuotedMessage"`
			QuestionResponsesCount      int           `json:"questionResponsesCount"`
			ReadQuestionResponsesCount  int           `json:"readQuestionResponsesCount"`
			ForwardsCount               int           `json:"forwardsCount"`
			Labels                      []interface{} `json:"labels"`
			HasReaction                 bool          `json:"hasReaction"`
			NewsletterAdminProfile      interface{}   `json:"newsletterAdminProfile"`
			ExpiredTimestamp            interface{}   `json:"expiredTimestamp"`
			ViewMode                    string        `json:"viewMode"`
			MessageSecret               struct {
				Field1  int `json:"0"`
				Field2  int `json:"1"`
				Field3  int `json:"2"`
				Field4  int `json:"3"`
				Field5  int `json:"4"`
				Field6  int `json:"5"`
				Field7  int `json:"6"`
				Field8  int `json:"7"`
				Field9  int `json:"8"`
				Field10 int `json:"9"`
				Field11 int `json:"10"`
				Field12 int `json:"11"`
				Field13 int `json:"12"`
				Field14 int `json:"13"`
				Field15 int `json:"14"`
				Field16 int `json:"15"`
				Field17 int `json:"16"`
				Field18 int `json:"17"`
				Field19 int `json:"18"`
				Field20 int `json:"19"`
				Field21 int `json:"20"`
				Field22 int `json:"21"`
				Field23 int `json:"22"`
				Field24 int `json:"23"`
				Field25 int `json:"24"`
				Field26 int `json:"25"`
				Field27 int `json:"26"`
				Field28 int `json:"27"`
				Field29 int `json:"28"`
				Field30 int `json:"29"`
				Field31 int `json:"30"`
				Field32 int `json:"31"`
			} `json:"messageSecret"`
			ProductHeaderImageRejected   bool        `json:"productHeaderImageRejected"`
			LastPlaybackProgress         int         `json:"lastPlaybackProgress"`
			IsDynamicReplyButtonsMsg     bool        `json:"isDynamicReplyButtonsMsg"`
			IsCarouselCard               bool        `json:"isCarouselCard"`
			ParentMsgId                  interface{} `json:"parentMsgId"`
			CallSilenceReason            interface{} `json:"callSilenceReason"`
			IsVideoCall                  bool        `json:"isVideoCall"`
			CallDuration                 interface{} `json:"callDuration"`
			CallCreator                  interface{} `json:"callCreator"`
			CallParticipants             interface{} `json:"callParticipants"`
			IsCallLink                   interface{} `json:"isCallLink"`
			CallLinkToken                interface{} `json:"callLinkToken"`
			TerminatedByDeviceSwitch     interface{} `json:"terminatedByDeviceSwitch"`
			SelfOtherDeviceConnected     interface{} `json:"selfOtherDeviceConnected"`
			BytesSent                    interface{} `json:"bytesSent"`
			BytesReceived                interface{} `json:"bytesReceived"`
			IsMdHistoryMsg               bool        `json:"isMdHistoryMsg"`
			StickerSentTs                int         `json:"stickerSentTs"`
			IsAvatar                     bool        `json:"isAvatar"`
			LastUpdateFromServerTs       int         `json:"lastUpdateFromServerTs"`
			InvokedBotWid                interface{} `json:"invokedBotWid"`
			BotTargetSenderJid           interface{} `json:"botTargetSenderJid"`
			MetaFrom                     interface{} `json:"metaFrom"`
			BizBotType                   interface{} `json:"bizBotType"`
			BotDeepLinkToken             interface{} `json:"botDeepLinkToken"`
			AiMediaCollectionInfo        interface{} `json:"aiMediaCollectionInfo"`
			ExpectedImageCount           interface{} `json:"expectedImageCount"`
			ExpectedVideoCount           interface{} `json:"expectedVideoCount"`
			BotResponseTargetId          interface{} `json:"botResponseTargetId"`
			BotPluginType                interface{} `json:"botPluginType"`
			BotPluginReferenceIndex      interface{} `json:"botPluginReferenceIndex"`
			BotPluginSearchProvider      interface{} `json:"botPluginSearchProvider"`
			BotPluginSearchUrl           interface{} `json:"botPluginSearchUrl"`
			BotPluginSearchQuery         interface{} `json:"botPluginSearchQuery"`
			BotPluginMaybeParent         bool        `json:"botPluginMaybeParent"`
			BotReelPluginThumbnailCdnUrl interface{} `json:"botReelPluginThumbnailCdnUrl"`
			BotMessageDisclaimerText     interface{} `json:"botMessageDisclaimerText"`
			BotSessionTransparencyType   interface{} `json:"botSessionTransparencyType"`
			BotMsgBodyType               interface{} `json:"botMsgBodyType"`
			BotModeSelection             interface{} `json:"botModeSelection"`
			BotModeOverride              interface{} `json:"botModeOverride"`
			ReportingTokenInfo           struct {
				ReportingTag struct {
				} `json:"reportingTag"`
				Version int `json:"version"`
			} `json:"reportingTokenInfo"`
			HatchMetadataSync                     interface{}   `json:"hatchMetadataSync"`
			RequiresDirectConnection              bool          `json:"requiresDirectConnection"`
			BizContentPlaceholderType             interface{}   `json:"bizContentPlaceholderType"`
			HostedBizEncStateMismatch             bool          `json:"hostedBizEncStateMismatch"`
			SenderOrRecipientAccountTypeHosted    bool          `json:"senderOrRecipientAccountTypeHosted"`
			PlaceholderCreatedWhenAccountIsHosted bool          `json:"placeholderCreatedWhenAccountIsHosted"`
			GroupHistoryBundleMessageKey          interface{}   `json:"groupHistoryBundleMessageKey"`
			GroupHistoryBundleMetadata            interface{}   `json:"groupHistoryBundleMetadata"`
			GroupHistoryIndividualMessageInfo     interface{}   `json:"groupHistoryIndividualMessageInfo"`
			NonJidMentions                        interface{}   `json:"nonJidMentions"`
			Links                                 []interface{} `json:"links"`
		} `json:"_data"`
		Id struct {
			FromMe      bool   `json:"fromMe"`
			Remote      string `json:"remote"`
			Id          string `json:"id"`
			Participant struct {
				Server     string `json:"server"`
				User       string `json:"user"`
				Serialized string `json:"_serialized"`
			} `json:"participant"`
			Field5     string `json:"$1"`
			Serialized string `json:"_serialized"`
		} `json:"id"`
		Ack             int           `json:"ack"`
		HasMedia        bool          `json:"hasMedia"`
		Body            string        `json:"body"`
		Type            string        `json:"type"`
		Timestamp       int           `json:"timestamp"`
		From            string        `json:"from"`
		To              string        `json:"to"`
		Author          string        `json:"author"`
		DeviceType      string        `json:"deviceType"`
		IsForwarded     bool          `json:"isForwarded"`
		ForwardingScore int           `json:"forwardingScore"`
		IsStatus        bool          `json:"isStatus"`
		IsStarred       bool          `json:"isStarred"`
		Broadcast       bool          `json:"broadcast"`
		FromMe          bool          `json:"fromMe"`
		HasQuotedMsg    bool          `json:"hasQuotedMsg"`
		HasReaction     bool          `json:"hasReaction"`
		VCards          []interface{} `json:"vCards"`
		MentionedIds    []interface{} `json:"mentionedIds"`
		GroupMentions   []interface{} `json:"groupMentions"`
		IsGif           bool          `json:"isGif"`
		Links           []interface{} `json:"links"`
	} `json:"lastMessage"`
}

type newGroupParticipant struct {
	Id struct {
		Server     string `json:"server"`
		User       string `json:"user"`
		Serialized string `json:"_serialized"`
	} `json:"id"`
	IsAdmin               bool        `json:"isAdmin"`
	IsSuperAdmin          bool        `json:"isSuperAdmin"`
	JoinTime              interface{} `json:"joinTime"`
	GroupHistorySentState int         `json:"groupHistorySentState"`
}

type newContactInfo struct {
	Id          string        `json:"id"`
	Number      string        `json:"number"`
	IsBusiness  bool          `json:"isBusiness"`
	Labels      []interface{} `json:"labels"`
	Name        string        `json:"name"`
	Pushname    string        `json:"pushname"`
	ShortName   string        `json:"shortName"`
	StatusMute  bool          `json:"statusMute"`
	Type        string        `json:"type"`
	IsMe        bool          `json:"isMe"`
	IsUser      bool          `json:"isUser"`
	IsGroup     bool          `json:"isGroup"`
	IsWAContact bool          `json:"isWAContact"`
	IsMyContact bool          `json:"isMyContact"`
	IsBlocked   bool          `json:"isBlocked"`
}
