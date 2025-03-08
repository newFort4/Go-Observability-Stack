source .env

envsubst < ./alertmanager.template.yml > ./alertmanager.yml