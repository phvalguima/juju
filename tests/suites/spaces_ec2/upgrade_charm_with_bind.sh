run_upgrade_charm_with_bind() {
    echo

    file="${TEST_DIR}/test-upgrade-charm-with-bind-ec2.txt"

    ensure "spaces-upgrade-charm-with-bind-ec2" "${file}"

    ## Setup spaces
    juju reload-spaces
    juju add-space isolated 172.31.254.0/24

    ## Create machine
    add_multi_nic_machine
    juju_machine_id=$(juju show-machine --format json | jq -r '.["machines"] | keys[0]')

    # Deploy test charm to dual-nic machine
    juju deploy cs:~juju-qa/bionic/space-defender-0 --bind "defend-a=alpha defend-b=isolated" --to "${juju_machine_id}"
    unit_index=$(get_unit_index "space-defender")
    wait_for "space-defender" "$(idle_condition "space-defender" 0 "${unit_index}")"

    assert_net_iface_for_endpoint_matches "space-defender" "defend-a" "ens5"
    assert_net_iface_for_endpoint_matches "space-defender" "defend-b" "ens6"

    assert_endpoint_binding_matches "space-defender" "" "alpha"
    assert_endpoint_binding_matches "space-defender" "defend-a" "alpha"
    assert_endpoint_binding_matches "space-defender" "defend-b" "isolated"

    # Upgrade the space-defender charm and modify its bindings
    juju upgrade-charm space-defender --bind "defend-a=alpha defend-b=alpha"
    wait_for "space-defender" "$(idle_condition "space-defender" 0 "${unit_index}")"

    # After the upgrade, defend-a should remain attached to ens5 but
    # defend-b which has now been bound to alpha should also get ens5
    assert_net_iface_for_endpoint_matches "space-defender" "defend-a" "ens5"
    assert_net_iface_for_endpoint_matches "space-defender" "defend-b" "ens5"

    assert_endpoint_binding_matches "space-defender" "" "alpha"
    assert_endpoint_binding_matches "space-defender" "defend-a" "alpha"
    assert_endpoint_binding_matches "space-defender" "defend-b" "alpha"

    destroy_model "spaces-upgrade-charm-with-bind-ec2"
}

test_upgrade_charm_with_bind() {
    if [ "$(skip 'test_upgrade_charm_with_bind')" ]; then
        echo "==> TEST SKIPPED: upgrade charm with --bind"
        return
    fi

    (
        set_verbosity

        cd .. || exit

        run "run_upgrade_charm_with_bind"
    )
}