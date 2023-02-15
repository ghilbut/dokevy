import json
from django.core.validators import slug_re
from django.db import IntegrityError
from django.http import HttpResponse, HttpResponseNotAllowed, HttpResponseNotFound
from django.shortcuts import redirect, render
from django.urls import reverse
from http import HTTPMethod, HTTPStatus
from .forms import StateNewForm
from .models import StateModel


def list(request):
    items = StateModel.objects.all()
    context = { 'items': items }
    return render(request, 'apps/terraform_states/list.html', context)


def new(request):
    if request.method == HTTPMethod.GET:
        return render(request, 'apps/terraform_states/new.html', context={}, content_type='text/html', status=HTTPStatus.OK)
    if request.method == HTTPMethod.POST:
        try:
            name = request.POST['name']
            desc = request.POST['desc']
            if name is None or name == '':
                raise KeyError('name is empty')
            if not slug_re.match(name):
                raise ValueError('name can be only with alphanumeric, hyphen and underscore')
            state = StateModel.objects.create(name=name, desc=desc)
            return redirect(reverse('terraform_states:item', kwargs={ 'name': state.name }))
        except (KeyError, ValueError, IntegrityError) as e:
            return render(request, 'apps/terraform_states/new.html', context={ 'name': name, 'desc': desc, 'error': e }, content_type='text/html', status=HTTPStatus.BAD_REQUEST)
    return HttpResponseNotAllowed(request.method)


def item(request, name):

    method = request.method

    if method in {HTTPMethod.POST, 'LOCK', 'UNLOCK'}:
        return __not_implemented__(request)

    if method == HTTPMethod.GET:
        return __get__(request, name)

    return HttpResponseNotAllowed(method)


def __get__(request, name):
    try:
        q = StateModel.objects.get(name=name)
        print(json.dumps(q.state, indent=4))
    except StateModel.DoesNotExist:
        return HttpResponseNotFound()

    return HttpResponse('OK')


def __not_implemented__(request):
    status = HTTPStatus.NOT_IMPLEMENTED
    return HttpResponse(status.description, content_type='plain/text', reason=status.phrase, status=status)


def edit(request, name):
    pass


def delete(request, name):
    pass
