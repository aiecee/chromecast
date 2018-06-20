package net

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/micro/mdns"
)

type DNSEntry struct {
	AddressV4  net.IP
	AddressV6  net.IP
	Port       int
	Name       string
	Host       string
	UUID       string
	Device     string
	Status     string
	DeviceName string
	InfoFields map[string]string
}

func GetAllEntries() []DNSEntry {
	entriesChannel := make(chan *mdns.ServiceEntry, 20)
	resultChannel := make(chan DNSEntry)
	entries := make([]DNSEntry, 0)
	go func() {
		for entry := range entriesChannel {
			parsedFields := parseInfoFields(entry.InfoFields)
			resultChannel <- DNSEntry{
				AddressV4:  entry.AddrV4,
				AddressV6:  entry.AddrV6,
				Port:       entry.Port,
				Name:       entry.Name,
				Host:       entry.Host,
				UUID:       parsedFields["id"],
				Device:     parsedFields["md"],
				DeviceName: parsedFields["fn"],
				Status:     parsedFields["rs"],
				InfoFields: parsedFields,
			}
		}
		close(resultChannel)
	}()

	err := mdns.Query(&mdns.QueryParam{
		Service: "_googlecast._tcp",
		Domain:  "local",
		Timeout: time.Second * 10,
		Entries: entriesChannel,
	})
	if err != nil {
		log.Fatal(err)
	}
	close(entriesChannel)
	for entry := range resultChannel {
		entries = append(entries, entry)
	}
	return entries
}

func parseInfoFields(infoFields []string) map[string]string {
	fields := make(map[string]string, len(infoFields))
	for _, field := range infoFields {
		split := strings.Split(field, "=")
		if len(split) != 2 {
			continue
		}
		fields[split[0]] = split[1]
	}
	return fields
}
