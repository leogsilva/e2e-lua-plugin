#!/bin/bash

INI_FILE=$HOME/.aws/credentials

while IFS=' = ' read key value
do
    if [[ $key == \[*] ]]; then
        section=$key
    elif [[ $value ]] && [[ $section == "[${1}]" ]]; then
        if [[ $key == 'aws_access_key_id' ]]; then
            AWS_ACCESS_KEY_ID=$value
        elif [[ $key == 'aws_secret_access_key' ]]; then
            AWS_SECRET_ACCESS_KEY=$value
        fi
    fi
done < $INI_FILE

kubectl create secret generic aws-secret \
	--from-literal=accesskey=$AWS_ACCESS_KEY_ID \
	--from-literal=secretkey=$AWS_SECRET_ACCESS_KEY \
	--kubeconfig=.kubeconfig \
	-n $2 -oyaml --dry-run=client | kubectl apply -f -