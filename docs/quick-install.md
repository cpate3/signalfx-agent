# SignalFx Smart Agent Quick Install


The SignalFx Smart Agent is a metric agent written in Go for monitoring infrastructure and application services in a variety of environments. It is a successor to our previous [collectd agent](https://github.com/signalfx/collectd), and still uses collectd internally on Linux; any existing Python or C-based collectd plugins will still work without modification. On Windows collectd is not included, but the agent can run python-based collectd plugins without collectd. C-based collectd plugins are not available on Windows.

 - [Concepts](#concepts)
 - [Installation](#installation)


## Concepts

The agent has three main components:

* _Monitors_ that collect metrics from the host and applications

* _Observers_ that discover applications and services running on the host

* a _Writer_ that sends the metrics collected by monitors to SignalFx


### Monitors

Monitors collect metrics from the host system and services.  They are
configured under the `monitors` list in the agent config.  For
application specific monitors, you can define discovery rules in your monitor
configuration. A separate monitor instance is created for each discovered
instance of applications that match a discovery rule. See [Auto
Discovery](./auto-discovery.md) for more information.

Many of the monitors are built around [collectd](https://collectd.org), an open
source third-party monitor, and use it to collect metrics. Some other monitors
do not use collectd. However, either type is configured in the same way.

For a list of supported monitors and their configurations,
see [Monitor Config](./monitor-config.md).

The agent is primarily intended to monitor services/applications running on the
same host as the agent.  This is in keeping with the collectd model.  The main
issue with monitoring services on other hosts is that the `host` dimension that
collectd sets on all metrics will currently get set to the hostname of the
machine that the agent is running on.  This allows everything to have a
consistent `host` dimension so that metrics can be matched to a specific
machine during metric analysis.

### Observers

Observers watch the various environments that we support to discover running
services and automatically configure the agent to send metrics for those
services.

For a list of supported observers and their configurations,
see [Observer Config](./observer-config.md).

### Writer
The writer collects metrics emitted by the configured monitors and sends them
to SignalFx on a regular basis.  There are a few things that can be
[configured](./config-schema.md#writer) in the writer, but this is generally only necessary if you have a very large number of metrics flowing through a single agent.


## Installation

The instructions below are for a command-line installation of a single agent. For other installation options, including bulk deployments, see the Advanced Installation in [Advanced Installation Options](./advanced-install-options.md).

### __Get started with Smart Agent using the 3 steps below.__

_Note: if you have previously configured another metric collection agent on your host such as collectd, uninstall or disable that agent before installing the SignalFx Smart Agent._

### Step 1. Install SignalFx Smart Agent on Single Host 

__Linux:__ Dependencies are completely bundled along with the agent, including a Java JRE runtime and a Python runtime, so there are no additional dependencies required. The agent works on any modern Linux distribution (kernel version 2.6+).

If you are not installing from the tile on the Integrations page:

- Get your API_TOKEN from: __Organization Settings => Access Token__ tab in the SignalFx application. 

- Determine YOUR\_SIGNAL_FX_REALM from: your profile page in the SignalFx web application.

To install the Smart Agent on a single Linux host, enter:

```sh
curl -sSL https://dl.signalfx.com/signalfx-agent.sh > /tmp/signalfx-agent.sh
sudo sh /tmp/signalfx-agent.sh --realm YOUR_SIGNALFX_REALM YOUR_SIGNALFX_API_TOKEN
````

__Windows:__ Ensure that the following dependencies are installed:

[.Net Framework 3.5](https://docs.microsoft.com/en-us/dotnet/framework/install/dotnet-35-windows-10) (Windows 8+)

[Visual C++ Compiler for Python 2.7](https://www.microsoft.com/EN-US/DOWNLOAD/DETAILS.ASPX?ID=44266)

* Get your API\_TOKEN from:  __Organization Settings => Access Token__ tab in the SignalFx application.

* Determine YOUR\_SIGNAL\_FX_REALM from: your profile page in the SignalFx web application.

To install the Smart Agent on a single Windows host, enter:

```sh
& {Set-ExecutionPolicy Bypass -Scope Process -Force; $script = ((New-Object System.Net.WebClient).DownloadString('https://dl.signalfx.com/signalfx-agent.ps1')); $params = @{access_token = "YOUR_SIGNALFX_API_TOKEN"};; 
ingest_url = "https://ingest.YOUR_SIGNALFX_REALM.signalfx.com"; api_url = "https://api.YOUR_SIGNALFX_REALM.signalfx.com"}; Invoke-Command -ScriptBlock ([scriptblock]::Create(". {$script} $(&{$args} @params)"))}
```

The agent will be installed as a Windows service and will log to the Windows Event Log.


### Step 2. Confirm your Installation 

To confirm the SignalFx Smart Agent installation is functional on either platform, enter:

```sh
sudo signalfx-agent status
````

The response you will see is similar to the one below:

```sh
SignalFx Agent version:           4.7.6
Agent uptime:                     8m44s
Observers active:                 host
Active Monitors:                  16
Configured Monitors:              33
Discovered Endpoint Count:        6
Bad Monitor Config:               None
Global Dimensions:                {host: my-host-1}
Datapoints sent (last minute):    1614
Events Sent (last minute):        0
Trace Spans Sent (last minute):   0
````

To verify the installation, you can run the following commands:

```sh
signalfx-agent status config - show resolved config in use by agent
signalfx-agent status endpoints - show discovered endpoints
signalfx-agent status monitors - show active monitors
signalfx-agent status all - show everything
````

#### Troubleshoot any discrepancies in the Installation 

##### Realm 

By default, the Smart Agent will send data to the us0 realm. If you are not in this realm, you will need to explicitly set the signalFxRealm option in your config like this:


```sh
signalFxRealm: YOUR_SIGNALFX_REALM
```


To determine if you are in a different realm and need to explicitly set the endpoints, check your profile page in the SignalFx web application.

_Configure your endpoints_

If you want to explicitly set the ingest, API server, and trace endpoint URLs, you can set them individually like so:


```sh
ingestUrl: "https://ingest.YOUR_SIGNALFX_REALM.signalfx.com"
apiUrl: "https://api.YOUR_SIGNALFX_REALM.signalfx.com"
traceEndpointUrl: "https://ingest.YOUR_SIGNALFX_REALM.signalfx.com/v1/trace"
````

This will default to the endpoints for the realm configured in signalFxRealm if not set.

To troubleshoot your installation further, check the FAQ about troubleshooting [here](https://docs.signalfx.com/en/latest/integrations/agent/faq.html).


### Step 3. Login to SignalFx and discover your data displays. 

Installation is complete.

To continue your exploration of SignalFx Smart Agent capabilities, see [Advanced Installation Options](https://docs.signalfx.com/en/latest/integrations/agent/advanced-install-options.html). 

