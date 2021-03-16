# Testing

## Provider Tests
In order to test the provider, you can simply run `make test`.

```bash
$ make test
```

## Acceptance Tests

You can run the complete suite of <REPLACE_ME> acceptance tests by doing the following:

```bash
$ make testacc TEST="./<REPLACE_ME>/" 2>&1 | tee test.log
```

To run a single acceptance test in isolation replace the last line above with:

```bash
$ make testacc TEST="./<REPLACE_ME>/" TESTARGS='-run=NAME_OF_TEST'
```

A set of tests can be selected by passing `TESTARGS` a substring. For example, to run all <REPLACE_ME> tests:

```bash
$ make testacc TEST="./<REPLACE_ME>/" TESTARGS='-run=NAME_OF_TEST'
```

### Test Parameters

The following parameters are available for running the test. The absence of some of the non-required parameters
will cause certain tests to be skipped.

* **TF_ACC** (`integer`) **Required** - must be set to `1`.

**For example:**
```bash
export TF_ACC=1
$ make testacc TEST="./NAME_OF_TEST/" 2>&1 | tee test.log
```
