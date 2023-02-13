from django.urls import include, path
from . import views
from .terraform_states import urls as terraform_states_urls


urlpatterns = [
    path('', views.index, name='index'),
    path('healthz/', views.health, name='health'),
    path('terraform_states/', include('apps.terraform_states.urls', namespace='terraform_states')),
]