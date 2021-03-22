# Testing

## Provider Tests
In order to test the provider, you can simply run `make test`.

```bash
$ make test
```

## Acceptance Tests

You can run the complete suite of split acceptance tests by doing the following:

```bash
$ make testacc TEST="./split/" 2>&1 | tee test.log
```

To run a single acceptance test in isolation replace the last line above with:

```bash
$ make testacc TEST="./split/" TESTARGS='-run=TestAccSplitEnvironment_Basic'
```

A set of tests can be selected by passing `TESTARGS` a substring. For example, to run all split tests:

```bash
$ make testacc TEST="./split/" TESTARGS='-run=TestAccSplitEnvironment_Basic'
```

### Test Parameters

The following parameters are available for running the test. The absence of some non-required parameters
will cause certain tests to be skipped.

* **TF_ACC** (`integer`) **Required** - must be set to `1`.
* **SPLIT_API_KEY** (`string`) **Required**  - A valid Split admin API key.
* **SPLIT_WORKSPACE_ID** (`string`) - A valid Split workspace ID.
* **SPLIT_WORKSPACE_NAME** (`string`) - A valid Split workspace name.

**For example:**
```bash
export TF_ACC=1
export SPLIT_API_KEY=<SOME_KEY>
$ make testacc TEST="./TestAccSplitEnvironment_Basic/" 2>&1 | tee test.log
```
