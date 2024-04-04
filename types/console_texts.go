package types

type CONSOLE_TEXT string

const (
	VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM CONSOLE_TEXT = "validate incomming remote function call request from '%s'\n"
	DETERMINE_THE_SCRIPT_CONTAINER                       CONSOLE_TEXT = "determine the script container '%s'\n"
	DETERMINE_THE_FUNCTION                               CONSOLE_TEXT = "&[%s]: determine the function '%s'\n"
	RPC_CALL_DONE_RESPONSE                               CONSOLE_TEXT = "done, response size = %s bytes\n"
)
