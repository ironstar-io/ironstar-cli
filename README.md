# Ironstar CLI

The Ironstar CLI is a user-friendly command-line interface that allows Ironstar customers to easily manage and interact with their hosting environments.

We welcome any feedback or suggestions from our customers to continue improving the CLI. Please contact the Ironstar support team to submit any feedback or feature requests.

## Installation

This command will install `iron` to `/usr/local/bin` for MacOS users

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ironstar-io/ironstar-cli/main/install/macos.sh)"
```

Linux and Windows users can find the binaries on the [releases page](https://github.com/ironstar-io/ironstar-cli/releases). There currently isn't and automated installer for these platforms.

## Basic Usage

`iron login`

Using your credentials and MFA code login to your Ironstar account

`iron sub list`

List the subscriptions available for your account

`iron sub link [sub_name]`

Link a subscription, commands from here will be made in the context of this subscription.

---

For a more detailed command list run `iron -h`

## Upgrading

There is an automatic command to check and upgrade the Ironstar CLI to the latest version

```
iron upgrade
```
