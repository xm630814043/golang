snmptrap -v 3 -a SHA -A trap_net_works -x AES -X net_works_trap -l authPriv -u pcb.snmpv3 -e 0x80001f8880315de44d53ce8394 127.0.0.1 "" linkUp.0 1.3.6.4.5.6 s "switch is off"

#snmptrap -v 2c -c public 127.0.0.1 '' SNMPv2-MIB::system SNMPv2-MIB::sysDescr.0 s "red laptop" SNMPv2-MIB::sysServices.0 i "5" SNMPv2-MIB::sysObjectID o "1.3.6.1.4.1.2.3.4.5"

#snmptrap -v 2c -c  public 127.0.0.1 "" .1.3.6.1.4.1.2021.251.1 sysLocation.0 s "this is test"