#from django.contrib.auth.decorators import login_required
from django.http import HttpResponse
from django.shortcuts import render
#from django.urls import reverse_lazy


#@login_required(login_url=reverse_lazy('login'), redirect_field_name=None)
def index(request):
    context = {}
    return render(request, 'index.html', context)


def health(request):
    return HttpResponse("OK")
