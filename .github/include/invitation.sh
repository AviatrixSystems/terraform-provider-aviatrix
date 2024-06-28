#!/bin/bash

#
#   This script is to help users to invite newcomers to a given GitHub organization
#
#   It will be called by a GitHub action.
#

ORG=${1:-AviatrixSystems}
EMAIL=${2}

if [ "${GITHUB_TOKEN}" == "" ]; then
    echo "GitHub must be set (export GITHUB_TOKEN env. variable)"
    exit 2
fi


if [ "${EMAIL}" != "" ]; then
    # Send the invitation
    curl -s -X POST \
        -H "Authorization: Token ${GITHUB_TOKEN}" \
        -d '{"email": "'${EMAIL}'", "role": "direct_member"}' \
        https://api.github.com/orgs/${ORG}/invitations > out.log

    if [ "$(grep errors out.log)" != "" ]; then
        # error
        echo -n "error: "
        cat out.log | jq -r .errors[].message
    else
        echo "User with email '${EMAIL}' was invited to join the organization '${ORG}'"
    fi

    rm -fr out.log
fi


# Show table with existing invitation (markdown)
invitations=$(curl -s \
    -H "Authorization: Token ${GITHUB_TOKEN}" \
    https://api.github.com/orgs/${ORG}/invitations |\
    jq -r '.[] | .login + ";" + .created_at'\
)

cat > sent-invitations.md << EOF

Pending invitations sent for the Organization ${ORG}:

| Org | User | Date |
|:---:|:----:|------|
EOF
for i in ${invitations}; do
    username=$(echo ${i} | cut -f 1 -d ';')
    created_at=$(echo ${i} | cut -f 2 -d ';' | tr 'T' ' ' | cut -c1-19)

    echo "| ${ORG} | [${username}](https://github.com/${username}) | ${created_at} |" >> sent-invitations.md
done
