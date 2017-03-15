[ "${TRAVIS_PULL_REQUEST}" = "false" ] && docker login -e $DOCKER_EMAIL -u $DOCKER_USERNAME -p $DOCKER_PASSWORD || false'
[ "${TRAVIS_PULL_REQUEST}" = "false" ] && docker push ${DOCKER_BUILD} || false'
