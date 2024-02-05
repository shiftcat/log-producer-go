package message

type LogType string

const (
	REQ = LogType("REQ")
	RES = LogType("RES")
)

type ServiceCode string

const (
	AP064210 = ServiceCode("AP064210")

	OU089012 = ServiceCode("OU089012")

	PG004040 = ServiceCode("PG004040")

	SP016739 = ServiceCode("SP016739")

	AJ054580 = ServiceCode("AJ054580")
)

type ResponseType string

const (
	I = ResponseType("I")
	S = ResponseType("S")
	E = ResponseType("E")
)
