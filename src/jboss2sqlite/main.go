package main

import "flag"

var (
	mode         = flag.String("mode", "import", "Mode of operation, (import, daemon), import works on old Data, daemon pulls periodically from the SITE")
	site         = flag.String("site", "TST", "The site that the parsed html files belong to")
	pollInterval = flag.Int("interval", 20, "Polling interval")
	sqliteDB     = flag.String("db", "jvm-ram2.db", "Database for Sqlite implementation")
)

func main() {
	flag.Parse()

	if *mode == "import" {
		importHtml()
	}
}
