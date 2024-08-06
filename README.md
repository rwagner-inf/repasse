#!/bin/bash

if [[ $CI_MERGE_REQUEST_LABELS =~ "debug" ]]; then
    set -x
fi

export ANSIBLE_FORCE_COLOR=true
export ANSIBLE_HOST_KEY_CHECKING=False

INVENTORY_FILE=$1
shift

if [[ ! -f $INVENTORY_FILE ]]; then
    echo "$INVENTORY_FILE not found - no filebeat inventories to deploy"
    exit 0
fi

echo "Inventories:"
cat $INVENTORY_FILE

# Run ansible plabooks in parallel
cat $INVENTORY_FILE | xargs -n1 -I{} -P8 bash -c 'ansible-playbook -i "{}" ansible/filebeat.yml --vault-password-file ~/.ssh/ansible-vault.yml --tags deploy -v > $(basename {}).log 2>&1'
exit_code=$?

# Print out logs
tail -n +1 *.log
grep FAILED *.log && echo "Errors reported!"

exit $exit_code
