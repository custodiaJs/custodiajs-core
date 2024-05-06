// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

type CONSOLE_TEXT string

const (
	VALIDATE_INCOMMING_REMOTE_FUNCTION_CALL_REQUEST_FROM CONSOLE_TEXT = "validate incomming remote function call request from '%s'\n"
	DETERMINE_THE_SCRIPT_CONTAINER                       CONSOLE_TEXT = "determine the script container '%s'\n"
	DETERMINE_THE_FUNCTION                               CONSOLE_TEXT = "&[%s]: determine the function '%s'\n"
	RPC_CALL_DONE_RESPONSE                               CONSOLE_TEXT = "done, response size = %s bytes\n"
)
