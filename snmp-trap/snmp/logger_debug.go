//go:build !gosnmp_nodebug
// +build !gosnmp_nodebug

package snmp

func (l *Logger) Print(v ...interface{}) {
    if l.logger != nil {
        l.logger.Print(v...)
    }
}

func (l *Logger) Printf(format string, v ...interface{}) {
    if l.logger != nil {
        l.logger.Printf(format, v...)
    }
}
