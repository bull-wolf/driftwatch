# driftwatch

> Detect config drift between deployed services and their source-of-truth manifests.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your manifest directory and a running cluster (or API endpoint) to scan for drift:

```bash
driftwatch scan --manifests ./manifests/ --target https://api.mycluster.internal
```

Compare a single service against its expected config:

```bash
driftwatch diff --manifest ./manifests/payments-service.yaml --service payments
```

Output formats:

```bash
driftwatch scan --manifests ./manifests/ --output json
driftwatch scan --manifests ./manifests/ --output table   # default
```

### Example Output

```
SERVICE          FIELD              EXPECTED        ACTUAL          STATUS
payments         replicas           3               1               DRIFT
auth-service     image.tag          v1.4.2          v1.3.9          DRIFT
inventory        env.LOG_LEVEL      info            debug           DRIFT
gateway          replicas           2               2               OK
```

Exit code `1` is returned when drift is detected, making it easy to integrate into CI pipelines.

---

## Configuration

`driftwatch` can be configured via a `.driftwatch.yaml` file in the project root:

```yaml
manifests: ./manifests/
target: https://api.mycluster.internal
ignore:
  - "*.timestamp"
  - "metadata.resourceVersion"
```

---

## License

MIT © yourorg