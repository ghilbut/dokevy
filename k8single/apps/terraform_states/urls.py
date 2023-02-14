from django.urls import path
from . import views


app_name = 'terraform_states'

urlpatterns = [
    path('', views.list, name='list'),
    path('new/', views.new, name='new'),
    path('<slug:name>/', views.item, name='item'),
    path('<slug:name>/edit/', views.edit, name='edit'),
    path('<slug:name>/delete/', views.delete, name='delete'),
]