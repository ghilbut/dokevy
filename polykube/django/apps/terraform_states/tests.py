from django.test import TestCase
from django.urls import reverse
from http import HTTPMethod, HTTPStatus
from pathlib import Path


__base_dir__ = Path(__file__).parent.absolute()
__fixtures_dir__ = __base_dir__ / 'fixtures'


class TerraformStateTest(TestCase):
    fixtures = [
        __fixtures_dir__ / 'states.yaml'
    ]


    @classmethod
    def setUpTestData(cls):
        pass


    @classmethod
    def tearDownClass(cls):
        pass


    def setUp(self) -> None:
        pass


    def tearDown(self) -> None:
        pass


    def test_not_found(self):
        response = self.client.get(reverse('terraform_states:item', kwargs={'name': '00'}))
        self.assertEqual(response.status_code, HTTPStatus.NOT_FOUND)


    def test_method_not_allowed(self):
        response = self.client.delete(reverse('terraform_states:item', kwargs={'name': '01'}))
        self.assertEqual(response.status_code, HTTPStatus.METHOD_NOT_ALLOWED, f'Method is {HTTPMethod.DELETE}')
        response = self.client.patch(reverse('terraform_states:item', kwargs={'name': '01'}))
        self.assertEqual(response.status_code, HTTPStatus.METHOD_NOT_ALLOWED, f'Method is {HTTPMethod.PATCH}')
        response = self.client.put(reverse('terraform_states:item', kwargs={'name': '01'}))
        self.assertEqual(response.status_code, HTTPStatus.METHOD_NOT_ALLOWED, f'Method is {HTTPMethod.PUT}')


    def test_not_implemented(self):
        response = self.client.post(reverse('terraform_states:item', kwargs={'name': '01'}))
        self.assertEqual(response.status_code, HTTPStatus.NOT_IMPLEMENTED, f'Method is {HTTPMethod.POST}')
        response = self.client.generic('LOCK', reverse('terraform_states:item', kwargs={'name': '01'}), content_type='plain/text')
        self.assertEqual(response.status_code, HTTPStatus.NOT_IMPLEMENTED, f'Method is LOCK')
        response = self.client.generic('UNLOCK', reverse('terraform_states:item', kwargs={'name': '01'}), content_type='plain/text')
        self.assertEqual(response.status_code, HTTPStatus.NOT_IMPLEMENTED, f'Method is UNLOCK')
