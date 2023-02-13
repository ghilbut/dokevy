from django.urls import path
from . import views


app_name = 'terraform_states'

urlpatterns = [
    path('', views.list, name='list'),
    path('<slug:name>/', views.item, name='item'),
]