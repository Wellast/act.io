ansible-playbook -i inventory.aws_ec2.yml cs2server.yml \
    -e "region=eu-central-1" \
    -e "reqOwner=project1" \
    -e "reqName=mycs2server" \
    -e "steamUser=" \
    -e "steamPass=" \
