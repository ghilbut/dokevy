from django.apps import AppConfig


class TerraformStatesConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    label = 'terraform_states'
    name = 'apps.terraform_states'
    verbose_name = 'Terraform States'
