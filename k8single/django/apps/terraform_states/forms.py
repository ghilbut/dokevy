from django import forms


class StateNewForm(forms.Form):
    name = forms.CharField(label='Name', max_length=128, min_length=1, required=True)
    desc = forms.Textarea()
    readonly = forms.BooleanField(required=False)
