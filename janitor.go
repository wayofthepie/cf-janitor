package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "strings"
    "time"
    "github.com/cloudfoundry/cli/cf/api/resources"
    //"github.com/cloudfoundry/cli/cf/terminal"
    "github.com/cloudfoundry/cli/plugin"

)

type JanitorPlugin struct{
    cliConnection plugin.CliConnection
    before 		  *string
}

/*
*	This function must be implemented by any plugin because it is part of the
*	plugin interface defined by the core CLI.
*
*	Run(....) is the entry point when the core CLI is invoking a command defined
*	by a plugin. The first parameter, plugin.CliConnection, is a struct that can
*	be used to invoke cli commands. The second parameter, args, is a slice of
*	strings. args[0] will be the name of the command, and will be followed by
*	any additional arguments a cli user typed in.
*
*	Any error handling should be handled with the plugin itself (this means printing
*	user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
*	1 should the plugin exits nonzero.
 */
func (c *JanitorPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    c.cliConnection = cliConnection

    fs := new(flag.FlagSet)
    c.before = fs.String("before", "", "")
    fs.Parse(args[1:])
    c.execute()

}


func (c *JanitorPlugin) execute() {
    if c.validArgs() {
        if *c.before != "" {
            space, err := c.cliConnection.GetCurrentSpace()
            if err != nil {
                fmt.Println(err.Error())
                return
            }

            var before time.Time
            if *c.before == "now" {
                before = time.Now()
            } else {
                before, err = time.Parse(time.RFC3339, *c.before)
                if err != nil {
                    fmt.Println(err.Error())
                    return
                }
            }
            c.findAppsBefore(space.Guid, before)
        }

    } else {
    }
}

func (c *JanitorPlugin) validArgs() bool {
    return (c.hasFlag(*c.before, "before"))
}

func (s *JanitorPlugin) hasFlag(fl string, name string) (ret bool) {
    if ret = (fl != ""); ret == false {
        return
    }
    return
}

/*
*	This function must be implemented as part of the	plugin interface
*	defined by the core CLI.
*
*	GetMetadata() returns a PluginMetadata struct. The first field, Name,
*	determines the name of the plugin which should generally be without spaces.
*	If there are spaces in the name a user will need to properly quote the name
*	during uninstall otherwise the name will be treated as seperate arguments.
*	The second value is a slice of Command structs. Our slice only contains one
*	Command Struct, but could contain any number of them. The first field Name
*	defines the command `cf basic-plugin-command` once installed into the CLI. The
*	second field, HelpText, is used by the core CLI to display help information
*	to the user in the core commands `cf help`, `cf`, or `cf -h`.
 */
func (c *JanitorPlugin) GetMetadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        Name: "Janitor",
        Version: plugin.VersionType{
            Major: 1,
            Minor: 0,
            Build: 0,
        },
        MinCliVersion: plugin.VersionType{
            Major: 6,
            Minor: 7,
            Build: 0,
        },
        Commands: []plugin.Command{
            plugin.Command{
                Name:     "janitor",
                HelpText: "Janitor command's help text",

                // UsageDetails is optional
                // It is used to show help of usage of each command
                UsageDetails: plugin.Usage{
                    Usage: "janitor",
                },
            },
        },
    }
}


func (c* JanitorPlugin) findAppsBefore(spaceGuid string, before time.Time) {

    appCmd := []string{"curl", "/v2/spaces/" + spaceGuid + "/apps"}
    appsJson, err := c.cliConnection.CliCommandWithoutTerminalOutput(appCmd...)

    if err != nil {
        return
    }

    res := &resources.PaginatedApplicationResources{}
    json.Unmarshal([]byte(strings.Join(appsJson,"")), &res)

    for _,appRes := range res.Resources {
        appName    	   := *appRes.Entity.Name
        lastUploadStr  := fmt.Sprint(appRes.Entity.PackageUpdatedAt)
        lastUploadTime := *appRes.Entity.PackageUpdatedAt
        if lastUploadTime.Before(before) {
            fmt.Println(appName + " last uploaded " + lastUploadStr)
        }
    }
}

/*
* Unlike most Go programs, the `Main()` function will not be used to run all of the
* commands provided in your plugin. Main will be used to initialize the plugin
* process, as well as any dependencies you might require for your
* plugin.
 */
func main() {
    // Any initialization for your plugin can be handled here
    //
    // Note: to run the plugin.Start method, we pass in a pointer to the struct
    // implementing the interface defined at "github.com/cloudfoundry/cli/plugin/plugin.go"
    //
    // Note: The plugin's main() method is invoked at install time to collect
    // metadata. The plugin will exit 0 and the Run([]string) method will not be
    // invoked.
    plugin.Start(new(JanitorPlugin))
    // Plugin code should be written in the Run([]string) method,
    // ensuring the plugin environment is bootstrapped.
}
