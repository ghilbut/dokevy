from django.core.validators import MinLengthValidator
from django.db import models


def __default_state__():
    return { 'version': 4, 'outputs': {} }


class StateModel(models.Model):
    name = models.CharField(max_length=256, blank=False, db_index=True, null=False, unique=True, validators=[MinLengthValidator(1)])
    desc = models.TextField(max_length=1024, default='')
    readonly = models.BooleanField(default=True)
    state = models.JSONField('Terraform State', default=__default_state__)

    def __str__(self):
        return f'Terraform State: {self.name}'

    class Meta:
        db_table = 'k8single_terraform_states'
