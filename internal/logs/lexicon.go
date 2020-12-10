package logs

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Lexicon() error {
	printCacheAccessLexion()
	printCacheErrorLexion()
	printNginxAccessLexion()
	printNginxErrorLexion()
	printFPMAccessLexion()
	printFPMErrorLexion()
	printDeployLexion()
	printCronLexion()

	return nil
}

func printCacheAccessLexion() {
	fmt.Println()
	fmt.Println("cache.access.log")

	data := [][]string{
		{"addr", "IP address of the remote user, or an upstream proxy such as a CDN cache node if the users visitors IP address couldn’t be determined"},
		{"rqid", "A unique request ID that is present in the cache, nginx, and fpm logs and also present in the X-Ironstar-Request-ID header"},
		{"stat", "The response HTTP status, such as 200 or 404"},
		{"meth", "The request method, such as GET or POST"},
		{"ruri", "The request URI, being the hostname, page, and any parameters"},
		{"cach", "Indicates if the request was served from the cache (HIT) or not found in cache (MISS) or bypassed (BYPASS), such as if the request belongs to a logged in user"},
		{"rqtm", "Request time shows how long it took for the entire require to be resolved"},
		{"bsnt", "Bytes sent shows the size of the body of the response"},
		{"xfor", "Displays the content of the X-Forwarded-For header"},
		{"cfra", "Displays the content of the CF-Ray-ID header for Cloudflare-enabled environments"},
		{"cfci", "Displays the content of the CF-Connecting-IP header for Cloudflare-enabled environments"},
		{"agnt", "The User Agent"},
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
		{"addr", "IP address of the remote user, or an upstream proxy such as a CDN cache node if the users visitors IP address couldn’t be determined"},
		{"rqid", "A unique request ID that is present in the cache, nginx, and fpm logs and also present in the x-ironstar-request-id header"},
		{"stat", "The response HTTP status, such as 200 or 404"},
		{"meth", "The request method, such as GET or POST"},
		{"ruri", "The request URI, being the hostname, page, and any parameters"},
		{"rqtm", "Request time shows how long it took for the entire require to be resolved"},
		{"usrt", "Upstream response time show show long PHP or NodeJS took to resolve the request. Will be 0 if request was served directly from disk"},
		{"bsnt", "Bytes sent shows the size of the body of the response"},
		{"ddch", "Displays the content of the x-drupal-dynamic-cacheheader (if set)"},
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
		{"mesg", "The body of the error message"},
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
		{"addr", "IP address of the remote user, or an upstream proxy such as a CDN cache node if the users visitors IP address couldn’t be determined"},
		{"rqid", "A unique request ID that is present in the cache, nginx, and fpm logs and also present in the x-ironstar-request-id header"},
		{"stat", "The response HTTP status, such as 200 or 404"},
		{"meth", "The request method, such as GET or POST"},
		{"ruri", "The request URI, being the hostname, page, and any parameters"},
		{"user", "Displays the content of the X-User header, if set, which you can use to track which user made the request. If set, this header is removed before the response to sent to the user. "},
		{"plen", "Displays the size of the POST body, if the request method is POST "},
		{"dura", "The duration, in seconds, that PHP took to resolve the request"},
		{"pmem", "The amount of memory, in megabytes that PHP used to resolve the request"},
		{"tcpu", "The percentage of CPU which was used to resolve the request"},
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
		{"mesg", "The body of the error message"},
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
		{"mesg", "The body of the deploy message"},
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
		{"levl", "The log level of this log entry"},
		{"mesg", "The log message, or in the case of the results of a cron execution, one line of the resulting output"},
		{"chan", "The output channel, which will be either stdout or stderr"},
		{"iter", "The iteration of this command since the Manager Instance was restarted"},
		{"comm", "The command that was executed (note that there may be multiple lines for the same command, one for each line of output from the command)"},
		{"posi", "The position of this command in the crontab file calculated at runtime"},
		{"schd", "The schedule that was used for this cron execution"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Purpose"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(true)

	table.AppendBulk(data)
	table.Render()
}
