package config

import (
	"flag"
	"os"
	"strings"

	"github.com/blazejsewera/go-test-proxy/colorfmt"
	"github.com/blazejsewera/go-test-proxy/colorfmt/log"
)

const applicationName = "https://github.com/blazejsewera/go-test-proxy"

type Configuration struct {
	Application    string   `json:"application"`
	Port           uint16   `json:"port"`
	Target         string   `json:"target"`
	Color          bool     `json:"color"`
	MockGroups     []string `json:"mockGroups"`
	EnableAllMocks bool     `json:"enableAllMocks"`

	Cfmt *colorfmt.Fmt `json:"-"`
}

func ParseConfig() Configuration {
	config := parseConfig()
	config.Cfmt = colorfmt.New(config.Color, os.Stdout, os.Stderr)
	log.SetFmt(config.Cfmt)
	validateConfig(config)
	return config
}

func parseConfig() Configuration {
	config := Configuration{Application: applicationName}

	target := flag.String("target", "", "the target host address, for example, https://example.com")
	port := flag.Int("port", 8000, "the port on which the proxy server will be running")
	color := flag.Bool("color", false, "use terminal color, requires the terminal emulator to support ANSI color")
	mockGroups := flag.String("mockGroups", "", "comma separated list of mock groups to enable, for example, 'group1,group2'")
	enableAllMocks := flag.Bool("allMocks", false, "enable all mocks; when enabled, mockGroups argument is ignored")
	flag.Parse()

	config.Target = *target

	portUint16 := uint16(*port)
	config.Port = portUint16

	config.Color = *color

	config.MockGroups = parseMockGroups(*mockGroups)
	config.EnableAllMocks = *enableAllMocks

	return config
}

func parseMockGroups(mockGroups string) []string {
	if mockGroups == "" {
		return []string{}
	}
	return strings.Split(mockGroups, ",")
}

func validateConfig(config Configuration) {
	if config.Target == "" {
		log.Fatalln("[FATAL] The target cannot be empty. " +
			"Specify the target server, e.g., https://example.com\n" +
			"Try: ./gotestproxy --target=https://example.com\n" +
			"Run with -h for help.")
	}
}
