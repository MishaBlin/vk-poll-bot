USER_EMAIL="test@example.com"
USER_PASS="Password1!"
USER_NAME="test"
TEAM_NAME="polls"

docker exec test_mm bash -c 'until mmctl --local system status 2> /dev/null; do echo "waiting for server to become available"; sleep 5; done'
docker exec test_mm mmctl --local team create --name "$TEAM_NAME" --display_name "Polls playground" --email "admin@example.com"
docker exec test_mm mmctl --local user create --email="${USER_EMAIL}" --password="${USER_PASS}" --username="${USER_NAME}"
docker exec test_mm mmctl --local roles system_admin "${USER_NAME}"
docker exec test_mm mmctl --local team users add "${TEAM_NAME}" "${USER_NAME}"

echo "Default user credentials"
echo "Email: ${USER_EMAIL} - Password: ${USER_PASS}"
