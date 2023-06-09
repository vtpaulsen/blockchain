# Systemtester with docker

## For testing-framework developers:

Below is a diagram of how the systemtester is constructed.

A setup is given as input to the systemtester, which has the responsibility of giving functionality for the testers (assertions etc). The docker layer uses topsort to topologically sort the peers from the setup. The responsibility of the docker layer is to handle I/O to docker.

```text
┌────────┐
│ setups │
└─┬──────┘
  │
  ▼
┌──────────────┐
│ systemtester ├───┐
└─┬────────────┘   │
  │   ▲            │
  ▼   │            │
┌─────┴──┐         │
│ docker │         │
└─┬───┬──┘         │
  │   │            ▼
  │   └────────────┐
  ▼                ▼
┌─────────┐    ┌────────┐
│ topsort │    │ helper │
└─────────┘    └────────┘
```
