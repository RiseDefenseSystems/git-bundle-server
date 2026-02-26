# Git Bundle Server

[bundle-uris]: https://git-scm.com/docs/bundle-uri
[codeowners]: CODEOWNERS
[contributing]: CONTRIBUTING.md
[license]: LICENSE
[support]: SUPPORT.md

## Background

The Git Bundle Server is a self-hosted bundle server for Git's
[bundle URI feature][bundle-uris].

## Usage

### Managing Repositories

The following command-line interface allows you to manage which repositories are
being managed by the bundle server.

* `git-bundle-server init <url> [<route>]`: Initialize a repository by cloning a
  bare repo from `<url>`. If `<route>` is specified, then it is the bundle
  server route to find the data for this repository. Otherwise, the route is
  inferred from `<url>` by removing the domain name.
* `git-bundle-server update <route>`: For the repository at the specified `<route>`,
  fetch the latest content from the remote and create a new set of bundles.
* `git-bundle-server stop <route>`: Stop computing bundles or serving content
  for the repository at the specified `<route>`.
* `git-bundle-server start <route>`: Start computing bundles and serving content
  for the repository at the specified `<route>`.
* `git-bundle-server delete <route>`: Remove the configuration for the given
  `<route>` and delete its repository data.

### Web Server Management

Independent of the management of the individual repositories hosted by the
server, you can manage the web server process itself using these commands:

* `git-bundle-server web-server start`: Start the web server process.
* `git-bundle-server web-server stop`: Stop the web server process.

### Additional Resources

Detailed guides to more complex administration tasks or user workflows can be
found in the [`docs/tutorials`](./docs/tutorials/) directory of this repository.

## Local Development

### Building

> To avoid environment issues building and executing Go code, we recommend that
> you clone inside the `src` directory of your `GOROOT`.

In the root of your cloned repository, you can build the `git-bundle-server` and
`git-bundle-web-server` executables a few ways.

The first is to use GNU Make; from the root of the repository, simply run:
