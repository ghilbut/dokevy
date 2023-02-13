import json
from django.http import HttpResponse, HttpResponseNotAllowed, HttpResponseNotFound
from http import HTTPMethod, HTTPStatus
from .models import State


def list(request):
    return HttpResponse('OK')


def item(request, name):

    method = request.method

    if method in {HTTPMethod.POST, 'LOCK', 'UNLOCK'}:
        return __not_implemented__(request)

    if method == HTTPMethod.GET:
        return __get__(request, name)

    return HttpResponseNotAllowed(method)


def __get__(request, name):
    try:
        q = State.objects.get(name=name)
        print(json.dumps(q.state, indent=4))
    except State.DoesNotExist:
        return HttpResponseNotFound()

    return HttpResponse('OK')


def __not_implemented__(request):
    status = HTTPStatus.NOT_IMPLEMENTED
    return HttpResponse(status.description, content_type='plain/text', reason=status.phrase, status=status)
