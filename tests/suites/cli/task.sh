test_cli() {
    if [ "$(skip 'test_cli')" ]; then
        echo "==> TEST SKIPPED: CLI tests"
        return
    fi

    set_verbosity

    echo "==> Checking for dependencies"
    check_dependencies juju

    file="${TEST_DIR}/test-local-charms.txt"

    bootstrap "test-local-charms" "${file}"

    test_local_charms

    destroy_controller "test-local-charms"
}
