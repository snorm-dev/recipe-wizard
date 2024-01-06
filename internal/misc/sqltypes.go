package misc

import "database/sql"

func SqlNullStringFromOkString(s string, ok bool) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  ok,
	}
}
