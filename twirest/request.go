package twirest

// Request a list of the account resources
type Accounts struct {
	FriendlyName string `FriendlyName=`
	Status       string `Status=`
}

// Account resource information for a single account
type Account struct {
	Sid string
}

// Request list of calls made to and from account
type Calls struct {
	To              string `To=`
	From            string `From=`
	Status          string `Status=`
	StartTime       string `StartTime=`
	StartTimeBefore string `StartTime<=`
	StartTimeAfter  string `StartTime>=`
	ParentCallSid   string `ParentCallSid=`
}

// Request call information about a single call
type Call struct {
	Sid           string // CallSid
	Recordings    bool
	Notifications bool
}

// Request to make a phone call
type MakeCall struct {
	From                 string `From=`
	To                   string `To=`
	Url                  string `Url=`
	ApplicationSid       string `ApplicationSid=`
	Method               string `Method=`
	FallbackUrl          string `FallbackUrl=`
	FallbackMethod       string `FallbackMethod=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
	SendDigits           string `SendDigits=`
	IfMachine            string `IfMachine=`
	Timeout              string `Timeout=`
	Record               string `Record=`
	SipAuthUsername      string `SipAuthUsername=`
	SipAuthPassword      string `SipAuthPassword=`
}

// Request to modify call in queue/progress
type ModifyCall struct {
	Sid                  string
	Url                  string `Url=`
	Method               string `Method=`
	Status               string `Status=`
	FallbackUrl          string `FallbackUrl=`
	FallbackMethod       string `FallbackMethod=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
}

// List conferences within an account
type Conferences struct {
	Status            string `Status=`
	FriendlyName      string `FriendlyName=`
	DateCreated       string `DateCreated=`
	DateCreatedBefore string `DateCreated<=`
	DateCreatedAfter  string `DateCreated>=`
	DateUpdated       string `DateUpdated=`
	DateUpdatedBefore string `DateUpdated<=`
	DateUpdatedAfter  string `DateUpdated>=`
}

// Resource for individual conference instance
type Conference struct {
	Sid string
}

// Request list of participants in a conference
type Participants struct {
	Sid   string // Conference Sid
	Muted string `Muted=`
}

// Resource about single conference participant
type Participant struct {
	Sid     string // Conference Sid
	CallSid string
}

// Remove a participant from a conference
type DeleteParticipant struct {
	Sid     string // Conference Sid
	CallSid string
}

// Request to change the status of a participant
type UpdateParticipant struct {
	Sid     string // Conference Sid
	CallSid string
	Muted   string `Muted=`
}

// Messages struct for request of list of messages
type Messages struct {
	To             string `To=`
	From           string `From=`
	DateSent       string `DateSent=`
	DateSentBefore string `DateSent<=`
	DateSentAfter  string `DateSent>=`
}

// Message struct for request of single message
type Message struct {
	Sid      string // MessageSid
	Media    bool
	MediaSid string
}

// Message struct for request to send a message
type SendMessage struct {
	Text           string `Body=`
	MediaUrl       string `MediaUrl=`
	From           string `From=`
	To             string `To=`
	ApplicationSid string `ApplicationSid=`
	StatusCallback string `StatusCallback=`
}

// Notifications struct for request of a possible list of notifications
type Notifications struct {
	Log           string `Log=`
	MsgDate       string `MessageDate=`
	MsgDateBefore string `MessageDate<=`
	MsgDateAfter  string `MessageDate>=`
}

// Notification struct for request of a specific notification
type Notification struct {
	Sid string
}

// DeleteNotification struct for removal of a notification
type DeleteNotification struct {
	Sid string
}

// Get outgoing caller IDs
type OutgoingCallerIds struct {
	PhoneNumber  string `PhoneNumber=`
	FriendlyName string `FriendlyName=`
}

// Get outgoing caller ID
type OutgoingCallerId struct {
	Sid string
}

type UpdateOutgoingCallerId struct {
	Sid          string
	FriendlyName string `FriendlyName=`
}

type DeleteOutgoingCallerId struct {
	Sid string
}

type AddOutgoingCallerId struct {
	PhoneNumber          string `PhoneNumber=`
	FriendlyName         string `FriendlyName=`
	CallDelay            string `CallDelay=`
	Extension            string `Extension=`
	StatusCallback       string `StatusCallback=`
	StatusCallbackMethod string `StatusCallbackMethod=`
}

// List recordings resource
type Recordings struct {
	Sid               string `CallSid=`
	DateCreated       string `DateCreated=`
	DateCreatedBefore string `DateCreated<=`
	DateCreatedAfter  string `DateCreated>=`
}

// Request resource for an individual recording
type Recording struct {
	Sid string // RecordingSid
}

// Delete a recording
type DeleteRecording struct {
	Sid string // RecordingSid
}

// Request usage by the account
type UsageRecords struct {
	SubResource string
	Category    string `Category=`
	StartDate   string `StartDate=`
	EndDate     string `EndDate=`
}

// List queues within an account
type Queues struct {
}

// Get resource for an individual Queue instance
type Queue struct {
	Sid string // QueueSid
}

// Create a new queue
type CreateQueue struct {
	FriendlyName string `FriendlyName=`
	MaxSize      string `MaxSize=`
}

// Request to change queue properties
type ChangeQueue struct {
	Sid          string
	FriendlyName string `FriendlyName=`
	MaxSize      string `MaxSize=`
}

// Remove a queue
type DeleteQueue struct {
	Sid string // QueueSid
}

// List members of a queue
type QueueMembers struct {
	Sid string // QueueSid
}

// Request resource for a queue member
type QueueMember struct {
	Sid     string // QueueSid
	CallSid string
	Front   bool
}

// Remove a member from a queue and redirect the member's call to a TwiML site
type DeQueue struct {
	Sid     string // Queue Sid
	CallSid string
	Front   bool
	Url     string `Url=`
	Method  string `Method=`
}
