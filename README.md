manual_gcp_autoscaler
====


It is a tool to change the minimum and maximum number of instances under [GCP Managed Instance Group](https://cloud.google.com/compute/docs/instance-groups/).

# Installation

```console
$ go get -u github.com/ohsawa0515/manual_gcp_autoscaler
$ cd manual_gcp_autoscaler
$ make dep
```

# Build

```console
$ make build
$ ls  ./bin/*/manual-gcp-autoscaler
bin/darwin_amd64/manual-gcp-autoscaler bin/linux_amd64/manual-gcp-autoscaler
```

# Set GCP credentials

If you can use service account, please set `GOOGLE_APPLICATION_CREDENTIALS` environment variable.

```console
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/<SERVICE_ACCOUNT>.json
```

# Usage

```console
$ manual-gcp-autoscaler --help
Usage of ./bin/darwin_amd64/manual-gcp-autoscaler:
  -max int
        The maximum number of instances that the autoscaler can scale up to.
  -mig string
        Name of managed instance group.
  -min int
        The minimum number of replicas that the autoscaler can scale down to. This cannot be less than 0.
  -project string
        Project ID.
  -region string
        Region of the managed regional instance group. e.g. asia-northeast1
  -zone string
        Zone of the managed instance group. e.g. asia-northeast1-a
```

If you want to operate the regional managed instance group you need to specify the `-region` option.

```console
$ manual_gcp_autoscaler --project sample-project -region asia-east1 -mig sample-instance-group -min 5 -max 10
```

If you want to operate the managed instance group based on zone you need to specify the `-zone` option.

```console
$ manual_gcp_autoscaler --project sample-project -zone asia-east1-a -mig sample-instance-group-zone-a -min 5 -max 10
```

# Contribution

1. Fork ([https://github.com/ohsawa0515/manual_gcp_autoscaler/fork](https://github.com/ohsawa0515/manual_gcp_autoscaler/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

# License

See [LICENSE](https://github.com/ohsawa0515/manual_gcp_autoscaler/blob/master/LICENSE).

# Author

Shuichi Ohsawa (@ohsawa0515)
