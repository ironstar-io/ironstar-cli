package logs

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Lexicon(args []string) error {
	if len(args) == 0 {
		printEntireLexicon()

		return nil
	}

	printSpecificLexicon(args[0])

	return nil
}

func printSpecificLexicon(field string) {
	switch field {
	case "cache.access.log":
		printCacheAccessLexion()
	case "cache.error.log":
		printCacheErrorLexion()
	case "nginx.access.log":
		printNginxAccessLexion()
	case "nginx.error.log":
		printNginxErrorLexion()
	case "fpm.access.log":
		printFPMAccessLexion()
	case "fpm.error.log":
		printFPMErrorLexion()
	case "deploy.log":
		printDeployLexion()
	case "cron.log":
		printCronLexion()
	default:
		if termLexicon[field] == "" {
			printEntireLexicon()

			return
		}

		fmt.Println(field)
		fmt.Println(termLexicon[field])
	}
}

func printEntireLexicon() {
	printCacheAccessLexion()
	printCacheErrorLexion()
	printNginxAccessLexion()
	printNginxErrorLexion()
	printFPMAccessLexion()
	printFPMErrorLexion()
	printDeployLexion()
	printCronLexion()
}

var termLexicon = map[string]string{
	"addr": "IP address of the remote user, or an upstream proxy such as a CDN cache node if the users visitors IP address couldn’t be determined",
	"rqid": "A unique request ID that is present in the cache, nginx, and fpm logs and also present in the X-Ironstar-Request-ID header",
	"stat": "The response HTTP status, such as 200 or 404",
	"meth": "The request method, such as GET or POST",
	"ruri": "The request URI, being the hostname, page, and any parameters",
	"cach": "Indicates if the request was served from the cache (HIT) or not found in cache (MISS) or bypassed (BYPASS), such as if the request belongs to a logged in user",
	"rqtm": "Request time shows how long it took for the entire require to be resolved",
	"bsnt": "Bytes sent shows the size of the body of the response",
	"xfor": "Displays the content of the X-Forwarded-For header",
	"cfra": "Displays the content of the CF-Ray-ID header for Cloudflare-enabled environments",
	"cfci": "Displays the content of the CF-Connecting-IP header for Cloudflare-enabled environments",
	"agnt": "The User Agent",
	"mesg": "The message body",
	"usrt": "Upstream response time show show long PHP or NodeJS took to resolve the request. Will be 0 if request was served directly from disk",
	"ddch": "Displays the content of the x-drupal-dynamic-cacheheader (if set)",
	"user": "Displays the content of the X-User header, if set, which you can use to track which user made the request. If set, this header is removed before the response to sent to the user. ",
	"plen": "Displays the size of the POST body, if the request method is POST ",
	"dura": "The duration, in seconds, that PHP took to resolve the request",
	"pmem": "The amount of memory, in megabytes that PHP used to resolve the request",
	"tcpu": "The percentage of CPU which was used to resolve the request",
	"levl": "The log level of this log entry",
	"chan": "The output channel, which will be either stdout or stderr",
	"iter": "The iteration of this command since the Manager Instance was restarted",
	"comm": "The command that was executed (note that there may be multiple lines for the same command, one for each line of output from the command)",
	"posi": "The position of this command in the crontab file calculated at runtime",
	"schd": "The schedule that was used for this cron execution",
}

func printCacheAccessLexion() {
	fmt.Println()
	fmt.Println("cache.access.log")

	data := [][]string{
		{"addr", termLexicon["addr"]},
		{"rqid", termLexicon["rqid"]},
		{"stat", termLexicon["stat"]},
		{"meth", termLexicon["meth"]},
		{"ruri", termLexicon["ruri"]},
		{"cach", termLexicon["cach"]},
		{"rqtm", termLexicon["rqtm"]},
		{"bsnt", termLexicon["bsnt"]},
		{"xfor", termLexicon["xfor"]},
		{"cfra", termLexicon["cfra"]},
		{"cfci", termLexicon["cfci"]},
		{"agnt", termLexicon["agnt"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printCacheErrorLexion() {
	fmt.Println()
	fmt.Println("cache.error.log")

	data := [][]string{
		{"mesg", "The body of the error message"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printNginxAccessLexion() {
	fmt.Println()
	fmt.Println("nginx.access.log")

	data := [][]string{
		{"addr", termLexicon["addr"]},
		{"rqid", termLexicon["rqid"]},
		{"stat", termLexicon["stat"]},
		{"meth", termLexicon["meth"]},
		{"ruri", termLexicon["ruri"]},
		{"rqtm", termLexicon["rqtm"]},
		{"usrt", termLexicon["usrt"]},
		{"bsnt", termLexicon["bsnt"]},
		{"ddch", termLexicon["ddch"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printNginxErrorLexion() {
	fmt.Println()
	fmt.Println("nginx.error.log")

	data := [][]string{
		{"mesg", termLexicon["mesg"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printFPMAccessLexion() {
	fmt.Println()
	fmt.Println("fpm.access.log")

	data := [][]string{
		{"addr", termLexicon["addr"]},
		{"rqid", termLexicon["rqid"]},
		{"stat", termLexicon["stat"]},
		{"meth", termLexicon["meth"]},
		{"ruri", termLexicon["ruri"]},
		{"user", termLexicon["user"]},
		{"plen", termLexicon["plen"]},
		{"dura", termLexicon["dura"]},
		{"pmem", termLexicon["pmem"]},
		{"tcpu", termLexicon["tcpu"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printFPMErrorLexion() {
	fmt.Println()
	fmt.Println("fpm.error.log")

	data := [][]string{
		{"mesg", termLexicon["mesg"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printDeployLexion() {
	fmt.Println()
	fmt.Println("deploy.log")

	data := [][]string{
		{"mesg", termLexicon["mesg"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}

func printCronLexion() {
	fmt.Println()
	fmt.Println("cron.log")

	data := [][]string{
		{"levl", termLexicon["levl"]},
		{"mesg", termLexicon["mesg"]},
		{"chan", termLexicon["chan"]},
		{"iter", termLexicon["iter"]},
		{"comm", termLexicon["comm"]},
		{"posi", termLexicon["posi"]},
		{"schd", termLexicon["schd"]},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}
