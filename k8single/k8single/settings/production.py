# noinspection PyUnresolvedReferences
from . import *

DEBUG = False

ALLOWED_HOSTS = ['*']

del LOGGING['loggers']['django']
LOGGING['loggers']['django.db.backends']['level'] = 'ERROR'
LOGGING['loggers']['elasticapm.errors']['level'] = 'ERROR'
LOGGING['loggers']['k8single']['level'] = 'WARNING'

ELASTIC_APM = ELASTIC_APM | {
    'SERVER_URL': 'http://localhost:8200',
    'ENABLED': False,
    'LOG_LEVEL': 'error',
    'ENVIRONMENT': 'production',
}
