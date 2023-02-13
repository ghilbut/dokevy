from django.http import HttpResponse, HttpResponseNotAllowed
from http import HTTPMethod, HTTPStatus
from . import models


def list(request):
    return HttpResponse('OK')


def item(request, name):

    method = request.method
    print(method)

    if method in {HTTPMethod.GET, str(HTTPMethod.POST), 'LOCK', 'UNLOCK'}:
        return __not_implemented__(request)

    if request.method == 'GET':
        return __get__(request, name)

    return HttpResponseNotAllowed(method)


def __get__(request, name):
    pass


def __not_implemented__(request):
    status = HTTPStatus.NOT_IMPLEMENTED
    return HttpResponse(status.description, content_type='plain/text', reason=status.phrase, status=status)
