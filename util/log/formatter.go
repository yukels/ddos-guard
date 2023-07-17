package log

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yukels/util/collection"
)

var (
	ignore = []string{"id", "ctx", "error", "func", "file", "line", "caller"}
)

type formatter struct {
}

func (f *formatter) Format(entry *log.Entry) ([]byte, error) {
	prefix := ""
	if err, ok := entry.Data["error"]; ok {
		prefix += fmt.Sprintf("[Error: %s] ", err)
	}
	where := fmt.Sprintf("[%s:%d - %s()]", entry.Data["file"], entry.Data["line"], entry.Data["func"])

	message := fmt.Sprintf("[%s] ~%s~ %s%s\t%s\t%s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"),
		strings.ToUpper(entry.Level.String()),
		prefix,
		entry.Message,
		fields(entry),
		where)

	return []byte(message), nil
}

func fields(entry *log.Entry) string {
	var fields []string
	for key, value := range entry.Data {
		if !collection.Contains(ignore, key) {
			fields = append(fields, fmt.Sprintf("[%s=%s]", key, value))
		}
	}
	return strings.Join(fields, " ")
}
