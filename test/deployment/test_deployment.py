from os import environ
from unittest import TestCase

import pytest

from grainger.dp.testutil.k8s import k8s_objects
from grainger.dp.testutil import flaky_ext

NAMESPACE = environ.get('TESTUTIL_NAMESPACE')
FQDN = environ.get('TESTUTIL_FQDN')
APP_NAME = environ.get('TESTUTIL_APP_NAME')
APP_VERSION = environ.get('TESTUTIL_APP_VERSION')
MAX_RETRIES = int(environ.get('TESTUTIL_MAX_RETRIES', 10))
WAIT_SECONDS = int(environ.get('TESTUTIL_WAIT_SECONDS', 30))



class TestDeployment(TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        assert NAMESPACE, 'env var TESTUTIL_NAMESPACE not set'
        assert FQDN, 'env var TESTUTIL_FQDN not set'
        assert APP_NAME, 'env var TESTUTIL_APP_NAME not set'
        assert APP_VERSION, 'env var TESTUTIL_APP_VERSION not set'

    @pytest.mark.flaky(max_runs=MAX_RETRIES, rerun_filter=flaky_ext.retry_unless_error(timeout=WAIT_SECONDS))
    def test_deployment_is_up_with_new_version(self):
        deployment = k8s_objects.Deployment(APP_NAME, NAMESPACE)
        assert deployment.pods()
        for c in deployment.ready_conditions():
            assert c.get('status') == 'True', 'pod "{}" {} status "{}"'.format(c.get('name'),
                                                                               c.get('type'),
                                                                               c.get('status'))
        for v in deployment.versions():
            assert v['version'] == APP_VERSION, 'pod "{}" version (git SHA1)'.format(v['name'])

    def test_service_port_agrees_with_api_container_port(self):
        service = k8s_objects.Service(APP_NAME, NAMESPACE)
        deployment = k8s_objects.Deployment(APP_NAME, NAMESPACE)
        assert service.first_port().get('targetPort') == deployment.first_port().get('containerPort')
        assert service.first_port().get('name') == deployment.first_port().get('name')

    @pytest.mark.flaky(max_runs=MAX_RETRIES, rerun_filter=flaky_ext.retry_unless_error(timeout=WAIT_SECONDS))
    def test_expected_replica_count(self):
        deployment = k8s_objects.Deployment(APP_NAME, NAMESPACE)
        assert deployment.deployment_status()
        configured_replicas_count = deployment.deployment_status().get('replicas')
        ready_replicas_count = deployment.deployment_status().get('readyReplicas')
        assert configured_replicas_count == ready_replicas_count
